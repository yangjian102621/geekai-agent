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
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	BaseHandler
	captchaConfig *types.CaptchaConfig
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, sysConfig *types.SystemConfig) *ConfigHandler {
	return &ConfigHandler{BaseHandler: BaseHandler{App: app, DB: db}, captchaConfig: &sysConfig.Captcha}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/config/")

	// 公开接口（无需登录）
	rg.GET("get", h.Get)
	rg.GET("captcha", h.GetCaptcha)
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	name := c.Query("name")
	if name != "system" && name != "notice" {
		resp.ERROR(c, "invalid name")
		return
	}
	var config model.Config
	res := h.DB.Where("name = ?", name).First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value map[string]any
	err := utils.JsonDecode(config.Value, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, value)
}

// 获取验证码配置
func (h *ConfigHandler) GetCaptcha(c *gin.Context) {
	var config model.Config
	res := h.DB.Where("name = ?", "captcha").First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value types.CaptchaConfig
	err := utils.JsonDecode(config.Value, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// fix:如果没有配置 API-KEY 则自动禁用行为验证码，修复无法登录管理后台的 Bug
	if value.ApiKey == "" {
		value.Enabled = false
		h.captchaConfig.Enabled = false
	}

	resp.SUCCESS(c, gin.H{
		"enabled": value.Enabled,
		"type":    value.Type,
	})
}
