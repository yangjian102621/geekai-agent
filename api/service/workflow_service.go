package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"

	"gorm.io/gorm"
)

// WorkflowService 面向前台的工作流和任务管理
type WorkflowService struct {
	db             *gorm.DB
	cozeService    *CozeService
	bailianService *BailianService
	userService    *UserService
}

func NewWorkflowService(db *gorm.DB, cozeService *CozeService, bailianService *BailianService, userService *UserService) *WorkflowService {
	svc := &WorkflowService{db: db, cozeService: cozeService, bailianService: bailianService, userService: userService}
	return svc
}

// ListEnabledWorkflows 获取可用工作流
func (s *WorkflowService) ListEnabledWorkflows(ctx context.Context) ([]vo.Workflow, error) {
	var workflows []model.Workflow
	err := s.db.WithContext(ctx).
		Where("enabled = ?", true).
		Order("id DESC").
		Find(&workflows).Error
	if err != nil {
		return nil, err
	}
	list := make([]vo.Workflow, 0, len(workflows))
	for _, item := range workflows {
		var params []vo.WorkflowParam
		if item.Params != "" {
			_ = utils.JsonDecode(item.Params, &params)
		}
		wfVo := vo.Workflow{
			Id:         item.Id,
			Name:       item.Name,
			Icon:       item.Icon,
			Enabled:    item.Enabled,
			Score:      item.Score,
			Summary:    item.Summary,
			Params:     params,
			Type:       item.Type,
			WorkflowId: item.WorkflowId,
			CreatedAt:  item.CreatedAt.Unix(),
			UpdatedAt:  item.UpdatedAt.Unix(),
			LastRunAt:  item.LastRunAt,
		}
		// 根据类型解析不同的授权配置
		if item.AuthConfig != "" {
			if item.Type == "bailian" {
				_ = utils.JsonDecode(item.AuthConfig, &wfVo.BailianAuthConfig)
			} else {
				_ = utils.JsonDecode(item.AuthConfig, &wfVo.AuthConfig)
			}
		}
		list = append(list, wfVo)
	}
	return list, nil
}

