package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppHandler struct {
	BaseHandler
	redis     *redis.Client
	sysConfig *types.SystemConfig
}

func NewAppHandler(app *core.AppServer, db *gorm.DB, client *redis.Client, sysConfig *types.SystemConfig) *AppHandler {
	return &AppHandler{
		BaseHandler: BaseHandler{DB: db, App: app},
		redis:       client,
		sysConfig:   sysConfig,
	}
}

// RegisterRoutes 注册路由
func (h *AppHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/app/")

	// 公开接口（无需登录）
	rg.GET("list", h.List)
	rg.GET("hot", h.HotApps)
	rg.GET("category/list", h.GetAppCategories)
}

// List 获取应用列表
func (h *AppHandler) List(c *gin.Context) {
	cid := h.GetInt(c, "cid", 0)
	keyword := h.GetTrim(c, "keyword")
	creatorId := h.GetInt(c, "creator_id", 0)
	var apps []model.App
	session := h.DB.Session(&gorm.Session{}).Where("enabled = ? AND `check` = ?", true, vo.CheckStatusPass)
	if cid > 0 {
		session = session.Where("cid", cid)
	}
	if keyword != "" {
		session = session.Where("name LIKE ?", "%"+keyword+"%")
	}
	if creatorId > 0 {
		session = session.Where("creator_id", creatorId)
	}
	if err := session.Find(&apps).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var appVos []vo.App
	for _, app := range apps {
		var appVo vo.App
		err := utils.CopyObject(app, &appVo)
		if err != nil {
			continue
		}
		// 应用配置包含第三方密钥，公开接口只返回应用元数据。
		appVo.Configs = vo.AppConfig{}
		appVo.CreatedAt = app.CreatedAt.Unix()
		appVo.UpdatedAt = app.UpdatedAt.Unix()
		appVos = append(appVos, appVo)
	}

	resp.SUCCESS(c, appVos)

}

// HotApps 获取热门应用列表
func (h *AppHandler) HotApps(c *gin.Context) {
	var apps []model.App
	session := h.DB.Session(&gorm.Session{}).Where("enabled = ? AND `check` = ?", true, vo.CheckStatusPass).Where("is_hot = ?", true)
	session = session.Debug().Order("id DESC").Limit(10)
	if err := session.Find(&apps).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var appVos []vo.App
	for _, app := range apps {
		var appVo vo.App
		err := utils.CopyObject(app, &appVo)
		if err != nil {
			continue
		}
		// 热门应用接口同样是公开接口，禁止返回应用密钥。
		appVo.Configs = vo.AppConfig{}
		appVo.CreatedAt = app.CreatedAt.Unix()
		appVo.UpdatedAt = app.UpdatedAt.Unix()
		appVos = append(appVos, appVo)
	}

	resp.SUCCESS(c, appVos)
}

// 获取应用分类
func (h *AppHandler) GetAppCategories(c *gin.Context) {
	var categories []model.AppCategory
	creatorId := h.GetInt(c, "creator_id", 0)
	session := h.DB.Session(&gorm.Session{}).Where("enabled", true)
	if creatorId >= 0 {
		session = session.Where("creator_id", creatorId)
	}
	if err := session.Find(&categories).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var categoryVos []vo.AppCategory
	for _, category := range categories {
		var categoryVo vo.AppCategory
		err := utils.CopyObject(category, &categoryVo)
		if err != nil {
			continue
		}
		categoryVos = append(categoryVos, categoryVo)
	}

	resp.SUCCESS(c, categoryVos)
}
