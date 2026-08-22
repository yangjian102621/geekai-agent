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
	"geekai/handler"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatorAppHandler struct {
	handler.BaseHandler
	appService       *service.AppService
	appConfigService *service.AppConfigService
}

func NewCreatorAppHandler(app *core.AppServer, db *gorm.DB, appService *service.AppService, appConfigService *service.AppConfigService) *CreatorAppHandler {
	return &CreatorAppHandler{
		BaseHandler:      handler.BaseHandler{DB: db, App: app},
		appService:       appService,
		appConfigService: appConfigService,
	}
}

func (h *CreatorAppHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/creator/app")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("/list", h.List)
		rg.POST("/check", h.Check)
		rg.GET("/remove", h.Remove)
	}
}

// List 获取应用列表
func (h *CreatorAppHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 10)
	name := h.GetTrim(c, "name")
	check := c.Query("check")
	offset := (page - 1) * pageSize

	session := h.DB.Session(&gorm.Session{}).Where("creator_id > ?", 0)
	if name != "" {
		session = session.Where("name LIKE ?", "%"+name+"%")
	}
	if check != "" {
		session = session.Where("`check` = ?", check)
	}

	var total int64
	if err := session.Model(&model.App{}).Count(&total).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var apps []model.App
	if err := session.Order("id DESC").Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var categories []model.AppCategory
	if err := h.DB.Order("id DESC").Find(&categories).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var categoryMap = make(map[uint]string)
	for _, category := range categories {
		categoryMap[category.Id] = category.Name
	}

	var appVos []vo.App
	for _, app := range apps {
		var appVo vo.App
		err := utils.CopyObject(app, &appVo)
		if err != nil {
			continue
		}
		if app.Configs != "" {
			err = h.appConfigService.Decode(app.Configs, &appVo.Configs)
			if err != nil {
				logger.Error(err)
				continue
			}
		}
		appVo.Configs = h.appConfigService.Mask(appVo.Configs)
		appVo.Id = app.Id
		appVo.CreatedAt = app.CreatedAt.Unix()
		appVo.UpdatedAt = app.UpdatedAt.Unix()
		appVo.Cname = categoryMap[app.Cid]
		appVos = append(appVos, appVo)
	}

	resp.SUCCESS(c, gin.H{
		"items": appVos,
		"total": total,
	})
}

// Check 审核应用
func (h *CreatorAppHandler) Check(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id == 0 {
		resp.ERROR(c, "无效的应用ID")
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

	var app model.App
	if err := h.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ERROR(c, "应用不存在")
		} else {
			resp.ERROR(c, "获取应用信息失败: "+err.Error())
		}
		return
	}

	// 更新审核状态
	updates := map[string]any{
		"check":      data.Check,
		"check_note": data.CheckNote,
	}

	err := h.DB.Model(&app).Updates(updates).Error
	if err != nil {
		resp.ERROR(c, "审核操作失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, "审核成功")
}

// Remove 删除应用
func (h *CreatorAppHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	creatorId := h.GetInt(c, "creator_id", 0)
	if id == 0 || creatorId == 0 {
		resp.ERROR(c, "无效的应用ID")
		return
	}

	err := h.appService.RemoveApp(id, uint(creatorId))
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
