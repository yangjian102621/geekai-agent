package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
import (
	"time"

	"geekai/service/dify"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"

	"github.com/gin-gonic/gin"
)

func (h *ChatHandler) difyMessage(c *gin.Context, app model.App, data ChatInput, userId uint) {
	var appConfig vo.AppConfig
	err := utils.JsonDecode(app.Configs, &appConfig)
	if err != nil {
		pushMessage(c, "error", "解析智能体参数失败："+err.Error())
		return
	}

	// 判断License是否有权限
	license := h.liceseService.GetLicense()
	if !license.IsActive {
		pushMessage(c, "error", "当前系统未授权，无法使用Dify应用")
		return
	}

	ctx := c.Request.Context()
	client := dify.NewClient(appConfig.ApiUrl, appConfig.Token)
	req := &dify.ChatMessageRequest{
		Query: data.Prompt,
		User:  data.ChatId,
	}
	// 获取会话ID
	var chatItem model.ChatItem
	h.DB.Where("chat_id = ?", data.ChatId).First(&chatItem)
	if chatItem.ConversationId != "" {
		req.ConversationID = chatItem.ConversationId
	}

	var res chan dify.ChatMessageStreamChannelResponse

	if res, err = client.Api().ChatMessagesStream(ctx, req); err != nil {
		pushMessage(c, "error", "启动对话失败："+err.Error())
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
			if err = stream.Err; err != nil {
				pushMessage(c, ChatEventError, err.Error())
				completed = true
			}
			if !isOpen {
				pushMessage(c, ChatEventMessageCompleted, "对话结束")
				completed = true
			}

			logger.Infof("流数据: %+v", stream.Event)

			// 工具调用
			if stream.Event == dify.MessageEventAgentToolStart {
				tool := vo.ToolCall{
					Name:   stream.ToolName,
					Status: vo.ToolCallInProgress,
					Spend:  time.Now().UnixMilli(),
				}
				tools = append(tools, tool)
				pushMessage(c, ChatEventTool, tool)
			} else if stream.Event == dify.MessageEventAgentToolEnd {
				tool := tools[len(tools)-1]
				tool.Status = vo.ToolCallSuccess
				tool.Spend = time.Now().UnixMilli() - tool.Spend
				if tool.Spend == 0 {
					tool.Spend = 1000
				}
				tools[len(tools)-1] = tool
				pushMessage(c, ChatEventTool, tool)
			} else if stream.Event == dify.MessageEventMessageDelta {
				if stream.Answer == "" {
					break
				}
				pushMessage(c, ChatEventMessageDelta, map[string]any{
					"type":    "text",
					"content": stream.Answer,
				})
			} else if stream.Event == dify.MessageEventAgentThought {
				finalContent = stream.Thought
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

func (h *ChatHandler) difyCancelSession(app model.App, data ChatInput) {

}
