package vo

type ChatMessage struct {
	Id         uint           `json:"id"`
	ChatId     string         `json:"chat_id"`
	UserId     uint           `json:"user_id"`
	AppId      uint           `json:"app_id"`
	Role       string         `json:"role"`
	Icon       string         `json:"icon"`
	Tokens     int            `json:"tokens"`
	Content    MessageContent `json:"content"`
	UseContext bool           `json:"use_context"`
	CreatedAt  int64          `json:"created_at"`
	UpdatedAt  int64          `json:"updated_at"`
}

type MessageContent struct {
	Tools   []ToolCall    `json:"tools,omitempty"`
	Files   []MessageFile `json:"files,omitempty"`
	Texts   []string      `json:"texts,omitempty"`
	Inputs  any           `json:"inputs,omitempty"`  // 参数输入表单
	Answers any           `json:"answers,omitempty"` // 问答答案
	Type    string        `json:"type"`              // 消息类型 text, file, tool, input, answer
}

type ToolCall struct {
	Name   string `json:"name"`
	Status string `json:"status"`          // 工具执行状态 IN_PROGRESS, SUCCESS, FAILED
	Spend  int64  `json:"spend,omitempty"` // 工具执行耗时，单位：毫秒
}

type MessageFile struct {
	Type string `json:"type"` // 文件类型 image, file
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

const (
	ToolCallInProgress = "IN_PROGRESS"
	ToolCallSuccess    = "SUCCESS"
	ToolCallFailed     = "FAILED"
)
