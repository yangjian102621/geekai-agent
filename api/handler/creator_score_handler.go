package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScoreHandler struct {
	BaseHandler
}

func NewScoreHandler(app *core.AppServer, db *gorm.DB) *ScoreHandler {
	return &ScoreHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

func (h *ScoreHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/creator/scores/")

	// 需要创作者登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("/list", h.List)
	}
}

// List 获取创作者收益日志列表
func (h *ScoreHandler) List(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var data struct {
		Page      int    `json:"page" binding:"required"`
		PageSize  int    `json:"page_size" binding:"required"`
		Type      string `json:"type"` // 收益类型 income/withdraw
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	// 构建查询
	query := h.DB.Session(&gorm.Session{}).Where("creator_id = ?", creator.Id)
	if data.Type != "" {
		query = query.Where("type = ?", data.Type)
	}
	if data.StartDate != "" {
		query = query.Where("created_at >= ?", data.StartDate)
	}
	if data.EndDate != "" {
		query = query.Where("created_at <= ?", data.EndDate)
	}

	var total int64
	query.Model(&model.CreatorScoreLog{}).Count(&total)

	var logs []model.CreatorScoreLog
	query.Order("id DESC").Offset((data.Page - 1) * data.PageSize).Limit(data.PageSize).Find(&logs)

	var items []vo.CreatorScoreLog
	for _, log := range logs {
		var item vo.CreatorScoreLog
		err := utils.CopyObject(log, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = log.CreatedAt.Unix()
		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}