// CreateTask 创建异步任务
func (s *WorkflowService) CreateTask(ctx context.Context, userId uint, workflowId string, params map[string]any) (*model.WorkflowTask, error) {
	if workflowId == "" {
		return nil, errors.New("workflow_id 不能为空")
	}
	var workflow model.Workflow
	if err := s.db.WithContext(ctx).Where("workflow_id = ? AND enabled = 1", workflowId).First(&workflow).Error; err != nil {
		return nil, err
	}

	// 格式化参数（确保类型正确）
	formattedParams := s.formatWorkflowParams(ctx, workflowId, params)

	// 根据工作流类型执行不同的工作流
	var executeId string
	var err error

	if workflow.Type == "bailian" {
		// 百炼工作流：同步执行，直接获取结果
		var bailianAuthConfig types.BailianApiConfig
		if workflow.AuthConfig == "" {
			return nil, errors.New("工作流未配置授权信息")
		}
		if err := utils.JsonDecode(workflow.AuthConfig, &bailianAuthConfig); err != nil {
			return nil, fmt.Errorf("授权配置解析失败: %v", err)
		}

		result, runErr := s.bailianService.RunWorkflowSync(workflow.WorkflowId, formattedParams, &bailianAuthConfig)
		if runErr != nil {
			return nil, fmt.Errorf("执行百炼工作流失败: %v", runErr)
		}
		executeId = result.RequestId

		// 百炼工作流同步执行完成，直接创建已完成的任务
		task := &model.WorkflowTask{
			UserId:       userId,
			TaskId:       utils.RandString(18),
			WorkflowId:   workflow.WorkflowId,
			WorkflowName: workflow.Name,
			Status:       vo.WorkflowTaskStatusCompleted,
			Progress:     100,
			Params:       utils.JsonEncode(formattedParams),
			Score:        workflow.Score,
			Output: utils.JsonEncode(map[string]any{
				"execute_id": executeId,
				"result":     result.Output,
			}),
		}
		if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
			return nil, err
		}

		// 扣减用户积分（如果工作流需要消耗积分）
		if workflow.Score > 0 {
			if err := s.userService.DecreaseScores(userId, workflow.Score, model.ScoreLog{
				Type:    types.ScoreConsume,
				Subject: workflow.Name,
				Remark:  fmt.Sprintf("执行工作流任务，任务ID：%s", task.TaskId),
			}); err != nil {
				logger.Errorf("扣减用户积分失败: %v", err)
			}
		}

		return task, nil
	}

	// Coze 工作流：异步执行
	var authConfig types.CozeApiConfig
	if workflow.AuthConfig == "" {
		return nil, errors.New("工作流未配置授权信息")
	}
	if err := utils.JsonDecode(workflow.AuthConfig, &authConfig); err != nil {
		return nil, fmt.Errorf("授权配置解析失败: %v", err)
	}

	// 异步执行工作流
	runId, err := s.cozeService.RunWorkflowAsync(workflow.WorkflowId, formattedParams, &authConfig)
	if err != nil {
		return nil, fmt.Errorf("启动工作流失败: %v", err)
	}
	executeId = runId

	task := &model.WorkflowTask{
		UserId:       userId,
		TaskId:       utils.RandString(18),
		WorkflowId:   workflow.WorkflowId,
		WorkflowName: workflow.Name,
		Status:       vo.WorkflowTaskStatusPending,
		Params:       utils.JsonEncode(formattedParams),
		Score:        workflow.Score,
		// 将 execute_id 存储在 Output 中
		Output: utils.JsonEncode(map[string]any{
			"execute_id": executeId,
		}),
	}
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}

	// 扣减用户积分（如果工作流需要消耗积分）
	if workflow.Score > 0 {
		if err := s.userService.DecreaseScores(userId, workflow.Score, model.ScoreLog{
			Type:    types.ScoreConsume,
			Subject: workflow.Name,
			Remark:  fmt.Sprintf("执行工作流任务，任务ID：%s", task.TaskId),
		}); err != nil {
			logger.Errorf("扣减用户积分失败: %v", err)
			// 积分扣减失败不影响任务创建，但记录错误日志
		}
	}

	return task, nil
}

// formatWorkflowParams 格式化工作流参数，确保类型正确
func (s *WorkflowService) formatWorkflowParams(ctx context.Context, workflowId string, params map[string]any) map[string]any {
	// 获取工作流定义以了解参数类型
	var workflow model.Workflow
	if err := s.db.WithContext(ctx).Where("workflow_id = ?", workflowId).First(&workflow).Error; err != nil {
		return params // 如果获取失败，返回原始参数
	}

	var paramDefs []vo.WorkflowParam
	if workflow.Params != "" {
		_ = utils.JsonDecode(workflow.Params, &paramDefs)
	}

	formatted := make(map[string]any)
	paramDefMap := make(map[string]vo.WorkflowParam)
	for _, def := range paramDefs {
		paramDefMap[def.Name] = def
	}

	// 格式化每个参数
	for key, value := range params {
		if def, ok := paramDefMap[key]; ok {
			// 根据类型格式化值
			switch def.Type {
			case "Number":
				// 确保是数字类型
				if num, ok := value.(float64); ok {
					formatted[key] = num
				} else if num, ok := value.(int); ok {
					formatted[key] = float64(num)
				} else {
					formatted[key] = value
				}
			case "Boolean":
				// 确保是布尔类型
				if b, ok := value.(bool); ok {
					formatted[key] = b
				} else {
					formatted[key] = value
				}
			case "CheckBox":
				// 确保是数组类型
				if arr, ok := value.([]any); ok {
					formatted[key] = arr
				} else {
					formatted[key] = value
				}
			case "Image", "Audio", "Video", "File", "Doc", "Zip":
				// 文件类型保持为字符串 URL
				if str, ok := value.(string); ok {
					formatted[key] = str
				} else {
					formatted[key] = value
				}
			default:
				// 其他类型保持原样
				formatted[key] = value
			}
		} else {
			// 未定义的参数也保留
			formatted[key] = value
		}
	}

	return formatted
}

