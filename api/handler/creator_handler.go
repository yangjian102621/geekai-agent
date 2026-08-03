package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"

	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatorHandler struct {
	BaseHandler
	creatorService *service.CreatorService
}

func NewCreatorHandler(app *core.AppServer, db *gorm.DB, creatorService *service.CreatorService) *CreatorHandler {
	return &CreatorHandler{
		BaseHandler:    BaseHandler{App: app, DB: db},
		creatorService: creatorService,
	}
}

func (h *CreatorHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/creator/")
	// C端展示接口
	rg.GET("/:username", h.GetCreatorInfo)
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.GET("/rand/name", h.GenerateName)
		rg.GET("/rand/logo", h.GenerateLogo)
		rg.POST("/apply", h.Apply)
		rg.GET("/status", h.GetStatus)
		// 创作者控制台info相关接口
		rg.GET("/info", h.CreatorInfo)
		rg.POST("/withdraw", h.SubmitWithdraw)
		rg.POST("/update/profile", h.UpdateProfile)
		rg.GET("/check/username", h.CheckUsername)
	}
}

// Apply 申请成为创作者
func (h *CreatorHandler) Apply(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Logo        string `json:"logo" binding:"required"`
		Username    string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	// 检查是否已经申请过
	var existCreator model.Creator
	err := h.DB.Where("user_id = ?", userId).First(&existCreator).Error
	if err == nil {
		// 已经有记录了
		switch existCreator.Check {
		case 0:
			resp.ERROR(c, "您的申请正在审核中，请耐心等待")
			return
		case 1:
			resp.ERROR(c, "您已经是创作者了")
			return
		case 2:
			resp.ERROR(c, "您的申请已被拒绝，拒绝原因："+existCreator.CheckNote)
			return
		}
	}

	// 创建新的申请记录
	creator := model.Creator{
		UserId:      uint(userId),
		Username:    req.Username,
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Check:       0, // 待审核
		Enabled:     false,
	}

	if err := h.DB.Create(&creator).Error; err != nil {
		resp.ERROR(c, "申请失败："+err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"message": "申请提交成功，请等待审核",
	})
}

// GetStatus 获取创作者申请状态
func (h *CreatorHandler) GetStatus(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	var creator model.Creator
	err := h.DB.Where("user_id = ?", userId).First(&creator).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.SUCCESS(c, gin.H{
				"status": "not_applied", // 未申请
			})
			return
		}
		resp.ERROR(c, "查询失败："+err.Error())
		return
	}

	var status string
	switch creator.Check {
	case 0:
		status = "pending"
	case 1:
		status = "approved"
	case 2:
		status = "rejected"
	}

	resp.SUCCESS(c, gin.H{
		"status":  status,
		"message": creator.CheckNote,
	})
}

// GenerateName 生成随机名称
func (h *CreatorHandler) GenerateName(c *gin.Context) {
	name := utils.GenerateNickname()
	resp.SUCCESS(c, gin.H{
		"name": name,
	})
}

// GenerateLogo 生成随机头像
func (h *CreatorHandler) GenerateLogo(c *gin.Context) {
	logo := utils.GenerateAvatar()
	resp.SUCCESS(c, gin.H{
		"logo": logo,
	})
}

// CreatorInfo 获取创作者CreatorInfo数据
func (h *CreatorHandler) CreatorInfo(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	// 统计数据
	var appCount int64
	h.DB.Model(&model.App{}).Where("creator_id = ?", h.Creator.Id).Count(&appCount)

	// 统计收益
	var totalEarnings int64
	h.DB.Model(&model.CreatorScoreLog{}).Where("creator_id = ? AND type = ?", creator.Id, types.CreatorScoreTypeIncome).Select("SUM(score)").Scan(&totalEarnings)

	// 今日收益
	var todayEarnings int64
	h.DB.Model(&model.CreatorScoreLog{}).Where("creator_id = ? AND type = ? AND DATE(created_at) = CURDATE()", creator.Id, types.CreatorScoreTypeIncome).Select("SUM(score)").Scan(&todayEarnings)

	var creatorVo vo.Creator
	err := utils.CopyObject(h.Creator, &creatorVo)
	if err != nil {
		resp.ERROR(c, "获取创作者信息失败："+err.Error())
		return
	}

	creatorVo.AppCount = int(appCount)
	creatorVo.TotalEarnings = int(totalEarnings)
	creatorVo.TodayEarnings = int(todayEarnings)
	creatorVo.WithdrawConfigs.ScoreToRMBRatio = types.ScoreToRMBRatio

	resp.SUCCESS(c, creatorVo)
}

