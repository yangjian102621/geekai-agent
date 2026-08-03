package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowHandler struct {
	BaseHandler
	cozeService     *service.CozeService
	uploadManager   *oss.UploaderManager
	workflowService *service.WorkflowService
}

var (
	WorkflowEventMessage = "message"
	WorkflowEventTool    = "tool"
	WorkflowEventDone    = "done"
	WorkflowEventError   = "error"
)

func NewWorkflowHandler(app *core.AppServer,
	db *gorm.DB,
	cozeService *service.CozeService,
	uploadManager *oss.UploaderManager,
	workflowSevice *service.WorkflowService,
) *WorkflowHandler {
	return &WorkflowHandler{
		BaseHandler:     BaseHandler{App: app, DB: db},
		cozeService:     cozeService,
		uploadManager:   uploadManager,
		workflowService: workflowSevice,
	}
}

func (h *WorkflowHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/workflow/")
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("run-workerflow", h.RunWorkerflow)
		rg.GET("list", h.ListUserWorkflows)
		rg.POST("tasks", h.CreateWorkflowTask)
		rg.GET("tasks", h.ListWorkflowTasks)
		rg.GET("tasks/:task_id", h.GetWorkflowTask)
		rg.POST("tasks/:task_id/cancel", h.CancelWorkflowTask)
		rg.POST("tasks/:task_id/retry", h.RetryWorkflowTask)
	}
}

// ListUserWorkflows 获取可用工作流
func (h *WorkflowHandler) ListUserWorkflows(c *gin.Context) {
	workflows, err := h.workflowService.ListEnabledWorkflows(c.Request.Context())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, workflows)
}

// CreateWorkflowTask 创建任务
func (h *WorkflowHandler) CreateWorkflowTask(c *gin.Context) {
	var data struct {
		WorkflowId string         `json:"workflow_id"`
		Params     map[string]any `json:"params"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 验证参数
	if err := h.validateWorkflowParams(c.Request.Context(), data.WorkflowId, data.Params); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	task, err := h.workflowService.CreateTask(c.Request.Context(), h.GetLoginUserId(c), data.WorkflowId, data.Params)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	taskVo, err := h.workflowService.GetTask(c.Request.Context(), h.GetLoginUserId(c), task.TaskId)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, taskVo)
}

// validateWorkflowParams 验证工作流参数
func (h *WorkflowHandler) validateWorkflowParams(ctx context.Context, workflowId string, params map[string]any) error {
	// 获取工作流定义
	var workflow model.Workflow
	if err := h.DB.WithContext(ctx).Where("workflow_id = ? AND enabled = 1", workflowId).First(&workflow).Error; err != nil {
		return errors.New("工作流不存在或已禁用")
	}

	// 解析工作流参数定义
	var paramDefs []vo.WorkflowParam
	if workflow.Params != "" {
		if err := utils.JsonDecode(workflow.Params, &paramDefs); err != nil {
			return errors.New("工作流参数配置错误")
		}
	}

	// 验证每个参数
	for _, paramDef := range paramDefs {
		value, exists := params[paramDef.Name]

		// 必填验证
		if paramDef.Required && (!exists || value == nil || value == "") {
			return fmt.Errorf("参数 %s (%s) 为必填项", paramDef.Name, paramDef.Label)
		}

		// 如果参数不存在或为空，跳过类型验证
		if !exists || value == nil || value == "" {
			continue
		}

		// 类型验证
		switch paramDef.Type {
		case "Number":
			if _, ok := value.(float64); !ok {
				if _, ok := value.(int); !ok {
					return fmt.Errorf("参数 %s 必须是数字类型", paramDef.Name)
				}
			}
		case "Boolean":
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("参数 %s 必须是布尔类型", paramDef.Name)
			}
		case "CheckBox":
			if _, ok := value.([]any); !ok {
				return fmt.Errorf("参数 %s 必须是数组类型", paramDef.Name)
			}
		case "Select", "Radio":
			// 验证值是否在选项列表中
			if len(paramDef.Options) > 0 {
				valueStr := fmt.Sprintf("%v", value)
				found := false
				for _, opt := range paramDef.Options {
					if fmt.Sprintf("%v", opt) == valueStr {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("参数 %s 的值不在允许的选项列表中", paramDef.Name)
				}
			}
		case "Image", "Audio", "Video", "File", "Doc", "Zip":
			// 验证文件 URL 格式
			urlStr, ok := value.(string)
			if !ok {
				return fmt.Errorf("参数 %s 必须是字符串类型（文件URL）", paramDef.Name)
			}
			if urlStr != "" && !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
				return fmt.Errorf("参数 %s 必须是有效的文件URL", paramDef.Name)
			}
		}
	}

	return nil
}

// ListWorkflowTasks 任务列表
func (h *WorkflowHandler) ListWorkflowTasks(c *gin.Context) {
	status := h.GetTrim(c, "status")
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)

	pageObj, err := h.workflowService.ListTasks(c.Request.Context(), h.GetLoginUserId(c), status, page, pageSize)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, pageObj)
}

// GetWorkflowTask 任务详情
func (h *WorkflowHandler) GetWorkflowTask(c *gin.Context) {
	taskId := c.Param("task_id")
	task, err := h.workflowService.GetTask(c.Request.Context(), h.GetLoginUserId(c), taskId)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, task)
}

func (h *WorkflowHandler) CancelWorkflowTask(c *gin.Context) {
	taskId := c.Param("task_id")
	if err := h.workflowService.CancelTask(c.Request.Context(), h.GetLoginUserId(c), taskId); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *WorkflowHandler) RetryWorkflowTask(c *gin.Context) {
	taskId := c.Param("task_id")
	if err := h.workflowService.RetryTask(c.Request.Context(), h.GetLoginUserId(c), taskId); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// RunWorkerflow 运行工作流（已废弃，使用异步执行）
// 注意：此方法已废弃，工作流现在使用异步执行方式
// 如需实时推送，请使用任务状态查询接口配合前端轮询
func (h *WorkflowHandler) RunWorkerflow(c *gin.Context) {
	resp.ERROR(c, "此接口已废弃，工作流现在使用异步执行方式，请使用任务创建接口")
}

// CozeUploadFile coze 上传文件
func (h *WorkflowHandler) CozeUploadFile(c *gin.Context) {
	fileInfo, err := h.cozeService.UploadFile(c, nil)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, fileInfo)
}
