package vo

import "geekai/core/types"

type App struct {
	Id            uint              `json:"id"`
	CreatedAt     int64             `json:"created_at"`
	UpdatedAt     int64             `json:"updated_at"`
	Name          string            `json:"name"`
	Type          types.AppType     `json:"type"`
	Cid           uint              `json:"cid"`   // 分类ID
	Cname         string            `json:"cname"` // 分类名称
	Enabled       bool              `json:"enabled"`
	Score         int               `json:"score"`
	BillingMode   types.BillingMode `json:"billing_mode"`
	BillingConfig BillingConfig     `json:"billing_config"`
	Icon          string            `json:"icon"`
	Configs       AppConfig         `json:"configs"`
	Params        []WorkflowParam   `json:"params"`
	Summary       string            `json:"summary"`
	BotId         string            `json:"bot_id"`
	Check         int8              `json:"check"`
	CheckNote     string            `json:"check_note"`
	CreatorId     uint              `json:"creator_id"`
	IsHot         bool              `json:"is_hot"`
	UseCount      int               `json:"use_count"`
}

type AppConfig struct {
	ApiUrl string `json:"api_url,omitempty"`
	Token  string `json:"token,omitempty"`
	// Coze 配置
	BotId       string `json:"bot_id,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`   // 授权私钥
	AppId       string `json:"app_id,omitempty"`        // 授权应用ID
	PublicKeyID string `json:"public_key_id,omitempty"` // 授权公钥ID
	// OpenAI 配置
	SystemPrompt     string `json:"system_prompt,omitempty"`      // 系统预设提示词
	ModelName        string `json:"model_name,omitempty"`         // 模型名称
	MaxLength        int    `json:"max_length,omitempty"`         // 最大输出长度
	EnableContext    bool   `json:"enable_context,omitempty"`     // 是否启用上下文
	MaxContextLength int    `json:"max_context_length,omitempty"` // 最大上下文长度
	HistoryDeep      int    `json:"history_deep,omitempty"`       // 保留历史对话轮数
	// 百炼配置
	BailianApiKey string `json:"bailian_api_key,omitempty"` // 百炼 API Key
	BailianAppId  string `json:"bailian_app_id,omitempty"`  // 百炼应用 ID
}

// BillingConfig 扣费配置
type BillingConfig struct {
	Suffixes []string `json:"suffixes"` // 文件后缀列表
	Marker   string   `json:"marker"`   // 字符串标记
}

// 审核状态
const (
	CheckStatusPending = iota // 待审核
	CheckStatusPass
	CheckStatusReject
)
