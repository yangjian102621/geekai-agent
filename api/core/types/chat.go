package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// ApiRequest API 请求实体
type ApiRequest struct {
	Model               string        `json:"model,omitempty"`
	Temperature         float32       `json:"temperature"`
	MaxTokens           int           `json:"max_tokens,omitempty"`
	MaxCompletionTokens int           `json:"max_completion_tokens,omitempty"` // 兼容GPT O1 模型
	Stream              bool          `json:"stream,omitempty"`
	Messages            []interface{} `json:"messages,omitempty"`
	Tools               []Tool        `json:"tools,omitempty"`
	Functions           []interface{} `json:"functions,omitempty"`       // 兼容中转平台
	ResponseFormat      interface{}   `json:"response_format,omitempty"` // 响应格式

	ToolChoice string `json:"tool_choice,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiResponse struct {
	Choices []ChoiceItem `json:"choices"`
}

// ChoiceItem API 响应实体
type ChoiceItem struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Role             string     `json:"role"`
	Name             string     `json:"name"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	FunctionCall     struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function_call,omitempty"`
}

type AppType string

const (
	AppCoze    = AppType("coze")
	AppDify    = AppType("dify")
	AppOpenAI  = AppType("openai")
	AppBailian = AppType("bailian")
)

// ScoreType 算力日志类型
type ScoreType int

const (
	ScoreRecharge = ScoreType(1) // 充值
	ScoreConsume  = ScoreType(2) // 消费
	ScoreRefund   = ScoreType(3) // 退款
	ScoreInvite   = ScoreType(4) // 邀请奖励
	ScoreRedeem   = ScoreType(5) // 兑换
	ScoreFineTune = ScoreType(6) // 系统调整
)

func (t ScoreType) String() string {
	switch t {
	case ScoreRecharge:
		return "充值"
	case ScoreConsume:
		return "消费"
	case ScoreRefund:
		return "退款"
	case ScoreRedeem:
		return "兑换"
	case ScoreFineTune:
		return "系统调整"
	}
	return "其他"
}

type ScoreMark int

const (
	ScoreSub  = ScoreMark(0) // 支出
	ScorePlus = ScoreMark(1) // 收入
)

type BillingMode string

const (
	BillingModeImmediate    = BillingMode("immediate")
	BillingModeFileSuffix   = BillingMode("file_suffix")
	BillingModeStringMarker = BillingMode("string_marker")
)
