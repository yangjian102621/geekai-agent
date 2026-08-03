package admin

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
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScoreLogHandler struct {
	handler.BaseHandler
}

func NewScoreLogHandler(app *core.AppServer, db *gorm.DB) *ScoreLogHandler {
	return &ScoreLogHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ScoreLogHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/score/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("list", h.List)
		rg.POST("batchRemove", h.BatchRemove)
	}
}

func (h *ScoreLogHandler) List(c *gin.Context) {
	var data struct {
		Username  string `json:"username"`
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
	if data.Username != "" {
		session = session.Where("username", data.Username)
	}
	if data.Type > 0 {
		session = session.Where("type", data.Type)
	}
	if data.StartTime != "" {
		session = session.Where("created_at >= ?", data.StartTime)
	}
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
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

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

	// 统计消费算力总和
	var totalScores float64
	if data.StartTime != "" && data.EndTime != "" {
		session.Where("mark", 0).Select("SUM(amount) as total_sum").Scan(&totalScores)
	}
	resp.SUCCESS(c, gin.H{"data": vo.NewPage(total, data.Page, data.PageSize, list), "stat": totalScores})
}

func (h *ScoreLogHandler) BatchRemove(c *gin.Context) {
	var data struct {
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	h.DB.Where("id IN (?)", data.Ids).Delete(&model.ScoreLog{})
	resp.SUCCESS(c, "删除成功")
}
