package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coze-dev/coze-go"
	"github.com/gin-gonic/gin"
)

// Coze 消息处理

type CozeRequest struct {
	BotId             string        `json:"bot_id"`              // 机器人ID
	UserId            string        `json:"user_id"`             // 用户ID
	Stream            bool          `json:"stream"`              // 是否流式响应
	AutoSaveHistory   bool          `json:"auto_save_history"`   // 是否自动保存历史记录
	AdditionalMessage []CozeMessage `json:"additional_messages"` // 附加消息
}

type CozeMessage struct {
	Role        string `json:"role"`         // 角色
	Content     string `json:"content"`      // 内容
	ContentType string `json:"content_type"` // 内容类型
}

type CozeResponse struct {
	Id             string `json:"id"`              // 消息ID
	ConversationId string `json:"conversation_id"` // 会话ID
	BotId          string `json:"bot_id"`          // 机器人ID
	Role           string `json:"role"`            // 角色
	Type           string `json:"type"`            // 消息类型
	Content        string `json:"content"`         // 消息内容
	ContentType    string `json:"content_type"`    // 消息内容类型
	ChatId         string `json:"chat_id"`         // 聊天ID
	SectionId      string `json:"section_id"`      // 分段ID
}

const (
	CozeEventChatCreated      = "conversation.chat.created"
	CozeEventChatInProgress   = "conversation.chat.in_progress"
	CozeEventMessageCompleted = "conversation.message.completed"
	CozeEventMessageDelta     = "conversation.message.delta"
	CozeEventChatCompleted    = "conversation.chat.completed"
	CozeEventDone             = "done"
)

// ActionContentResult 解析 actionContent 的结果
type ActionContentResult struct {
	EventType string // 消息类型：ChatEventInput 或 ChatEventAnswer
	Content   any    // 消息内容
}

// parseActionContent 解析 actionContent，支持两种格式：
// 1. 参数输入节点：JSON 数组格式，如 [{"type":"string","name":"input","required":true}, ...]
// 2. 问答节点：文本格式，如 "对生成的 PPT 是否满意？\n- 不满意，重新生成\n- 满意，直接下载"
func parseActionContent(actionContent string) *ActionContentResult {
	if actionContent == "" {
		return nil
	}

	// 尝试解析为 JSON 数组（参数输入节点）
	var params []map[string]any
	if err := json.Unmarshal([]byte(actionContent), &params); err == nil {
		// 验证是否是有效的参数数组格式
		if len(params) > 0 {
			// 检查第一个元素是否有 type 和 name 字段（参数输入节点的特征）
			firstParam := params[0]
			if _, hasType := firstParam["type"]; hasType {
				if _, hasName := firstParam["name"]; hasName {
					// 是参数输入节点
					return &ActionContentResult{
						EventType: ChatEventInput,
						Content:   params,
					}
				}
			}
		}
	}

	logger.Debugf("actionContent: %s", actionContent)
	// 解析为问答节点（文本格式）
	lines := strings.Split(actionContent, "\n-")
	if len(lines) < 3 {
		return nil
	}

	var question string
	var answers []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 第一行非空行作为问题
		if question == "" {
			question = line
			continue
		}

		// 后面的作为选项
		answer := strings.TrimSpace(line)
		if answer != "" {
			answers = append(answers, answer)
		}
	}

	// 如果解析到了问题和选项，返回问答节点
	if question != "" && len(answers) > 0 {
		return &ActionContentResult{
			EventType: ChatEventAnswer,
			Content: map[string]any{
				"title":   question,
				"answers": answers,
			},
		}
	}

	// 如果只解析到了问题，也返回问答节点（可能只有问题没有选项）
	if question != "" {
		return &ActionContentResult{
			EventType: ChatEventAnswer,
			Content: map[string]interface{}{
				"title":   question,
				"answers": []string{},
			},
		}
	}

	// 无法解析，返回 nil
	return nil
}

