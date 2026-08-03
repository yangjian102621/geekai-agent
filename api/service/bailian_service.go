package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/service/bailian"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BailianService struct {
	db *gorm.DB
}

func NewBailianService(db *gorm.DB) *BailianService {
	return &BailianService{db: db}
}

// RunWorkflowAsync 执行百炼工作流（同步流式调用，在服务端完成流式读取后存储结果）
func (s *BailianService) RunWorkflowAsync(workflowId string, params map[string]any, authConfig *types.BailianApiConfig) (string, error) {
	if authConfig == nil {
		return "", fmt.Errorf("config is nil")
	}

	if authConfig.ApiKey == "" || authConfig.AppId == "" {
		return "", fmt.Errorf("百炼 API Key 或 App ID 为空")
	}

	client := bailian.NewClient(authConfig.ApiKey, authConfig.AppId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 构建请求参数
	prompt := buildBailianPrompt(params)
	req := &bailian.CompletionRequest{
		Input: bailian.CompletionInput{
			Prompt: prompt,
		},
		Parameters: bailian.CompletionParameters{
			IncrementalOutput: false,
			HasThoughts:       false,
		},
	}

	// 使用流式调用，但在服务端完成读取
	streamCh, err := client.API().CompletionStream(ctx, req)
	if err != nil {
		return "", fmt.Errorf("启动百炼工作流失败: %v", err)
	}

	// 读取完整结果
	var result strings.Builder
	var requestId string
	for stream := range streamCh {
		if stream.Err != nil {
			return requestId, fmt.Errorf("百炼工作流执行错误: %v", stream.Err)
		}
		if stream.Output.Text != "" {
			result.WriteString(stream.Output.Text)
		}
		if stream.RequestId != "" {
			requestId = stream.RequestId
		}
		if stream.Output.FinishReason == "stop" {
			break
		}
	}

	// 将结果存储到 task 的 output 中（通过 requestId 标识）
	// 返回 requestId 作为 execute_id
	return requestId, nil
}

// BailianWorkflowResult 百炼工作流执行结果
type BailianWorkflowResult struct {
	RequestId string `json:"request_id"`
	Output    string `json:"output"`
	Status    string `json:"status"` // completed, failed
}

// RunWorkflowSync 同步执行百炼工作流，返回完整结果
func (s *BailianService) RunWorkflowSync(workflowId string, params map[string]any, authConfig *types.BailianApiConfig) (*BailianWorkflowResult, error) {
	if authConfig == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if authConfig.ApiKey == "" || authConfig.AppId == "" {
		return nil, fmt.Errorf("百炼 API Key 或 App ID 为空")
	}

	client := bailian.NewClient(authConfig.ApiKey, authConfig.AppId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	prompt := buildBailianPrompt(params)
	req := &bailian.CompletionRequest{
		Input: bailian.CompletionInput{
			Prompt: prompt,
		},
		Parameters: bailian.CompletionParameters{
			IncrementalOutput: false,
			HasThoughts:       false,
		},
	}

	streamCh, err := client.API().CompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("启动百炼工作流失败: %v", err)
	}

	var result strings.Builder
	var requestId string
	for stream := range streamCh {
		if stream.Err != nil {
			return &BailianWorkflowResult{
				RequestId: requestId,
				Status:    "failed",
			}, stream.Err
		}
		if stream.Output.Text != "" {
			result.WriteString(stream.Output.Text)
		}
		if stream.RequestId != "" {
			requestId = stream.RequestId
		}
		if stream.Output.FinishReason == "stop" {
			break
		}
	}

	return &BailianWorkflowResult{
		RequestId: requestId,
		Output:    result.String(),
		Status:    "completed",
	}, nil
}

// buildBailianPrompt 将参数构建为百炼工作流的 prompt
func buildBailianPrompt(params map[string]any) string {
	if len(params) == 0 {
		return ""
	}

	// 尝试将参数序列化为 JSON 作为 prompt
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Sprintf("%v", params)
	}
	return string(data)
}

// GetWorkflowRunStatus 获取百炼工作流执行状态
// 百炼没有独立的状态查询 API，任务结果在 RunWorkflowSync 中已直接存储
// 此方法从数据库中读取任务结果
func (s *BailianService) GetWorkflowRunStatus(taskId string, authConfig *types.BailianApiConfig) (*WorkflowRunStatus, error) {
	// 百炼工作流是同步执行的，结果已在 RunWorkflowAsync 中存储
	// 这里返回 nil 表示不需要轮询
	return nil, nil
}

// ParseWorkflowOutput 解析百炼工作流输出
func ParseWorkflowOutput(output string) map[string]any {
	if output == "" {
		return map[string]any{}
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(output), &result); err == nil {
		return result
	}

	// 如果不是 JSON，作为原始文本存储
	return map[string]any{"raw": output}
}

// shouldPollBailian 百炼工作流是否需要轮询（不需要，因为是同步执行）
func shouldPollBailian(workflowType string) bool {
	return !slices.Contains([]string{"bailian"}, workflowType)
}
