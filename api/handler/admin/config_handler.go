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
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/payment"
	"geekai/service/sms"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	handler.BaseHandler
	licenseService *service.LicenseService
	sysConfig      *types.SystemConfig
	alipayService  *payment.AlipayService
	wxpayService   *payment.WxPayService
	epayService    *payment.EPayService
	smsAliyun      *sms.AliYunSmsService
	smsBao         *sms.BaoSmsService
	smsManager     *sms.SmsManager
	localOss       *oss.LocalStorage
	qiniuOss       *oss.QiNiuOss
	aliyunOss      *oss.AliYunOss
	minioOss       *oss.MiniOss
	smtpService    *service.SmtpService
	captchaService *service.CaptchaService
	cozeService    *service.CozeService
	wxLoginService *service.WxLoginService
}

func NewConfigHandler(
	app *core.AppServer,
	db *gorm.DB,
	licenseService *service.LicenseService,
	sysConfig *types.SystemConfig,
	alipayService *payment.AlipayService,
	wxpayService *payment.WxPayService,
	epayService *payment.EPayService,
	smsAliyun *sms.AliYunSmsService,
	smsBao *sms.BaoSmsService,
	smsManager *sms.SmsManager,
	localOss *oss.LocalStorage,
	qiniuOss *oss.QiNiuOss,
	aliyunOss *oss.AliYunOss,
	minioOss *oss.MiniOss,
	smtpService *service.SmtpService,
	captchaService *service.CaptchaService,
	cozeService *service.CozeService,
	wxLoginService *service.WxLoginService,
) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		licenseService: licenseService,
		sysConfig:      sysConfig,
		alipayService:  alipayService,
		wxpayService:   wxpayService,
		epayService:    epayService,
		smsAliyun:      smsAliyun,
		smsBao:         smsBao,
		smsManager:     smsManager,
		localOss:       localOss,
		qiniuOss:       qiniuOss,
		aliyunOss:      aliyunOss,
		minioOss:       minioOss,
		smtpService:    smtpService,
		captchaService: captchaService,
		cozeService:    cozeService,
		wxLoginService: wxLoginService,
	}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/config")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("update/base", h.UpdateBase)
		rg.POST("update/notice", h.UpdateNotice)
		rg.POST("update/captcha", h.UpdateCaptcha)
		rg.POST("update/wechat", h.UpdateWechat)
		rg.POST("update/coze", h.UpdateCoze)
		rg.POST("update/payment", h.UpdatePayment)
		rg.POST("update/sms", h.UpdateSms)
		rg.POST("update/oss", h.UpdateOss)
		rg.POST("update/smtp", h.UpdateStmp)
		rg.GET("get", h.Get)
		rg.POST("license/active", h.Active)
		rg.GET("license/get", h.GetLicense)
		rg.POST("license/apply", h.ApplyFreeLicense)
	}
}