// ListTasks 查询任务列表
func (s *WorkflowService) ListTasks(ctx context.Context, userId uint, status string, page, pageSize int) (*vo.Page, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := s.db.WithContext(ctx).Model(&model.WorkflowTask{}).Where("user_id = ?", userId)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var tasks []model.WorkflowTask
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, err
	}

	respTasks := make([]vo.WorkflowTask, 0, len(tasks))
	for _, task := range tasks {
		respTasks = append(respTasks, s.buildTaskVo(task))
	}

	pageObj := vo.NewPage(total, page, pageSize, respTasks)
	return &pageObj, nil
}

// GetTask 查询任务详情
func (s *WorkflowService) GetTask(ctx context.Context, userId uint, taskId string) (*vo.WorkflowTask, error) {
	var task model.WorkflowTask
	if err := s.db.WithContext(ctx).Where("task_id = ? AND user_id = ?", taskId, userId).First(&task).Error; err != nil {
		return nil, err
	}

	voTask := s.buildTaskVo(task)
	return &voTask, nil
}

// CancelTask 取消任务
func (s *WorkflowService) CancelTask(ctx context.Context, userId uint, taskId string) error {
	res := s.db.WithContext(ctx).Model(&model.WorkflowTask{}).
		Where("task_id = ? AND user_id = ? AND status IN ?", taskId, userId, []vo.WorkflowTaskStatus{vo.WorkflowTaskStatusPending, vo.WorkflowTaskStatusRunning}).
		Updates(map[string]any{
			"status":   vo.WorkflowTaskStatusCanceled,
			"progress": 0,
		})
	return res.Error
}

// RetryTask 重试任务
func (s *WorkflowService) RetryTask(ctx context.Context, userId uint, taskId string) error {
	var task model.WorkflowTask
	if err := s.db.WithContext(ctx).Where("task_id = ? AND user_id = ?", taskId, userId).First(&task).Error; err != nil {
		return err
	}

	// 解析授权配置
	var workflow model.Workflow
	if err := s.db.WithContext(ctx).Where("workflow_id = ?", task.WorkflowId).First(&workflow).Error; err != nil {
		return fmt.Errorf("工作流不存在: %v", err)
	}

	// 解析任务参数
	var params map[string]any
	if task.Params != "" {
		if err := utils.JsonDecode(task.Params, &params); err != nil {
			return fmt.Errorf("参数解析失败: %v", err)
		}
	}

	var executeId string
	var err error

	if workflow.Type == "bailian" {
		// 百炼工作流：同步执行
		var bailianAuthConfig types.BailianApiConfig
		if workflow.AuthConfig == "" {
			return errors.New("工作流未配置授权信息")
		}
		if err := utils.JsonDecode(workflow.AuthConfig, &bailianAuthConfig); err != nil {
			return fmt.Errorf("授权配置解析失败: %v", err)
		}

		result, runErr := s.bailianService.RunWorkflowSync(task.WorkflowId, params, &bailianAuthConfig)
		if runErr != nil {
			return fmt.Errorf("执行百炼工作流失败: %v", runErr)
		}
		executeId = result.RequestId

		// 保存原始状态，用于判断是否需要重新扣费
		originalStatus := task.Status

		task.WorkflowName = workflow.Name
		task.Status = vo.WorkflowTaskStatusCompleted
		task.Progress = 100
		task.Error = ""
		task.Output = utils.JsonEncode(map[string]any{
			"execute_id": executeId,
			"result":     result.Output,
		})
		task.UpdatedAt = time.Now()

		if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
			return err
		}

		// 如果任务之前是失败状态（积分已被退回），重试时需要再次扣费
		if task.Score > 0 && originalStatus == vo.WorkflowTaskStatusFailed {
			if err := s.userService.DecreaseScores(userId, task.Score, model.ScoreLog{
				Type:    types.ScoreConsume,
				Subject: task.WorkflowName,
				Remark:  fmt.Sprintf("重试工作流任务，任务ID：%s", task.TaskId),
			}); err != nil {
				logger.Errorf("重试任务扣减用户积分失败: %v", err)
			}
		}

		return nil
	}

	// Coze 工作流：异步执行
	var authConfig types.CozeApiConfig
	if workflow.AuthConfig == "" {
		return errors.New("工作流未配置授权信息")
	}
	if err := utils.JsonDecode(workflow.AuthConfig, &authConfig); err != nil {
		return fmt.Errorf("授权配置解析失败: %v", err)
	}

	// 异步执行工作流
	runId, err := s.cozeService.RunWorkflowAsync(task.WorkflowId, params, &authConfig)
	if err != nil {
		return fmt.Errorf("启动工作流失败: %v", err)
	}
	executeId = runId

	// 保存原始状态，用于判断是否需要重新扣费
	originalStatus := task.Status

	// 更新工作流名称（可能已变更）
	task.WorkflowName = workflow.Name
	task.Status = vo.WorkflowTaskStatusPending
	task.Progress = 0
	task.Error = ""
	task.Output = utils.JsonEncode(map[string]any{
		"execute_id": executeId,
	})
	task.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&task).Error; err != nil {
		return err
	}

	// 如果任务之前是失败状态（积分已被退回），重试时需要再次扣费
	// 只有失败的任务才会退回积分，所以重试失败的任务时需要重新扣费
	if task.Score > 0 && originalStatus == vo.WorkflowTaskStatusFailed {
		if err := s.userService.DecreaseScores(userId, task.Score, model.ScoreLog{
			Type:    types.ScoreConsume,
			Subject: task.WorkflowName,
			Remark:  fmt.Sprintf("重试工作流任务，任务ID：%s", task.TaskId),
		}); err != nil {
			logger.Errorf("重试任务扣减用户积分失败: %v", err)
			// 积分扣减失败不影响任务重试，但记录错误日志
		}
	}

	return nil
}

