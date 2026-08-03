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
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScoreLogHandler struct {
	BaseHandler
}

func NewScoreLogHandler(app *core.AppServer, db *gorm.DB) *ScoreLogHandler {
	return &ScoreLogHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ScoreLogHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/score/")

	// 需要用户登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("list", h.List)
	}
}

func (h *ScoreLogHandler) List(c *gin.Context) {
	var data struct {
		Type      int    `json:"type"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	userId := h.GetLoginUserId(c)
	session = session.Where("user_id", userId)
	if data.Type != 0 {
		session = session.Where("type", data.Type)
	}
	if data.StartTime != "" {
		session = session.Where("created_at >= ?", data.StartTime)
	}
	// Fix: 无法查出当天的数据
	if data.EndTime != "" {
		// endTime + 24h，确保能查出当天的数据
		t, err := time.Parse("2006-01-02", data.EndTime)
		if err == nil {
			t = t.Add(24 * time.Hour)
			session = session.Where("created_at < ?", t)
		}
	}

	var total int64
	session.Model(&model.ScoreLog{}).Count(&total)
	var items []model.ScoreLog
	var list = make([]vo.ScoreLog, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var log vo.ScoreLog
			err := utils.CopyObject(item, &log)
			if err != nil {
				continue
			}
			log.Id = item.Id
			log.CreatedAt = item.CreatedAt.Unix()
			log.TypeStr = item.Type.String()
			list = append(list, log)
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}
