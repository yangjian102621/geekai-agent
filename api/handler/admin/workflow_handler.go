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
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowHandler struct {
	handler.BaseHandler
	redis       *redis.Client
	cozeService *service.CozeService
	uploader    *oss.UploaderManager
	sysConfig   *types.SystemConfig
}

func NewWorkflowHandler(app *core.AppServer,
	db *gorm.DB,
	client *redis.Client,
	cozeService *service.CozeService,
	uploader *oss.UploaderManager,
	sysConfig *types.SystemConfig) *WorkflowHandler {
	return &WorkflowHandler{
		BaseHandler: handler.BaseHandler{DB: db, App: app},
		redis:       client,
		cozeService: cozeService,
		uploader:    uploader,
		sysConfig:   sysConfig,
	}
}

func (h *WorkflowHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/workflow/")
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("save", h.Save)
		rg.GET("remove", h.Remove)
		rg.POST("enable", h.Enable)
		rg.GET("list", h.List)
		rg.POST("batch-remove", h.BatchRemove)
		rg.POST("import", h.ImportWorkflows)
		rg.POST("batch-import", h.BatchImportWorkflows)
	}
}

func (h *WorkflowHandler) Save(c *gin.Context) {
	var data struct {
		Id                uint                    `json:"id"`
		Name              string                  `json:"name"`
		WorkflowId        string                  `json:"workflow_id"`
		Icon              string                  `json:"icon"`
		Enabled           bool                    `json:"enabled"`
		Score             int                     `json:"score"`
		Summary           string                  `json:"summary"`
		Type              string                  `json:"type"`
		Params            []vo.WorkflowParam      `json:"params"`
		AuthConfig        types.CozeApiConfig     `json:"auth_config"`
		BailianAuthConfig types.BailianApiConfig  `json:"bailian_auth_config"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.WorkflowId == "" {
		resp.ERROR(c, "工作流ID不能为空")
		return
	}

	// 根据类型验证授权配置
	workflowType := data.Type
	if workflowType == "" {
		workflowType = "coze" // 默认为 Coze
	}

	var authConfigStr string
	if workflowType == "bailian" {
		authConfigStr = utils.JsonEncode(data.BailianAuthConfig)
	} else {
		// Coze 类型：验证授权配置
		if data.AuthConfig.AppId != "" || data.AuthConfig.PrivateKey != "" {
			_, err := h.cozeService.GetAccessToken(&data.AuthConfig)
			if err != nil {
				resp.ERROR(c, "获取Coze授权失败: "+err.Error())
				return
			}
		}
		authConfigStr = utils.JsonEncode(data.AuthConfig)
	}

	app := model.Workflow{
		Name:       data.Name,
		Icon:       data.Icon,
		Enabled:    data.Enabled,
		Score:      data.Score,
		Params:     utils.JsonEncode(data.Params),
		AuthConfig: authConfigStr,
		Summary:    data.Summary,
		WorkflowId: data.WorkflowId,
		Type:       workflowType,
	}

	var err error
	if data.Id > 0 {
		err = h.DB.Model(&model.Workflow{}).
			Select("name", "icon", "enabled", "score", "summary", "params", "auth_config", "workflow_id", "type").
			Where("id", data.Id).Updates(&app).Error
	} else {
		err = h.DB.Create(&app).Error
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Remove 删除工作流
func (h *WorkflowHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Where("id", id).Delete(&model.Workflow{})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// Enable 启用/禁用
func (h *WorkflowHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.Workflow{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}

// List 获取工作流列表
func (h *WorkflowHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 10)
	offset := (page - 1) * pageSize

	var total int64
	if err := h.DB.Model(&model.Workflow{}).Count(&total).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var apps []model.Workflow
	if err := h.DB.Order("id DESC").Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var appVos []vo.Workflow
	for _, app := range apps {
		var appVo vo.Workflow
		err := utils.CopyObject(app, &appVo)
		if err != nil {
			continue
		}
		// 解析工作流参数
		if app.Params != "" {
			err = utils.JsonDecode(app.Params, &appVo.Params)
			if err != nil {
				logger.Error(err)
				appVo.Params = []vo.WorkflowParam{}
			}
		}
		// 解析授权配置（根据类型解析为不同的结构体）
		if app.AuthConfig != "" {
			if app.Type == "bailian" {
				err = utils.JsonDecode(app.AuthConfig, &appVo.BailianAuthConfig)
			} else {
				err = utils.JsonDecode(app.AuthConfig, &appVo.AuthConfig)
			}
			if err != nil {
				logger.Error(err)
			}
		}
		appVo.Id = app.Id
		appVo.CreatedAt = app.CreatedAt.Unix()
		appVos = append(appVos, appVo)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, appVos))
}

// BatchRemove 批量删除工作流
func (h *WorkflowHandler) BatchRemove(c *gin.Context) {
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

	res := h.DB.Where("id IN ?", data.Ids).Delete(&model.Workflow{})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// ImportWorkflows 从 Coze 导入工作流列表
func (h *WorkflowHandler) ImportWorkflows(c *gin.Context) {
	var data struct {
		AuthConfig types.CozeApiConfig `json:"auth_config"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	workflows, err := h.cozeService.GetWorkflowList(&data.AuthConfig)
	if err != nil {
		resp.ERROR(c, "获取工作流列表失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, workflows)
}

// ImportResult 导入结果结构体
type ImportResult struct {
	WorkflowName string `json:"workflow_name"`
	Status       string `json:"status"` // "imported" | "updated" | "failed"
	Error        string `json:"error,omitempty"`
}

// BatchImportWorkflows 批量导入工作流
func (h *WorkflowHandler) BatchImportWorkflows(c *gin.Context) {
	var data struct {
		Workflows []struct {
			WorkflowID   string `json:"workflow_id"`
			WorkflowName string `json:"workflow_name"`
			Description  string `json:"description"`
			IconURL      string `json:"icon_url"`
		} `json:"workflows"`
		AuthConfig types.CozeApiConfig `json:"auth_config"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Workflows) == 0 {
		resp.ERROR(c, "请选择要导入的工作流")
		return
	}

	// 验证授权配置
	if data.AuthConfig.AppId == "" || data.AuthConfig.PrivateKey == "" {
		resp.ERROR(c, "请配置授权信息")
		return
	}

	var imported int
	var updated int
	var results []ImportResult

	// 批量导入工作流
	for _, workflow := range data.Workflows {
		// 检查是否已存在
		var existingWorkflow model.Workflow
		err := h.DB.Where("workflow_id = ?", workflow.WorkflowID).First(&existingWorkflow).Error
		if err == nil && existingWorkflow.Id > 0 {
			// 已存在，只更新授权配置
			existingWorkflow.AuthConfig = utils.JsonEncode(data.AuthConfig)
			if err := h.DB.Save(&existingWorkflow).Error; err != nil {
				logger.Error("更新工作流授权配置失败: ", err)
				results = append(results, ImportResult{
					WorkflowName: workflow.WorkflowName,
					Status:       "failed",
					Error:        "更新工作流授权配置失败: " + err.Error(),
				})
				continue
			}
			results = append(results, ImportResult{
				WorkflowName: workflow.WorkflowName,
				Status:       "updated",
			})
			updated++
			continue
		}

		// 不存在，获取详情并创建新记录
		detail, err := h.cozeService.GetWorkflowDetail(workflow.WorkflowID, &data.AuthConfig)
		if err != nil {
			logger.Error("获取工作流详情失败: ", err)
			results = append(results, ImportResult{
				WorkflowName: workflow.WorkflowName,
				Status:       "failed",
				Error:        "获取工作流详情失败: " + err.Error(),
			})
			continue
		}

		// 转换参数格式
		params := service.ConvertCozeParamsToWorkflowParams(detail.InputSchema)

		// 下载工作流图标
		iconURL := workflow.IconURL
		if iconURL == "" {
			iconURL = detail.IconURL
		}

		var icon string
		if iconURL != "" {
			uploadHandler := h.uploader.GetUploadHandler()
			uploadedIcon, err := uploadHandler.PutUrlFile(iconURL, false)
			if err != nil {
				logger.Error("下载工作流图标失败: ", err)
				// 图标下载失败不影响导入，使用默认图标
				icon = "/images/app-placeholder.png"
			} else {
				icon = uploadedIcon
			}
		} else {
			icon = "/images/app-placeholder.png"
		}

		// 创建新工作流
		newWorkflow := model.Workflow{
			Name:       workflow.WorkflowName,
			WorkflowId: workflow.WorkflowID,
			Icon:       icon,
			Enabled:    true,
			Score:      1,
			Summary:    workflow.Description,
			Params:     utils.JsonEncode(params),
			AuthConfig: utils.JsonEncode(data.AuthConfig),
		}

		if err := h.DB.Create(&newWorkflow).Error; err != nil {
			logger.Error("创建工作流失败: ", err)
			results = append(results, ImportResult{
				WorkflowName: workflow.WorkflowName,
				Status:       "failed",
				Error:        "创建工作流失败: " + err.Error(),
			})
			continue
		}
		results = append(results, ImportResult{
			WorkflowName: workflow.WorkflowName,
			Status:       "imported",
		})
		imported++
	}

	// 统计失败数量
	failedCount := 0
	for _, result := range results {
		if result.Status == "failed" {
			failedCount++
		}
	}

	// 返回导入结果
	result := gin.H{
		"results": results,
		"summary": gin.H{
			"imported": imported,
			"updated":  updated,
			"failed":   failedCount,
			"total":    len(data.Workflows),
		},
	}

	resp.SUCCESS(c, result)
}
