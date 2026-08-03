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
	"geekai/service"
	"geekai/service/sms"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const CodeStorePrefix = "/verify/codes/"

type SmsHandler struct {
	BaseHandler
	redis         *redis.Client
	sms           *sms.SmsManager
	smtp          *service.SmtpService
	captcha       *service.CaptchaService
	captchaConfig *types.CaptchaConfig
	baseConfig    *types.BaseConfig
}

func NewSmsHandler(
	app *core.AppServer,
	client *redis.Client,
	sms *sms.SmsManager,
	smtp *service.SmtpService,
	captcha *service.CaptchaService,
	sysConfig *types.SystemConfig) *SmsHandler {
	return &SmsHandler{
		redis:         client,
		sms:           sms,
		captcha:       captcha,
		smtp:          smtp,
		BaseHandler:   BaseHandler{App: app},
		captchaConfig: &sysConfig.Captcha,
		baseConfig:    &sysConfig.Base,
	}
}

// RegisterRoutes 注册路由
func (h *SmsHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/sms/")

	// 公开接口（无需登录）
	rg.POST("code", h.SendCode)
}

// SendCode 发送验证码
func (h *SmsHandler) SendCode(c *gin.Context) {
	var data struct {
		Receiver string `json:"receiver"` // 接收者
		Key      string `json:"key"`
		Dots     string `json:"dots,omitempty"`
		X        int    `json:"x,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if h.captchaConfig.Enabled {
		var check bool
		if data.X != 0 {
			check = h.captcha.SlideCheck(data)
		} else {
			check = h.captcha.Check(data)
		}
		if !check {
			resp.ERROR(c, "请先完人机验证")
			return
		}
	}

	code := utils.RandomNumber(6)
	var err error
	if utils.IsValidEmail(data.Receiver) { // email
		// 检查邮箱后缀是否在白名单
		if len(h.baseConfig.EmailWhiteList) > 0 {
			inWhiteList := false
			for _, suffix := range h.baseConfig.EmailWhiteList {
				if strings.HasSuffix(data.Receiver, suffix) {
					inWhiteList = true
					break
				}
			}
			if !inWhiteList {
				resp.ERROR(c, "邮箱后缀不在白名单中")
				return
			}
		}
		err = h.smtp.SendVerifyCode(data.Receiver, code)
	} else if utils.IsValidMobile(data.Receiver) {
		err = h.sms.GetService().SendVerifyCode(data.Receiver, code)
	} else {
		resp.ERROR(c, "请输入正确的手机号或邮箱")
		return
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 存储验证码，等待后面注册验证
	_, err = h.redis.Set(c, CodeStorePrefix+data.Receiver, code, 0).Result()
	if err != nil {
		resp.ERROR(c, "验证码保存失败")
		return
	}

	if h.App.Debug {
		resp.SUCCESS(c, code)
	} else {
		resp.SUCCESS(c)
	}
}
