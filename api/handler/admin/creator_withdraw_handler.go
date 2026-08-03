package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"time"

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

type CreatorWithdrawHandler struct {
	handler.BaseHandler
	creatorService *service.CreatorService
}

func NewWithdrawHandler(app *core.AppServer, db *gorm.DB, creatorService *service.CreatorService) *CreatorWithdrawHandler {
	return &CreatorWithdrawHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		creatorService: creatorService,
	}
}

func (h *CreatorWithdrawHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/creator/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("/withdraws", h.List)
		rg.POST("/withdraws/proccess", h.Proccess)
	}
}

// List 获取提现申请列表
func (h *CreatorWithdrawHandler) List(c *gin.Context) {
	var data struct {
		Page      int    `json:"page" binding:"required"`
		PageSize  int    `json:"page_size" binding:"required"`
		Status    string `json:"status"`
		Method    string `json:"method"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		CreatorId string `json:"creator_id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	var withdraws []model.CreatorWithdraw
	var total int64

	// 构建查询
	query := h.DB.Session(&gorm.Session{})
	// 筛选条件
	if data.Status != "" {
		query = query.Where("status = ?", data.Status)
	}
	if data.Method != "" {
		query = query.Where("method = ?", data.Method)
	}
	if data.StartTime != "" {
		query = query.Where("created_at >= ?", data.StartTime)
	}
	if data.EndTime != "" {
		// endTime + 24h，确保能查出当天的数据
		t, err := time.Parse("2006-01-02", data.EndTime)
		if err == nil {
			t = t.Add(24 * time.Hour)
			query = query.Where("created_at < ?", t)
		}
	}
	if data.CreatorId != "" {
		query = query.Where("creator_id = ?", data.CreatorId)
	}

	// 统计总数
	query.Model(&model.CreatorWithdraw{}).Count(&total)

	// 分页查询
	query.Offset((data.Page - 1) * data.PageSize).Limit(data.PageSize).Order("id DESC").Find(&withdraws)

	var items []vo.CreatorWithdraw
	for _, withdraw := range withdraws {
		var item vo.CreatorWithdraw
		err := utils.CopyObject(withdraw, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = withdraw.CreatedAt.Unix()
		item.UpdatedAt = withdraw.UpdatedAt.Unix()
		items = append(items, item)
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}

// Proccess 处理提现申请
func (h *CreatorWithdrawHandler) Proccess(c *gin.Context) {
	var data struct {
		Id     uint   `json:"id" binding:"required"`
		Status string `json:"status" binding:"required"`
		Note   string `json:"note" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	// 查询提现申请
	var withdraw model.CreatorWithdraw
	if err := h.DB.First(&withdraw, data.Id).Error; err != nil {
		resp.ERROR(c, "提现申请不存在")
		return
	}

	if withdraw.Status != types.WithdrawStatusPending {
		resp.ERROR(c, "该申请已处理")
		return
	}

	// 更新提现申请状态
	if err := h.DB.Model(&withdraw).Updates(map[string]any{
		"status": data.Status,
		"note":   data.Note,
	}).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 如果状态为拒绝，则退回创作者积分
	if data.Status == types.WithdrawStatusReject {
		h.creatorService.IncreaseScores(withdraw.CreatorId, withdraw.Scores, model.CreatorScoreLog{
			Type:    types.CreatorScoreTypeRefund,
			Subject: "提现失败",
			Remark:  fmt.Sprintf("提现被拒绝，退回积分：%d", withdraw.Scores),
		})
	}

	resp.SUCCESS(c)
}