// SubmitWithdraw 提现申请
func (h *CreatorHandler) SubmitWithdraw(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var data struct {
		Scores      int     `json:"scores" binding:"required"`
		Method      string  `json:"method" binding:"required"`
		Account     string  `json:"account" binding:"required"`
		AccountName string  `json:"account_name" binding:"required"`
		QrCode      string  `json:"qrcode" binding:"required"`
		TotalMoney  float64 `json:"total_money" binding:"required"`
		RealMoney   float64 `json:"real_money" binding:"required"`
		Fee         float64 `json:"fee"`
		Note        string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}
	if h.Creator.Scores < data.Scores {
		resp.ERROR(c, "积分不足")
		return
	}

	withdrawLog := model.CreatorWithdraw{
		CreatorId:   h.Creator.Id,
		Scores:      data.Scores,
		TotalMoney:  data.TotalMoney,
		RealMoney:   data.RealMoney,
		Fee:         data.Fee,
		Method:      data.Method,
		Account:     data.Account,
		AccountName: data.AccountName,
		QrCode:      data.QrCode,
		Status:      types.WithdrawStatusPending,
		Note:        data.Note,
	}

	tx := h.DB.Begin()
	if err := tx.Create(&withdrawLog).Error; err != nil {
		resp.ERROR(c, "提现失败："+err.Error())
		return
	}
	err := h.creatorService.DecreaseScores(creator.Id, data.Scores, model.CreatorScoreLog{
		Type:    types.CreatorScoreTypeWithdraw,
		Subject: "提现",
		Remark:  fmt.Sprintf("提现%d积分，到账金额%.2f元，手续费%.2f元", data.Scores, data.RealMoney, data.Fee),
	})
	if err != nil {
		resp.ERROR(c, "提现失败："+err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	resp.SUCCESS(c, gin.H{
		"message": "提现申请成功，请等待审核",
	})
}

// UpdateProfile 更新创作者信息
func (h *CreatorHandler) UpdateProfile(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var data struct {
		Name            string         `json:"name" binding:"required"`
		Description     string         `json:"description" binding:"required"`
		Logo            string         `json:"logo" binding:"required"`
		Username        string         `json:"username" binding:"required"`
		WithdrawConfigs map[string]any `json:"withdraw_configs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	// 只更新允许的字段
	updates := map[string]any{
		"name":             data.Name,
		"username":         data.Username,
		"description":      data.Description,
		"logo":             data.Logo,
		"withdraw_configs": utils.JsonEncode(data.WithdrawConfigs),
	}

	if err := h.DB.Model(&model.Creator{}).Where("id = ?", creator.Id).Updates(updates).Error; err != nil {
		resp.ERROR(c, "更新失败："+err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"message": "更新成功",
	})
}

// CheckUsername 检查用户名是否可用
func (h *CreatorHandler) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		resp.ERROR(c, "用户名不能为空")
		return
	}

	var creator model.Creator
	err := h.DB.Where("username = ?", username).First(&creator).Error
	if err == nil {
		resp.ERROR(c, fmt.Sprintf("用户名 %s 已被使用", username))
		return
	}

	resp.SUCCESS(c, gin.H{
		"message": "用户名可用",
	})
}

// GetCreatorInfo 获取创作者信息（C端）
func (h *CreatorHandler) GetCreatorInfo(c *gin.Context) {
	username := c.Param("username")
	var creator model.Creator
	if err := h.DB.Where("username = ? AND enabled = 1 AND `check` = 1", username).First(&creator).Error; err != nil {
		resp.ERROR(c, "创作者不存在或未启用")
		return
	}

	var creatorVo vo.Creator
	err := utils.CopyObject(creator, &creatorVo)
	if err != nil {
		resp.ERROR(c, "获取创作者信息失败："+err.Error())
		return
	}

	resp.SUCCESS(c, creatorVo)
}