func (h *ChatHandler) cozeMessage(c *gin.Context, app model.App, data ChatInput, userId uint) {
	var chatItem model.ChatItem
	h.DB.Where("chat_id", data.ChatId).First(&chatItem)
	var appConfig vo.AppConfig
	err := utils.JsonDecode(app.Configs, &appConfig)
	if err != nil {
		pushMessage(c, ChatEventError, "应用配置解析出错！")
		return
	}

	accessToken, err := h.cozeService.GetAccessToken(&types.CozeApiConfig{
		ApiUrl:      appConfig.ApiUrl,
		AppId:       appConfig.AppId,
		PublicKeyID: appConfig.PublicKeyID,
		PrivateKey:  appConfig.PrivateKey,
	})
	if err != nil {
		pushMessage(c, ChatEventError, "获取 access token 失败："+err.Error())
		return
	}

	customClient := &http.Client{
		Timeout: time.Minute * 20,
	}

	authCli := coze.NewTokenAuth(accessToken)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(appConfig.ApiUrl), coze.WithHttpClient(customClient))

	var messages *coze.Message
	if len(data.Files) == 0 {
		messages = coze.BuildUserQuestionText(data.Prompt, nil)
	} else {

		objects := make([]*coze.MessageObjectString, 0)
		objects = append(objects, coze.NewTextMessageObject(data.Prompt))
		for _, file := range data.Files {
			if utils.IsImage(file.URL) {
				objects = append(objects, coze.NewImageMessageObjectByURL(file.URL))
			} else if utils.IsAudio(file.URL) {
				objects = append(objects, coze.NewAudioMessageObjectByURL(file.URL))
			} else {
				objects = append(objects, coze.NewFileMessageObjectByURL(file.URL))
			}
		}
		messages = coze.BuildUserQuestionObjects(objects, nil)
	}

	req := &coze.CreateChatsReq{
		BotID:          appConfig.BotId,
		ConversationID: chatItem.ConversationId,
		UserID:         fmt.Sprintf("%d", userId),
		Messages:       []*coze.Message{messages},
	}

	logger.Debugf("请求参数: %s", utils.JsonEncode(req))

	resp, err := cozeCli.Chat.Stream(c.Request.Context(), req)
	if err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}

	defer resp.Close()

	finalContent := make([]string, 0)
	pushMessage(c, ChatEventStart, "开始响应")
	tools := make([]vo.ToolCall, 0)
	directReply := false                  // 是否直接回复
	var actionResult *ActionContentResult // 动作解析结果

	tempContent := make([]string, 0)
	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			logger.Info("对话输出结束")
			break
		}
		if err != nil {
			pushMessage(c, ChatEventError, "读取响应失败:"+err.Error())
			break
		}

		// 对话结束
		if event.Event == coze.ChatEventDone {
			break
		}

		logger.Debugf("事件内容: %s", utils.JsonEncode(event))

		switch event.Event {
		// 创建新的对话
		case coze.ChatEventConversationChatCreated:
			chatItem.ConversationId = event.Chat.ConversationID
			session := h.ChatSessions[data.ChatId]
			session.BotId = event.Chat.BotID
			session.ChatId = event.Chat.ID
			session.ConversationId = chatItem.ConversationId
			h.ChatSessions[data.ChatId] = session
		case coze.ChatEventConversationMessageCompleted:
			// 获取最终响应内容
			switch event.Message.Type {
			case coze.MessageTypeAnswer:
				finalContent = append(finalContent, strings.Join(tempContent, ""))
				tempContent = make([]string, 0)
				pushMessage(c, ChatEventMessageCompleted, nil)
			case coze.MessageTypeFunctionCall:
				var plugin map[string]any
				json.Unmarshal([]byte(event.Message.Content), &plugin)
				tool := vo.ToolCall{
					Name:   plugin["plugin"].(string),
					Status: vo.ToolCallInProgress,
					Spend:  time.Now().UnixMilli(),
				}
				tools = append(tools, tool)
				pushMessage(c, ChatEventTool, tool)
			case coze.MessageTypeToolResponse:
				if len(tools) > 0 {
					tool := tools[len(tools)-1]
					tool.Status = vo.ToolCallSuccess
					tool.Spend = time.Now().UnixMilli() - tool.Spend
					tools[len(tools)-1] = tool
					pushMessage(c, ChatEventTool, tool)
				}

				// 处理工作流节点直接输出的内容
				if strings.Contains(event.Message.Content, "directly streaming reply") {
					directReply = true
				}

			case coze.MessageTypeFollowUp:
				pushMessage(c, ChatEventFollowUp, event.Message.Content)
			}

		case coze.ChatEventConversationChatCompleted:
			// TODO: 统计本地对话消耗的 Token

		case coze.ChatEventConversationMessageDelta:
			// 解析工作流输出节点输出的内容
			if directReply {
				actionResult = parseActionContent(event.Message.Content)
				if actionResult != nil {
					pushMessage(c, actionResult.EventType, actionResult.Content)
					break
				}
			}

			// 输出消息
			pushMessage(c, ChatEventMessageDelta, map[string]any{
				"type":    event.Message.ContentType,
				"content": event.Message.Content,
			})
			tempContent = append(tempContent, event.Message.Content)

		case coze.ChatEventConversationChatRequiresAction:
			// 退出直接回复模式，避免重复解析
			directReply = false
			// 推送交互信息到前端
			if actionResult != nil {
				pushMessage(c, actionResult.EventType, actionResult.Content)
			}
		}

	}

	// 更新工具调用状态
	for _, tool := range tools {
		if tool.Status == vo.ToolCallInProgress {
			tool.Status = vo.ToolCallSuccess
			tool.Spend = time.Now().UnixMilli() - tool.Spend
			tools[len(tools)-1] = tool
			pushMessage(c, ChatEventTool, tool)
		}
	}

	// 保存对话
	h.saveChat(c, app, data, userId, &chatItem, finalContent, tools, actionResult)
}

// 取消进行中的对话
func (h *ChatHandler) cozeCancelSession(app model.App, data ChatInput) {
	var appConfig vo.AppConfig
	err := utils.JsonDecode(app.Configs, &appConfig)
	if err != nil {
		logger.Error(err)
		return
	}

	session, ok := h.ChatSessions[data.ChatId]
	if !ok {
		logger.Error("找不到会话信息")
		return
	}

	accessToken, err := h.cozeService.GetAccessToken(&types.CozeApiConfig{
		ApiUrl:      appConfig.ApiUrl,
		AppId:       appConfig.AppId,
		PublicKeyID: appConfig.PublicKeyID,
		PrivateKey:  appConfig.PrivateKey,
	})
	if err != nil {
		logger.Errorf("获取 access token 失败:%v", err)
		return
	}

	authCli := coze.NewTokenAuth(accessToken)
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(appConfig.ApiUrl))
	// 创建取消会话请求
	req := &coze.CancelChatsReq{
		ConversationID: session.ConversationId,
		ChatID:         session.ChatId,
	}

	// 发送取消会话请求
	resp, err := cozeCli.Chat.Cancel(context.Background(), req)
	if err != nil {
		logger.Errorf("取消会话失败:%v", err)
		return
	}

	logger.Infof("取消会话结果: %+v", resp)
}
