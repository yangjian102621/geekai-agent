package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/store/vo"
	"geekai/utils"
	"io"
	"net/http"
	"time"

	"github.com/coze-dev/coze-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type CozeService struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewCozeService(db *gorm.DB, rds *redis.Client) *CozeService {
	return &CozeService{db: db, rds: rds}
}

// GetAccessToken 获取 access token
func (s *CozeService) GetAccessToken(config *types.CozeApiConfig) (string, error) {

	if config == nil {
		return "", fmt.Errorf("config is nil")
	}

	// 先去 redis 中获取
	key := fmt.Sprintf("access_token_%s", config.AppId)

	accessToken, err := s.rds.Get(context.Background(), key).Result()
	if err == nil {
		return accessToken, nil
	}

	oauth, err := coze.NewJWTOAuthClient(coze.NewJWTOAuthClientParam{
		ClientID: config.AppId, PublicKey: config.PublicKeyID, PrivateKeyPEM: config.PrivateKey,
	}, coze.WithAuthBaseURL(config.ApiUrl))
	if err != nil {
		return "", fmt.Errorf("error creating JWT OAuth client: %v", err)
	}

	expiresIn := 900 // token 有效期 15 分钟
	resp, err := oauth.GetAccessToken(context.Background(), &coze.GetJWTAccessTokenReq{
		TTL: expiresIn,
	})
	if err != nil {
		return "", fmt.Errorf("error getting access token: %v", err)
	}
	// 缓存 token
	s.rds.Set(context.Background(), key, resp.AccessToken, time.Duration(expiresIn)*time.Second)
	return resp.AccessToken, nil
}

// CozeAgent 表示 Coze 平台上的智能体
type CozeAgent struct {
	BotID       string `json:"bot_id"`       // 智能体ID
	BotName     string `json:"bot_name"`     // 智能体名称
	Description string `json:"description"`  // 智能体描述
	IconURL     string `json:"icon_url"`     // 智能体图标URL
	PublishTime string `json:"publish_time"` // 发布时间
}

// GetAgentList 获取智能体列表
func (s *CozeService) GetAgentList(config *types.CozeApiConfig) ([]CozeAgent, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	accessToken, err := s.GetAccessToken(config)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %v", err)
	}

	authCli := coze.NewTokenAuth(accessToken)
	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(config.ApiUrl))
	ctx := context.Background()
	botList, err := cozeCli.Bots.List(ctx, &coze.ListBotsReq{
		SpaceID:  config.SpaceId,
		PageNum:  1,
		PageSize: 100,
	})

	if err != nil {
		return nil, fmt.Errorf("error getting bot list: %v", err)
	}

	if botList.Err() != nil {
		return nil, fmt.Errorf("error listing bots: %v", botList.Err())
	}

	var agents []CozeAgent
	for botList.Next() {
		bot := botList.Current()
		if bot != nil {
			agents = append(agents, CozeAgent{
				BotID:       bot.BotID,
				BotName:     bot.BotName,
				Description: bot.Description,
				IconURL:     bot.IconURL,
				PublishTime: bot.PublishTime,
			})
		}
	}

	return agents, nil
}

// UploadFile 上传文件
func (s *CozeService) UploadFile(c *gin.Context, config *types.CozeApiConfig) (*coze.FileInfo, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	token, err := s.GetAccessToken(config)
	if err != nil {
		return nil, err
	}
	authCli := coze.NewTokenAuth(token)
	// 初始化 Coze 客户端
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(config.ApiUrl))
	ctx := context.Background()
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 上传文件
	uploadResp, err := cozeCli.Files.Upload(ctx, &coze.UploadFilesReq{
		File: coze.NewUploadFile(src, fileHeader.Filename),
	})
	if err != nil {
		return nil, err
	}

	return &uploadResp.FileInfo, nil
}

// CozeWorkflow 表示 Coze 平台上的工作流
type CozeWorkflow struct {
	WorkflowID   string `json:"workflow_id"`   // 工作流ID
	WorkflowName string `json:"workflow_name"` // 工作流名称
	Description  string `json:"description"`   // 工作流描述
	IconURL      string `json:"icon_url"`      // 工作流图标URL
}

// CozeWorkflowDetail 表示 Coze 工作流详情
type CozeWorkflowDetail struct {
	WorkflowID   string              `json:"workflow_id"`   // 工作流ID
	WorkflowName string              `json:"workflow_name"` // 工作流名称
	Description  string              `json:"description"`   // 工作流描述
	IconURL      string              `json:"icon_url"`      // 工作流图标URL
	InputSchema  []CozeWorkflowParam `json:"input_schema"`  // 输入参数
}