func (s *WorkflowService) buildTaskVo(task model.WorkflowTask) vo.WorkflowTask {
	result := vo.WorkflowTask{
		Id:           task.Id,
		TaskId:       task.TaskId,
		WorkflowId:   task.WorkflowId,
		WorkflowName: task.WorkflowName,
		Status:       task.Status,
		Progress:     task.Progress,
		Error:        task.Error,
		Score:        task.Score,
		CreatedAt:    task.CreatedAt.Unix(),
		UpdatedAt:    task.UpdatedAt.Unix(),
		Params:       map[string]any{},
		Output:       map[string]any{},
	}

	if task.Params != "" {
		_ = utils.JsonDecode(task.Params, &result.Params)
	}
	if task.Output != "" {
		_ = utils.JsonDecode(task.Output, &result.Output)
	}

	return result
}

// pollWorkflowStatusOnce 轮询一次工作流执行状态
func (s *WorkflowService) pollWorkflowStatusOnce(task model.WorkflowTask) {
	// 百炼工作流是同步执行的，不需要轮询
	var workflow model.Workflow
	if err := s.db.Where("workflow_id = ?", task.WorkflowId).First(&workflow).Error; err != nil {
		logger.Errorf("获取工作流配置失败: %v", err)
		return
	}
	if workflow.Type == "bailian" {
		return
	}

	// 检查任务是否被取消
	if task.Status == vo.WorkflowTaskStatusCanceled {
		return
	}

	// 解析 execute_id
	var output map[string]any
	if task.Output != "" {
		if err := utils.JsonDecode(task.Output, &output); err != nil {
			logger.Errorf("解析任务输出失败: %v", err)
			return
		}
	}

	executeId, ok := output["execute_id"].(string)
	if !ok || executeId == "" {
		// 兼容旧数据，尝试 run_id
		if runId, ok := output["run_id"].(string); ok && runId != "" {
			executeId = runId
		} else {
			logger.Errorf("任务 %s 缺少 execute_id", task.TaskId)
			return
		}
	}

	// 获取工作流配置（已在函数开头获取）
	var authConfig types.CozeApiConfig
	if workflow.AuthConfig == "" {
		logger.Errorf("工作流 %s 未配置授权信息", task.WorkflowId)
		return
	}
	if err := utils.JsonDecode(workflow.AuthConfig, &authConfig); err != nil {
		logger.Errorf("解析授权配置失败: %v", err)
		return
	}

	// 查询工作流执行状态
	status, err := s.cozeService.GetWorkflowRunStatus(executeId, task.WorkflowId, &authConfig)
	if err != nil {
		// 如果查询失败，记录日志但不更新状态（可能是临时网络问题）
		logger.Errorf("查询工作流状态失败: %v", err)
		return
	}

	// 更新任务状态
	update := map[string]any{
		"progress": status.Progress,
	}

	switch status.Status {
	case "Success":
		update["status"] = vo.WorkflowTaskStatusCompleted
		update["output"] = utils.JsonEncode(status.Output)
		_ = s.db.Model(&model.WorkflowTask{}).Where("task_id = ?", task.TaskId).Updates(update).Error
	case "Running":
		update["status"] = vo.WorkflowTaskStatusRunning
		_ = s.db.Model(&model.WorkflowTask{}).Where("task_id = ?", task.TaskId).Updates(update).Error
	case "Fail":
		// 在更新状态前，检查当前任务状态，避免重复退回积分
		var currentTask model.WorkflowTask
		needRefund := false
		if err := s.db.Where("task_id = ?", task.TaskId).First(&currentTask).Error; err == nil {
			// 如果当前状态不是失败，且任务有消耗积分，则需要退回
			if currentTask.Status != vo.WorkflowTaskStatusFailed && currentTask.Score > 0 {
				needRefund = true
			}
		}

		update["status"] = vo.WorkflowTaskStatusFailed
		update["error"] = status.Error
		_ = s.db.Model(&model.WorkflowTask{}).Where("task_id = ?", task.TaskId).Updates(update).Error

		// 如果任务失败且之前不是失败状态，退回用户积分
		if needRefund {
			if err := s.userService.IncreaseScores(int(task.UserId), task.Score, model.ScoreLog{
				Type:    types.ScoreRefund,
				Subject: task.WorkflowName,
				Remark:  fmt.Sprintf("工作流任务执行失败，退回积分，任务ID：%s", task.TaskId),
			}); err != nil {
				logger.Errorf("退回用户积分失败: %v", err)
			} else {
				logger.Infof("任务 %s 执行失败，已退回 %d 积分给用户 %d", task.TaskId, task.Score, task.UserId)
			}
		}
	default:
		// 其他状态（如 pending）保持当前状态
		update["status"] = vo.WorkflowTaskStatusPending
		_ = s.db.Model(&model.WorkflowTask{}).Where("task_id = ?", task.TaskId).Updates(update).Error
	}
}

// StartTaskPolling 启动任务轮询，循环查询未完成的任务并更新状态
func (s *WorkflowService) StartTaskPolling() {
	ticker := time.NewTicker(5 * time.Second) // 每5秒查询一次数据库
	defer ticker.Stop()

	for range ticker.C {
		// 查询所有待处理的任务
		var tasks []model.WorkflowTask
		if err := s.db.Where("status IN ?", []vo.WorkflowTaskStatus{
			vo.WorkflowTaskStatusPending,
			vo.WorkflowTaskStatusRunning,
		}).Find(&tasks).Error; err != nil {
			logger.Errorf("查询待处理任务失败: %v", err)
			continue
		}

		// 逐个轮询任务状态
		for _, task := range tasks {
			s.pollWorkflowStatusOnce(task)
			// 每个任务轮询完成后 sleep 0.5秒，防止 coze 并发限制
			time.Sleep(500 * time.Millisecond)
		}
	}
}
