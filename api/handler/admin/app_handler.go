package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppHandler struct {
	handler.BaseHandler
	redis          *redis.Client
	captcha        *service.CaptchaService
	cozeService    *service.CozeService
	uploader       *oss.UploaderManager
	sysConfig      *types.SystemConfig
	licenseService *service.LicenseService
	appService     *service.AppService
}

func NewAppHandler(app *core.AppServer, db *gorm.DB, client *redis.Client, captcha *service.CaptchaService, cozeService *service.CozeService, uploader *oss.UploaderManager, sysConfig *types.SystemConfig, licenseService *service.LicenseService, appService *service.AppService) *AppHandler {
	return &AppHandler{
		BaseHandler:    handler.BaseHandler{DB: db, App: app},
		redis:          client,
		captcha:        captcha,
		cozeService:    cozeService,
		uploader:       uploader,
		sysConfig:      sysConfig,
		licenseService: licenseService,
		appService:     appService,
	}
}

// RegisterRoutes 注册路由
func (h *AppHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/app/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("list", h.List)
		rg.POST("save", h.Save)
		rg.POST("set", h.Set)
		rg.GET("remove", h.Remove)
		rg.POST("copy", h.Copy)
		rg.GET("coze/agents", h.GetCozeAgents)
		rg.POST("batch-remove", h.BatchRemove)
		rg.POST("coze/import", h.ImportCozeAgents)
	}
}

func (h *AppHandler) Save(c *gin.Context) {
	var data struct {
		Id            uint               `json:"id"`
		Name          string             `json:"name"`
		Type          string             `json:"type"`
		Enabled       bool               `json:"enabled"`
		Score         int                `json:"score"`
		Icon          string             `json:"icon"`
		Summary       string             `json:"summary"`
		Configs       vo.AppConfig       `json:"configs"`
		Params        []vo.WorkflowParam `json:"params"`
		Cid           uint               `json:"cid"`
		IsHot         bool               `json:"is_hot"`
		BillingMode   types.BillingMode  `json:"billing_mode"`
		BillingConfig vo.BillingConfig   `json:"billing_config"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Name == "" {
		resp.ERROR(c, "应用名称不能为空")
		return
	}

	app := model.App{
		Name:          data.Name,
		Type:          types.AppType(data.Type),
		Enabled:       data.Enabled,
		Score:         data.Score,
		Icon:          data.Icon,
		Configs:       utils.JsonEncode(data.Configs),
		Params:        utils.JsonEncode(data.Params),
		BillingMode:   data.BillingMode,
		BillingConfig: utils.JsonEncode(data.BillingConfig),
		Summary:       data.Summary,
		Cid:           data.Cid,
		IsHot:         data.IsHot,
	}

	if app.Type == types.AppCoze {
		_, err := h.cozeService.GetAccessToken(&types.CozeApiConfig{
			ApiUrl:      data.Configs.ApiUrl,
			AppId:       data.Configs.AppId,
			PublicKeyID: data.Configs.PublicKeyID,
			PrivateKey:  data.Configs.PrivateKey,
		})
		if err != nil {
			resp.ERROR(c, "获取Coze授权失败: "+err.Error())
			return
		}
	}
	license := h.licenseService.GetLicense()

	if app.Type == types.AppDify && !license.IsActive {
		resp.ERROR(c, "当前系统未授权，无法使用Dify应用，请联系管理员开通")
		return
	}

	var err error
	if data.Id > 0 {
		err = h.DB.Model(&model.App{}).Select("name", "type", "enabled", "configs", "params", "score", "icon", "summary", "cid", "billing_mode", "billing_config").Where("id", data.Id).Updates(&app).Error
	} else {
		app.CreatorId = 0
		app.Check = vo.CheckStatusPass // 管理员创建的自动审核通过
		app.CheckNote = "管理员创建"
		err = h.DB.Create(&app).Error
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Copy 复制应用
func (h *AppHandler) Copy(c *gin.Context) {
	var data struct {
		Id uint `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 读取要复制的应用信息
	app := model.App{}
	if err := h.DB.Where("id", data.Id).First(&app).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	app.Id = 0
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.Name = app.Name + "（副本）"

	if err := h.DB.Create(&app).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Remove 删除管理员
func (h *AppHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.appService.RemoveApp(id, 0)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Set 设置应用状态/推荐状态
func (h *AppHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.App{}).Where("id", data.Id).UpdateColumn(data.Name, data.Value)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}

// List 获取应用列表
func (h *AppHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 10)
	name := h.GetTrim(c, "name")
	cid := h.GetInt(c, "cid", 0)
	offset := (page - 1) * pageSize

	session := h.DB.Session(&gorm.Session{}).Where("creator_id = 0")
	if name != "" {
		session = session.Where("name LIKE ?", "%"+name+"%")
	}
	if cid > 0 {
		session = session.Where("cid = ?", cid)
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
		appVo.Params = []vo.WorkflowParam{}
		if app.Configs != "" {
			err = utils.JsonDecode(app.Configs, &appVo.Configs)
			if err != nil {
				logger.Error(err)
				continue
			}
		}
		if app.Params != "" {
			err = utils.JsonDecode(app.Params, &appVo.Params)
			if err != nil {
				logger.Error(err)
				appVo.Params = []vo.WorkflowParam{}
			}
		}
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

// GetCozeAgents 获取 Coze 智能体列表
func (h *AppHandler) GetCozeAgents(c *gin.Context) {
	// 首先确认是否配置了 Coze 的 AppId 和 PublicKeyID
	if h.sysConfig.Coze.PublicKeyID == "" || h.sysConfig.Coze.PrivateKey == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.NotConfig})
		return
	}

	agents, err := h.cozeService.GetAgentList(&h.sysConfig.Coze)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("获取智能体列表失败: %v", err))
		return
	}

	resp.SUCCESS(c, agents)
}

