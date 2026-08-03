package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/service"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LicenseHandler struct {
	BaseHandler
	licenseService *service.LicenseService
}

func NewLicenseHandler(app *core.AppServer, db *gorm.DB, licenseService *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{BaseHandler: BaseHandler{App: app, DB: db}, licenseService: licenseService}
}

// RegisterRoutes 注册路由
func (h *LicenseHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/license/")

	// 公开接口（无需登录）
	rg.GET("get", h.Get)
}

// Get 获取 License 信息
func (h *LicenseHandler) Get(c *gin.Context) {
	license := h.licenseService.GetLicense()
	if license == nil {
		resp.ERROR(c, "License not found")
		return
	}

	resp.SUCCESS(c, gin.H{
		"name":       license.Name,
		"active_at":  license.ActiveAt,
		"is_active":  license.IsActive,
		"expired_at": license.ExpiredAt,
		"configs":    license.Configs,
	})
}