// UpdateBase 更新基础配置
func (h *ConfigHandler) UpdateBase(c *gin.Context) {
	var data types.BaseConfig

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 未授权的话不允许修改版权
	license := h.licenseService.GetLicense()
	if !license.IsActive && data.Copyright != h.sysConfig.Base.Copyright {
		resp.ERROR(c, "未授权系统不允许修改版权信息")
		return
	}

	// 未授权的话不允许修改 Logo
	if !license.IsActive && data.Logo != h.sysConfig.Base.Logo {
		resp.ERROR(c, "未授权系统不允许修改 Logo")
		return
	}

	err := h.Update("system", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	h.sysConfig.Base = data

	resp.SUCCESS(c, data)
}

// UpdateNotice 更新公告配置
func (h *ConfigHandler) UpdateNotice(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update("notice", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// UpdateCaptcha 更新行为验证码配置
func (h *ConfigHandler) UpdateCaptcha(c *gin.Context) {
	var data types.CaptchaConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	data.ApiURL = h.App.Config.GeekApiHost
	err := h.Update("captcha", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if data.Enabled {
		err = utils.CheckGeekAIAPIKey(data.ApiURL, data.ApiKey)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
		h.captchaService.UpdateConfig(data)
	}
	h.sysConfig.Captcha = data
	resp.SUCCESS(c, data)

}

// UpdateCoze 更新 Coze 配置
func (h *ConfigHandler) UpdateCoze(c *gin.Context) {
	var data types.CozeApiConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update("coze", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 验证密钥是否正确
	_, err = h.cozeService.GetAccessToken(&data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	h.sysConfig.Coze = data
	resp.SUCCESS(c, data)
}

// UpdatePayment 更新支付配置
func (h *ConfigHandler) UpdatePayment(c *gin.Context) {
	var data types.PaymentConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var config model.Config
	oldData := types.PaymentConfig{}
	err := h.DB.Where("name", "payment").First(&config).Error
	if err == nil {
		utils.JsonDecode(config.Value, &oldData)
	}

	err = h.Update("payment", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新支付服务配置
	if data.WxPayConfig.Enabled {
		err = h.wxpayService.UpdateConfig(&data.WxPayConfig)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	if data.EPayConfig.Enabled {
		h.epayService.UpdateConfig(&data.EPayConfig)
	}
	if data.AlipayConfig.Enabled {
		err = h.alipayService.UpdateConfig(&data.AlipayConfig)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	h.sysConfig.Payment = data
	resp.SUCCESS(c, data)
}

// UpdateSms 更新短信配置
func (h *ConfigHandler) UpdateSms(c *gin.Context) {
	var data types.SMSConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var config model.Config
	oldData := types.SMSConfig{}
	err := h.DB.Where("name", "sms").First(&config).Error
	if err == nil {
		utils.JsonDecode(config.Value, &oldData)
	}

	err = h.Update("sms", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新服务配置
	if data.Active == sms.AliYun {
		err = h.smsAliyun.UpdateConfig(&data.Aliyun)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	if data.Active == sms.Bao {
		h.smsBao.UpdateConfig(&data.Bao)
	}

	h.smsManager.SetActive(data.Active)

	resp.SUCCESS(c, data)
}

// UpdateOss 更新 Oss 配置
func (h *ConfigHandler) UpdateOss(c *gin.Context) {
	var data types.OSSConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var config model.Config
	oldData := types.OSSConfig{}
	err := h.DB.Where("name", "oss").First(&config).Error
	if err == nil {
		utils.JsonDecode(config.Value, &oldData)
	}

	err = h.Update("oss", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新服务配置
	if data.Active == oss.Local {
		h.localOss.UpdateConfig(&data.Local)
	}

	if data.Active == oss.QiNiu {
		h.qiniuOss.UpdateConfig(&data.QiNiu)
	}

	if data.Active == oss.AliYun {
		err := h.aliyunOss.UpdateConfig(&data.AliYun)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	if data.Active == oss.Minio {
		err := h.minioOss.UpdateConfig(&data.Minio)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	h.sysConfig.OSS = data

	resp.SUCCESS(c, data)
}

// UpdateStmp 更新 Stmp 配置
func (h *ConfigHandler) UpdateStmp(c *gin.Context) {
	var data types.SmtpConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var config model.Config
	oldData := types.SmtpConfig{}
	err := h.DB.Where("name", "smtp").First(&config).Error
	if err == nil {
		utils.JsonDecode(config.Value, &oldData)
	}

	err = h.Update("smtp", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if !data.Equal(&oldData) {
		h.smtpService.UpdateConfig(&data)
	}

	h.sysConfig.Smtp = data
	resp.SUCCESS(c, data)
}

// UpdateWechat 更新微信登录配置
func (h *ConfigHandler) UpdateWechat(c *gin.Context) {
	var data types.WxLoginConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	data.ApiURL = h.App.Config.GeekApiHost
	err := h.Update("wx_login", data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if data.Enabled {
		err = utils.CheckGeekAIAPIKey(data.ApiURL, data.ApiKey)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
		h.wxLoginService.UpdateConfig(data)
	}

	h.sysConfig.WxLogin = data
	resp.SUCCESS(c, data)
}

// Update 更新系统配置
func (h *ConfigHandler) Update(name string, value any) error {
	var config model.Config
	err := h.DB.Where("name", name).First(&config).Error
	if err != nil { // 不存在则创建
		config.Name = name
		config.Value = utils.JsonEncode(value)
		return h.DB.Create(&config).Error
	} else { // 存在则更新
		config.Value = utils.JsonEncode(value)
		return h.DB.Updates(&config).Error
	}

}

// Get 获取指定名称的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	name := c.Query("name")
	var config model.Config
	res := h.DB.Where("name", name).First(&config)
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

// Active 激活系统
func (h *ConfigHandler) Active(c *gin.Context) {
	var data struct {
		License string `json:"license"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	license, err := h.licenseService.ActiveLicense(data.License)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if err := h.Update("license", license); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 更新系统配置
	h.sysConfig.License = *license

	resp.SUCCESS(c, license.MachineId)

}

// ApplyFreeLicense 申请免费License
func (h *ConfigHandler) ApplyFreeLicense(c *gin.Context) {
	license, err := h.licenseService.ApplyFreeLicense()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if err := h.Update("license", license); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if !license.IsActive {
		resp.ERROR(c, "申请免费License失败，每个设备只能申请一次免费 License")
		return
	}
	h.sysConfig.License = *license
	resp.SUCCESS(c)
}

// GetLicense 获取 License 信息
func (h *ConfigHandler) GetLicense(c *gin.Context) {
	license := h.licenseService.GetLicense()
	resp.SUCCESS(c, license)
}

//
