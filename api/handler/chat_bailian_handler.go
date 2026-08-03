package handler

import (
	"geekai/service/bailian"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *ChatHandler) bailianMessage(c *gin.Context, app model.App, data ChatInput, userId uint) {
	var appConfig vo.AppConfig
	err := utils.JsonDecode(app.Configs, &appConfig)
	if err != nil {
		pushMessage(c, ChatEventError, "解析智能体参数失败："+err.Error())
		return
	}

	if appConfig.BailianApiKey == "" || appConfig.BailianAppId == "" {
		pushMessage(c, ChatEventError, "百炼应用配置不完整，请检查 API Key 和 App ID")
		return
	}

	ctx := c.Request.Context()
	client := bailian.NewClient(appConfig.BailianApiKey, appConfig.BailianAppId)

	// 构建请求
	req := &bailian.CompletionRequest{
		Input: bailian.CompletionInput{
			Prompt: data.Prompt,
		},
		Parameters: bailian.CompletionParameters{
			IncrementalOutput: true,
			HasThoughts:       true,
		},
	}

	// 获取会话ID，用于多轮对话
	var chatItem model.ChatItem
	h.DB.Where("chat_id = ?", data.ChatId).First(&chatItem)
	if chatItem.ConversationId != "" {
		req.Input.SessionId = chatItem.ConversationId
	}

	// 处理文件上传
	if len(data.Files) > 0 {
		for _, file := range data.Files {
			if utils.IsImage(file.URL) {
				req.Input.ImageList = append(req.Input.ImageList, bailian.ImageInfo{
					Image: file.URL,
				})
			} else {
				req.Input.FileList = append(req.Input.FileList, bailian.FileInfo{
					Name: file.Name,
					Url:  file.URL,
				})
			}
		}
	}

	var res chan bailian.StreamChannelResponse
	if res, err = client.API().CompletionStream(ctx, req); err != nil {
		pushMessage(c, ChatEventError, "启动对话失败："+err.Error())
		return
	}

	pushMessage(c, ChatEventStart, "开始响应")
	finalContent := ""
	tools := make([]vo.ToolCall, 0)
	completed := false
	for !completed {
		select {
		case <-ctx.Done():
			completed = true
		case stream, isOpen := <-res:
			if !isOpen {
				pushMessage(c, ChatEventMessageCompleted, "对话结束")
				completed = true
				break
			}
			if stream.Err != nil {
				pushMessage(c, ChatEventError, stream.Err.Error())
				completed = true
				break
			}

			// 更新 session_id（用于多轮对话）
			if stream.Output.SessionId != "" && chatItem.ConversationId != stream.Output.SessionId {
				chatItem.ConversationId = stream.Output.SessionId
				session := h.ChatSessions[data.ChatId]
				session.ConversationId = stream.Output.SessionId
				h.ChatSessions[data.ChatId] = session
			}

			// 处理工具调用（thoughts）
			for _, thought := range stream.Output.Thoughts {
				tool := vo.ToolCall{
					Name:   thought.ActionName,
					Status: vo.ToolCallInProgress,
					Spend:  time.Now().UnixMilli(),
				}
				tools = append(tools, tool)
				pushMessage(c, ChatEventTool, tool)

				// 工具调用完成
				if thought.Response != "" {
					tool.Status = vo.ToolCallSuccess
					tool.Spend = time.Now().UnixMilli() - tool.Spend
					if tool.Spend == 0 {
						tool.Spend = 1000
					}
					tools[len(tools)-1] = tool
					pushMessage(c, ChatEventTool, tool)
				}
			}

			// 输出增量文本
			if stream.Output.Text != "" {
				pushMessage(c, ChatEventMessageDelta, map[string]any{
					"type":    "text",
					"content": stream.Output.Text,
				})
				finalContent += stream.Output.Text
			}

			// 对话完成
			if stream.Output.FinishReason == "stop" {
				pushMessage(c, ChatEventMessageCompleted, "对话结束")
				completed = true
			}
		}
	}

	// 如果最终内容为空，则不保存对话
	if finalContent == "" {
		return
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
	h.saveChat(c, app, data, userId, &chatItem, []string{finalContent}, tools, nil)
}

func (h *ChatHandler) bailianCancelSession(app model.App, data ChatInput) {
	// 百炼无需主动取消，session_id 1小时后自动过期
}
