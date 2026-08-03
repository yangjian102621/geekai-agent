package vo

import "geekai/core/types"

type AppType string

const (
	AppTypeBot      AppType = "bot"      // 智能体
	AppTypeWorkflow AppType = "workflow" // 工作流
)

type WorkflowParam struct {
	Label       string   `json:"label,omitempty"`        // 参数说明
	Name        string   `json:"name,omitempty"`         // 参数名称
	Type        string   `json:"type,omitempty"`         // 参数类型
	Default     any      `json:"default,omitempty"`      // 默认值
	Required    bool     `json:"required,omitempty"`     // 是否必填
	MaxFilesize int      `json:"max_filesize,omitempty"` // 最大文件大小
	Options     []string `json:"options,omitempty"`      // 选项
}

type Workflow struct {
	Id                uint                    `json:"id"`
	CreatedAt         int64                   `json:"created_at"`
	UpdatedAt         int64                   `json:"updated_at"`
	Name              string                  `json:"name"`
	Icon              string                  `json:"icon"`
	Enabled           bool                    `json:"enabled"`
	Score             int                     `json:"score"`                          // 消耗积分
	Params            []WorkflowParam         `json:"params"`                         // 工作流参数
	Type              string                  `json:"type"`                           // 工作流类型：coze, bailian
	AuthConfig        types.CozeApiConfig     `json:"auth_config,omitempty"`          // Coze 授权配置
	BailianAuthConfig types.BailianApiConfig  `json:"bailian_auth_config,omitempty"`  // 百炼授权配置
	Summary           string                  `json:"summary"`
	WorkflowId        string                  `json:"workflow_id"`
	LastRunAt         int64                   `json:"last_run_at"`
}