// CozeWorkflowParam 表示 Coze 工作流参数
type CozeWorkflowParam struct {
	Name        string `json:"name"`        // 参数名称
	Type        string `json:"type"`        // 参数类型
	Description string `json:"description"` // 参数说明
	Required    bool   `json:"required"`    // 是否必填
	Default     any    `json:"default"`     // 默认值
}

// GetWorkflowList 获取工作流列表
func (s *CozeService) GetWorkflowList(config *types.CozeApiConfig) ([]CozeWorkflow, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	token, err := s.GetAccessToken(config)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %v", err)
	}

	authCli := coze.NewTokenAuth(token)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(config.ApiUrl))
	ctx := context.Background()

	// 准备请求参数
	publishStatus := coze.PublishStatusPublishedOnline
	req := &coze.ListWorkflowReq{
		WorkspaceID:   &config.SpaceId,
		PublishStatus: &publishStatus,
		PageNum:       1,
		PageSize:      30,
	}

	// List workflows
	listResp, err := cozeCli.Workflows.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error listing workflows: %v", err)
	}

	var workflows []CozeWorkflow

	// 处理第一页数据
	items := listResp.Items()
	for _, item := range items {
		if item != nil {
			workflows = append(workflows, CozeWorkflow{
				WorkflowID:   item.WorkflowID,
				WorkflowName: item.WorkflowName,
				Description:  item.Description,
				IconURL:      item.IconURL,
			})
		}
	}

	// 如果还有更多数据，继续获取（最多获取前几页）
	pageNum := 2
	for listResp.HasMore() && pageNum <= 10 { // 最多获取10页
		req.PageNum = pageNum
		listResp, err = cozeCli.Workflows.List(ctx, req)
		if err != nil {
			// 如果后续分页失败，返回已获取的数据
			break
		}

		items = listResp.Items()
		for _, item := range items {
			if item != nil {
				workflows = append(workflows, CozeWorkflow{
					WorkflowID:   item.WorkflowID,
					WorkflowName: item.WorkflowName,
					Description:  item.Description,
					IconURL:      item.IconURL,
				})
			}
		}

		pageNum++
	}

	return workflows, nil
}