// BatchRemove 批量删除应用
func (h *AppHandler) BatchRemove(c *gin.Context) {
	var data struct {
		Ids []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Ids) == 0 {
		resp.ERROR(c, "请选择要删除的应用")
		return
	}

	res := h.DB.Where("id IN ?", data.Ids).Delete(&model.App{})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// ImportCozeAgents 批量导入 Coze 智能体
func (h *AppHandler) ImportCozeAgents(c *gin.Context) {
	var data struct {
		Agents []struct {
			BotID       string `json:"bot_id"`
			BotName     string `json:"bot_name"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
			Cid         uint   `json:"cid"`
		} `json:"agents"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Agents) == 0 {
		resp.ERROR(c, "请选择要导入的智能体")
		return
	}

	// 获取 Coze 配置
	cozeConfig := h.sysConfig.Coze
	if cozeConfig.PublicKeyID == "" || cozeConfig.PrivateKey == "" {
		resp.ERROR(c, "请先配置 Coze 的 AppId 和 PublicKeyID")
		return
	}

	// 批量导入智能体
	for _, agent := range data.Agents {
		// 检查是否已存在
		var existingApp model.App
		err := h.DB.Where("type = ? AND bot_id = ?", types.AppCoze, agent.BotID).First(&existingApp).Error
		if err == nil || existingApp.Id > 0 {
			// 已存在，则更新秘钥信息
			existingApp.Configs = utils.JsonEncode(vo.AppConfig{
				ApiUrl:      cozeConfig.ApiUrl,
				AppId:       cozeConfig.AppId,
				PublicKeyID: cozeConfig.PublicKeyID,
				PrivateKey:  cozeConfig.PrivateKey,
				BotId:       agent.BotID,
			})
			h.DB.Save(&existingApp)
			continue
		}

		// 创建新的应用
		configs := vo.AppConfig{
			ApiUrl:      cozeConfig.ApiUrl,
			BotId:       agent.BotID,
			AppId:       cozeConfig.AppId,
			PublicKeyID: cozeConfig.PublicKeyID,
			PrivateKey:  cozeConfig.PrivateKey,
		}

		// 下载应用图标
		handler := h.uploader.GetUploadHandler()
		icon, err := handler.PutUrlFile(agent.Icon, false)
		if err != nil {
			logger.Error("下载应用图标失败: ", err)
			continue
		}
		app := model.App{
			Name:      agent.BotName,
			Type:      types.AppCoze,
			Enabled:   true,
			BotId:     agent.BotID,
			Score:     1,
			Icon:      icon,
			Configs:   utils.JsonEncode(configs),
			Summary:   agent.Description,
			Cid:       agent.Cid,
			Check:     vo.CheckStatusPass,
			CheckNote: "管理员导入",
		}

		if err := h.DB.Create(&app).Error; err != nil {
			resp.ERROR(c, fmt.Sprintf("导入智能体失败: %v", err))
			return
		}
	}

	resp.SUCCESS(c)
}
