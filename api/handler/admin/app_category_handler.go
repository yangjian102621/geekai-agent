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

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type AppCategoryHandler struct {
	handler.BaseHandler
}

func NewAppCategoryHandler(app *core.AppServer, db *gorm.DB) *AppCategoryHandler {
	return &AppCategoryHandler{
		BaseHandler: handler.BaseHandler{DB: db, App: app},
	}
}

// RegisterRoutes 注册路由
func (h *AppCategoryHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/app/category/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("list", h.List)
		rg.POST("save", h.Save)
		rg.POST("enable", h.Enable)
		rg.GET("remove", h.Remove)
	}
}

// List 数据列表
func (h *AppCategoryHandler) List(c *gin.Context) {
	system := h.GetBool(c, "system")
	var items []model.AppCategory
	session := h.DB.Session(&gorm.Session{})
	if system {
		session = session.Where("creator_id = ?", 0)
	} else {
		session = session.Where("creator_id > ?", 0)
	}
	res := session.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	categories := make([]vo.AppCategory, 0)
	for _, item := range items {
		var u vo.AppCategory
		err := utils.CopyObject(item, &u)
		if err != nil {
			continue
		}
		u.Id = item.Id
		u.CreatedAt = item.CreatedAt.Unix()
		u.UpdatedAt = item.UpdatedAt.Unix()
		categories = append(categories, u)
	}

	resp.SUCCESS(c, categories)
}

// Save 保存
func (h *AppCategoryHandler) Save(c *gin.Context) {
	var data struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Name == "" {
		resp.ERROR(c, "分类名称不能为空")
		return
	}

	// 更新
	var err error
	if data.Id > 0 {
		err = h.DB.Model(&model.AppCategory{}).Where("id", data.Id).Select("name", "enabled").Updates(&model.AppCategory{
			Name:    data.Name,
			Enabled: data.Enabled,
		}).Error
	} else {
		err = h.DB.Create(&model.AppCategory{
			Name:    data.Name,
			Enabled: data.Enabled,
		}).Error
	}
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Remove 删除管理员
func (h *AppCategoryHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Where("id", id).Delete(&model.AppCategory{})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// Enable 启用/禁用
func (h *AppCategoryHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.AppCategory{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}