// GetWorkflowDetail 获取工作流详情
// 参考文档: https://www.coze.cn/open/docs/developer_guides/get_workflow_info
func (s *CozeService) GetWorkflowDetail(workflowId string, config *types.CozeApiConfig) (*CozeWorkflowDetail, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if workflowId == "" {
		return nil, fmt.Errorf("workflow_id is empty")
	}

	token, err := s.GetAccessToken(config)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %v", err)
	}

	// 使用 HTTP API 调用 Coze 获取工作流详情接口
	url := fmt.Sprintf("%s/v1/workflows/%s?include_input_output=true", config.ApiUrl, workflowId)

	// 发送 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var apiResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			WorkflowDetail struct {
				WorkflowID   string `json:"workflow_id"`
				WorkflowName string `json:"workflow_name"`
				Description  string `json:"description"`
				IconURL      string `json:"icon_url"`
				AppID        string `json:"app_id"`
			} `json:"workflow_detail"`
			Input struct {
				Parameters map[string]struct {
					Required     bool   `json:"required"`
					Description  string `json:"description"`
					DefaultValue any    `json:"default_value"`
					Type         string `json:"type"`
					Items        *struct {
						Type string `json:"type"`
					} `json:"items,omitempty"` // for array type
				} `json:"parameters"`
			} `json:"input"`
			Output struct {
				Parameters    map[string]any `json:"parameters"`
				TerminatePlan string         `json:"terminate_plan"`
			} `json:"output"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", apiResp.Msg)
	}

	logger.Debugf("workflow detail: %s", utils.JsonEncode(apiResp.Data))

	// 转换 input.parameters 为 InputSchema 格式
	var inputSchema []CozeWorkflowParam
	for name, param := range apiResp.Data.Input.Parameters {
		inputSchema = append(inputSchema, CozeWorkflowParam{
			Name:        name,
			Type:        param.Type,
			Description: param.Description,
			Required:    param.Required,
			Default:     param.DefaultValue,
		})
	}

	detail := &CozeWorkflowDetail{
		WorkflowID:   apiResp.Data.WorkflowDetail.WorkflowID,
		WorkflowName: apiResp.Data.WorkflowDetail.WorkflowName,
		Description:  apiResp.Data.WorkflowDetail.Description,
		IconURL:      apiResp.Data.WorkflowDetail.IconURL,
		InputSchema:  inputSchema,
	}

	return detail, nil
}

// ConvertCozeParamsToWorkflowParams 将 Coze 参数格式转换为系统参数格式
func ConvertCozeParamsToWorkflowParams(cozeParams []CozeWorkflowParam) []vo.WorkflowParam {
	var params []vo.WorkflowParam
	for _, p := range cozeParams {
		if p.Name == "CONVERSATION_NAME" || p.Name == "USER_INPUT" || p.Name == "BOT_USER_INPUT" {
			continue
		}
		param := vo.WorkflowParam{
			Name:     p.Name,
			Label:    p.Description,
			Type:     cases.Title(language.English).String(p.Type),
			Required: p.Required,
			Default:  p.Default,
		}
		params = append(params, param)
	}
	return params
}

// RunWorkflowAsync 异步执行工作流
// 参考文档: https://www.coze.cn/open/docs/developer_guides/workflow_run
func (s *CozeService) RunWorkflowAsync(workflowId string, params map[string]any, config *types.CozeApiConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is nil")
	}

	token, err := s.GetAccessToken(config)
	if err != nil {
		return "", fmt.Errorf("error getting access token: %v", err)
	}

	authCli := coze.NewTokenAuth(token)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(config.ApiUrl))
	ctx := context.Background()

	req := &coze.RunWorkflowsReq{
		WorkflowID: workflowId,
		Parameters: params,
		IsAsync:    true, // 异步执行
	}

	// 使用 SDK 的 Create 方法创建异步工作流运行
	resp, err := cozeCli.Workflows.Runs.Create(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error creating workflow run: %v", err)
	}

	if resp.ExecuteID == "" {
		return "", fmt.Errorf("workflow execute ID is empty")
	}

	return resp.ExecuteID, nil
}

// WorkflowRunStatus 工作流执行状态
type WorkflowRunStatus struct {
	RunID      string           `json:"run_id"`      // 运行ID
	Status     string           `json:"status"`      // 状态：pending, running, completed, failed
	Progress   int              `json:"progress"`    // 进度 0-100
	Output     map[string]any   `json:"output"`      // 输出结果
	Error      string           `json:"error"`       // 错误信息
	CreatedAt  int64            `json:"created_at"`  // 创建时间
	FinishedAt int64            `json:"finished_at"` // 完成时间
	Events     []map[string]any `json:"events"`      // 执行事件列表
}

// GetWorkflowRunStatus 查询工作流执行状态
// 参考文档: https://www.coze.cn/open/docs/developer_guides/workflow_history
func (s *CozeService) GetWorkflowRunStatus(executeId string, workflowId string, config *types.CozeApiConfig) (*WorkflowRunStatus, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if executeId == "" {
		return nil, fmt.Errorf("execute_id is empty")
	}

	if workflowId == "" {
		return nil, fmt.Errorf("workflow_id is empty")
	}

	token, err := s.GetAccessToken(config)
	if err != nil {
		return nil, fmt.Errorf("error getting access token: %v", err)
	}

	authCli := coze.NewTokenAuth(token)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(config.ApiUrl))
	ctx := context.Background()

	// 使用 SDK 的 Histories.Retrieve 方法查询工作流运行历史
	historyResp, err := cozeCli.Workflows.Runs.Histories.Retrieve(ctx, &coze.RetrieveWorkflowsRunsHistoriesReq{
		WorkflowID: workflowId,
		ExecuteID:  executeId,
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving workflow run history: %v", err)
	}

	if len(historyResp.Histories) == 0 {
		return nil, fmt.Errorf("workflow run history not found")
	}

	history := historyResp.Histories[0]

	status := &WorkflowRunStatus{
		RunID:      executeId,
		Status:     string(history.ExecuteStatus),
		Output:     make(map[string]any),
		Events:     []map[string]any{},
		CreatedAt:  int64(history.CreateTime),
		FinishedAt: int64(history.UpdateTime),
		Error:      history.ErrorMessage,
	}

	// 根据状态设置进度
	switch history.ExecuteStatus {
	case coze.WorkflowExecuteStatusSuccess:
		status.Progress = 100
		// 解析输出结果
		if history.Output != "" {
			var outputMap map[string]any
			if err := json.Unmarshal([]byte(history.Output), &outputMap); err == nil {
				status.Output = outputMap
			} else {
				// 如果不是 JSON，直接作为字符串存储
				status.Output = map[string]any{"raw": history.Output}
			}
		}
	case coze.WorkflowExecuteStatusRunning:
		status.Progress = 50
	case coze.WorkflowExecuteStatusFail:
		status.Progress = 0
	default:
		status.Progress = 0
	}

	return status, nil
}
