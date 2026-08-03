package admin

import (
	"fmt"

	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type CreatorHandler struct {
	handler.BaseHandler
	creatorService *service.CreatorService
}

func NewCreatorHandler(app *core.AppServer, db *gorm.DB, creatorService *service.CreatorService) *CreatorHandler {
	return &CreatorHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		creatorService: creatorService,
	}
}

func (h *CreatorHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/creator")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("/list", h.List)
		rg.POST("/check", h.Check)
		rg.POST("/enable", h.Enable)
		rg.POST("/update", h.Update)
		rg.GET("/delete", h.Delete)
		rg.POST("/score", h.ChangeScore)
	}
}

// List 获取创作者列表（支持按审核状态筛选）
func (h *CreatorHandler) List(c *gin.Context) {
	page := utils.IntValue(c.Query("page"), 1)
	pageSize := utils.IntValue(c.Query("page_size"), 20)
	session := h.DB.Session(&gorm.Session{})
	// 搜索功能
	name := c.Query("name")
	if name != "" {
		session = session.Where("name LIKE ?", "%"+name+"%")
	}
	check := c.Query("check")
	if check != "" {
		session = session.Where("`check` = ?", check)
	}

	var total int64
	session.Model(&model.Creator{}).Count(&total)

	var creators []model.Creator
	offset := (page - 1) * pageSize
	err := session.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&creators).Error
	if err != nil {
		resp.ERROR(c, "获取创作者列表失败: "+err.Error())
		return
	}

	// 关联用户信息
	userIds := make([]uint, 0)
	userMap := make(map[uint]model.User)
	for _, creator := range creators {
		userIds = append(userIds, creator.UserId)
	}
	var users []model.User
	err = h.DB.Where("id IN (?)", userIds).Find(&users).Error
	if err != nil {
		resp.ERROR(c, "获取用户信息失败: "+err.Error())
		return
	}
	for _, user := range users {
		userMap[user.Id] = user
	}

	// 转换为 VO
	var items []vo.Creator
	for _, creator := range creators {
		var item vo.Creator
		err := utils.CopyObject(creator, &item)
		if err != nil {
			continue
		}
		// 转换时间格式
		item.CreatedAt = creator.CreatedAt.Unix()
		item.UpdatedAt = creator.UpdatedAt.Unix()
		item.Username = userMap[creator.UserId].Username
		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

// Check 审核创作者
func (h *CreatorHandler) Check(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id == 0 {
		resp.ERROR(c, "无效的创作者ID")
		return
	}

	var data struct {
		Check     int8   `json:"check" binding:"required"` // 1: 通过, -1: 不通过
		CheckNote string `json:"check_note"`               // 审核备注
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数解析失败: "+err.Error())
		return
	}

	// 如果是拒绝，必须填写备注
	if data.Check == vo.CheckStatusReject && data.CheckNote == "" {
		resp.ERROR(c, "审核不通过时必须说明原因")
		return
	}

	var creator model.Creator
	if err := h.DB.First(&creator, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ERROR(c, "创作者不存在")
		} else {
			resp.ERROR(c, "获取创作者信息失败: "+err.Error())
		}
		return
	}

	// 更新审核状态
	updates := map[string]any{
		"check":      data.Check,
		"check_note": data.CheckNote,
	}

	// 如果审核通过，默认启用
	if data.Check == vo.CheckStatusPass {
		updates["enabled"] = true
	} else {
		updates["enabled"] = false
	}

	err := h.DB.Model(&creator).Updates(updates).Error
	if err != nil {
		resp.ERROR(c, "审核操作失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, "审核成功")
}

// Enable 启用/禁用创作者
func (h *CreatorHandler) Enable(c *gin.Context) {
	var data struct {
		Id      int  `json:"id"`
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数解析失败: "+err.Error())
		return
	}

	err := h.DB.Model(&model.Creator{}).Where("id = ?", data.Id).Update("enabled", data.Enabled).Error
	if err != nil {
		resp.ERROR(c, "操作失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, "操作成功")
}

// Update 更新创作者信息
func (h *CreatorHandler) Update(c *gin.Context) {
	var data struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		Fee         int    `json:"fee"`
		Check       int    `json:"check"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数解析失败: "+err.Error())
		return
	}

	var creator model.Creator
	if err := h.DB.Where("id = ?", data.Id).First(&creator).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ERROR(c, "创作者不存在")
		} else {
			resp.ERROR(c, "获取创作者信息失败: "+err.Error())
		}
		return
	}

	// 更新字段
	updates := make(map[string]any)
	updates["check"] = data.Check
	if data.Name != "" {
		updates["name"] = data.Name
	}
	if data.Description != "" {
		updates["description"] = data.Description
	}
	if data.Logo != "" {
		updates["logo"] = data.Logo
	}
	if data.Fee <= 0 || data.Fee >= 100 {
		resp.ERROR(c, "提现费率必须在0-100之间")
		return
	} else {
		updates["fee"] = data.Fee
	}

	if len(updates) > 0 {
		if err := h.DB.Model(&creator).Updates(updates).Error; err != nil {
			resp.ERROR(c, "更新创作者失败: "+err.Error())
			return
		}
	}

	resp.SUCCESS(c, "创作者更新成功")
}

// Delete 删除创作者
func (h *CreatorHandler) Delete(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id == 0 {
		resp.ERROR(c, "无效的创作者ID")
		return
	}

	tx := h.DB.Begin()

	if err := tx.Where("id = ?", id).Delete(&model.Creator{}).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "删除创作者失败: "+err.Error())
		return
	}

	// 删除创作者应用分类
	if err := tx.Where("creator_id = ?", id).Delete(&model.AppCategory{}).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "删除创作者应用分类失败: "+err.Error())
		return
	}

	var apps []model.App
	if err := tx.Where("creator_id = ?", id).Find(&apps).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "获取创作者创建的应用失败: "+err.Error())
		return
	}

	// 删除创作者创建的智能体
	if err := tx.Where("creator_id = ?", id).Delete(&model.App{}).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "删除创作者创建的应用失败: "+err.Error())
		return
	}

	var appIds []uint
	for _, app := range apps {
		appIds = append(appIds, app.Id)
	}

	// 删除创作者创建的对话
	if err := tx.Where("app_id IN (?)", appIds).Delete(&model.ChatItem{}).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "删除创作者创建的对话失败: "+err.Error())
		return
	}

	// 删除创作者创建的对话记录
	if err := tx.Where("app_id IN (?)", appIds).Delete(&model.ChatMessage{}).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, "删除创作者创建的对话失败: "+err.Error())
		return
	}

	tx.Commit()
	resp.SUCCESS(c, "创作者删除成功")
}

// ChangeScore 修改创作者积分
func (h *CreatorHandler) ChangeScore(c *gin.Context) {
	var data struct {
		CreatorId int    `json:"creator_id" binding:"required"`
		Action    string `json:"action" binding:"required"` // inc/dec
		Score     int    `json:"score" binding:"required"`
		Remark    string `json:"remark" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数解析失败: "+err.Error())
		return
	}
	if data.Score <= 0 {
		resp.ERROR(c, "积分数值必须大于0")
		return
	}
	var creator model.Creator
	if err := h.DB.First(&creator, data.CreatorId).Error; err != nil {
		resp.ERROR(c, "创作者不存在")
		return
	}
	var err error
	subject := "系统调整"

	switch data.Action {
	case "inc":
		err = h.creatorService.IncreaseScores(uint(data.CreatorId), data.Score, model.CreatorScoreLog{
			Type:    types.CreatorScoreFineTune,
			Subject: subject,
			Remark:  fmt.Sprintf("增加%d积分: %s", data.Score, data.Remark),
		})
	case "dec":
		if creator.Scores < data.Score {
			resp.ERROR(c, "积分不足，无法扣减")
			return
		}
		err = h.creatorService.DecreaseScores(uint(data.CreatorId), data.Score, model.CreatorScoreLog{
			Type:    types.CreatorScoreFineTune,
			Subject: subject,
			Remark:  fmt.Sprintf("扣除%d积分: %s", data.Score, data.Remark),
		})
	default:
		resp.ERROR(c, "不支持的操作类型")
		return
	}
	if err != nil {
		resp.ERROR(c, "操作失败: "+err.Error())
		return
	}
	resp.SUCCESS(c, "操作成功")
}
