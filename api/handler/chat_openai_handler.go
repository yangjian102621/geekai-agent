package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *ChatHandler) openAIMessage(c *gin.Context, ctx context.Context, app model.App, data ChatInput, userId uint) {
	var appConfig vo.AppConfig
	err := utils.JsonDecode(app.Configs, &appConfig)
	if err != nil {
		pushMessage(c, ChatEventError, "解析应用配置失败")
		return
	}

	promptTokens, _ := utils.CalcTokens(data.Prompt, appConfig.ModelName)
	if appConfig.MaxContextLength > 0 && promptTokens > appConfig.MaxContextLength {
		pushMessage(c, ChatEventError, "对话内容超出了当前模型允许的最大上下文长度！")
		return
	}

	// 处理上传文件
	var content any
	if len(data.Files) > 0 {
		fileContents := make([]string, 0)
		messages := []any{
			gin.H{
				"type": "text",
				"text": strings.TrimSpace(data.Prompt),
			},
		}
		for _, file := range data.Files {
			if utils.IsImage(file.URL) {
				messages = append(messages, gin.H{
					"type":      "image_url",
					"image_url": file.URL,
				})
			} else {
				content, err := utils.ReadFileContent(file.URL, h.App.Config.TikaHost)
				if err != nil {
					logger.Error("error with read file: ", err)
				} else {
					fileContents = append(fileContents, content)
				}
			}
		}
		if len(fileContents) > 0 {
			fullPrompt := fmt.Sprintf("请根据提供的文件内容信息回答问题(其中Excel 已转成 HTML)：\n\n %s\n\n 问题：%s", strings.Join(fileContents, "\n"), data.Prompt)
			tokens, _ := utils.CalcTokens(fullPrompt, appConfig.ModelName)
			if tokens > appConfig.MaxContextLength {
				pushMessage(c, ChatEventError, "文件的长度超出模型允许的最大上下文长度，请减少文件内容数量或文件大小。")
				return
			}
			content = fullPrompt
			promptTokens = tokens
		} else {
			content = messages
		}

	} else {
		content = data.Prompt
	}

	// 重新生成逻辑
	dbSession := h.DB.Session(&gorm.Session{}).Where("chat_id", data.ChatId)
	if data.LastMsgId > 0 {
		dbSession = dbSession.Where("id < ?", data.LastMsgId)
		// 删除对应的聊天记录
		h.DB.Where("chat_id", data.ChatId).Where("id >= ?", data.LastMsgId).Delete(&model.ChatMessage{})
	}

	chatCtx := make([]any, 0)
	// 系统预设提示词
	chatCtx = append(chatCtx, model.ChatMessage{
		Role:    "system",
		Content: appConfig.SystemPrompt,
	})
	systemPrompTokens, _ := utils.CalcTokens(appConfig.SystemPrompt, appConfig.ModelName)

	// 检查当前应用是否启用了上下文支持
	if appConfig.EnableContext {
		// 获取历史对话
		num := appConfig.HistoryDeep * 2
		if num == 0 { // 如果历史深度为0，则默认取6条(3轮)
			num = 6
		}

		var chatMessages []model.ChatMessage
		var messages []any

		if err := dbSession.Order("created_at DESC").Limit(num).Find(&chatMessages).Error; err != nil {
			messages = append(messages, model.ChatMessage{
				Role:    "system",
				Content: appConfig.SystemPrompt,
			})
		} else {
			for _, chatMessage := range chatMessages {
				var content struct {
					Type  string   `json:"type"`
					Texts []string `json:"texts"`
				}
				err := utils.JsonDecode(chatMessage.Content, &content)
				if err != nil {
					logger.Error("解析聊天消息失败:", err)
					continue
				}
				messages = append(messages, gin.H{
					"role":    chatMessage.Role,
					"content": content.Texts[0],
				})
			}
		}

		// 计算当前请求的 token 总长度，确保不会超出最大上下文长度
		// MaxContextLength =  Prompt + Response + SystemPrompt + Context
		tokens := promptTokens + systemPrompTokens + appConfig.MaxLength // 最大响应长度
		for i := len(messages) - 1; i >= 0; i-- {
			v := messages[i]
			tks, _ := utils.CalcTokens(utils.JsonEncode(v), appConfig.ModelName)
			// 上下文 token 超出了模型的最大上下文长度
			if tokens+tks >= appConfig.MaxContextLength {
				break
			}
			tokens += tks
			chatCtx = append(chatCtx, v)
		}
		logger.Debugf("聊天上下文：%s", utils.JsonEncode(chatCtx))
	}

	reqMgs := make([]any, 0)
	for i := len(chatCtx) - 1; i >= 0; i-- {
		reqMgs = append(reqMgs, chatCtx[i])
	}
	reqMgs = append(reqMgs, gin.H{
		"role":    "user",
		"content": content,
	})

	// 构造请求体
	requestBody := types.ApiRequest{
		Model:       appConfig.ModelName,
		Messages:    reqMgs,
		Temperature: 0.9,
		Stream:      true,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		pushMessage(c, ChatEventError, "请求参数序列化失败")
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/chat/completions", appConfig.ApiUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		pushMessage(c, ChatEventError, "创建请求失败")
		return
	}
	// 设置请求上下文, 用于中断请求
	req = req.WithContext(ctx)

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", appConfig.Token))

	tm := time.Now()
	logger.Infof("开始发送请求：%s", appConfig.ApiUrl)
	logger.Debug("请求体:", string(jsonData))

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		pushMessage(c, ChatEventError, "请求 API失败:"+err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		pushMessage(c, ChatEventError, fmt.Sprintf("请求 API 失败：%d, %v", response.StatusCode, string(body)))
		return
	}

	logger.Infof("请求结束, 耗时：%dms", time.Since(tm).Milliseconds())
	contentType := response.Header.Get("Content-Type")
	reasoning := false
	if strings.Contains(contentType, "text/event-stream") {

		contents := make([]string, 0)
		// 读取并转发响应流
		pushMessage(c, ChatEventStart, "开始响应")
		reader := bufio.NewReader(response.Body)
		for {
			byteData, _, err := reader.ReadLine()
			line := string(byteData)
			// logger.Info("响应内容:", line)
			if err != nil {
				if err == io.EOF {
					break
				}
				pushMessage(c, ChatEventError, "读取响应失败:"+err.Error())
				break
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			if line == "data: [DONE]" {
				break
			}

			// 解析消息
			var responseBody = types.ApiResponse{}
			err = json.Unmarshal([]byte(line[6:]), &responseBody)
			if err != nil {
				logger.Error("数据解析出错:", err)
				continue
			}
			if len(responseBody.Choices) == 0 { // Fixed: 兼容 Azure API 第一个输出空行
				continue
			}
			if responseBody.Choices[0].Delta.Content == "" &&
				responseBody.Choices[0].Delta.ReasoningContent == "" &&
				len(responseBody.Choices[0].Delta.ToolCalls) == 0 {
				continue
			}

			if responseBody.Choices[0].FinishReason == "stop" && len(contents) == 0 {
				pushMessage(c, ChatEventError, "抱歉😔😔😔，AI助手由于未知原因已经停止输出内容。")
				break
			}
			if responseBody.Choices[0].Delta.ReasoningContent != "" {
				reasoningContent := responseBody.Choices[0].Delta.ReasoningContent
				if !reasoning {
					reasoningContent = fmt.Sprintf("<think>%s", reasoningContent)
					reasoning = true
				}

				pushMessage(c, ChatEventMessageDelta, map[string]any{
					"type":    "text",
					"content": reasoningContent,
				})
				contents = append(contents, reasoningContent)
			} else if responseBody.Choices[0].Delta.Content != "" {
				finalContent := responseBody.Choices[0].Delta.Content
				if reasoning {
					finalContent = fmt.Sprintf("</think>%s", responseBody.Choices[0].Delta.Content)
					reasoning = false
				}
				pushMessage(c, ChatEventMessageDelta, map[string]any{
					"type":    "text",
					"content": finalContent,
				})
				contents = append(contents, finalContent)
			}
		}

		// 更新上下文消息
		if appConfig.EnableContext {
			reqMgs = append(reqMgs, types.Message{
				Role:    "assistant",
				Content: strings.Join(contents, ""),
			})
			h.ChatContexts.Put(data.ChatId, reqMgs)
		}

		var chatItem model.ChatItem
		h.DB.Where("chat_id", data.ChatId).First(&chatItem)

		// 保存对话
		h.saveChat(c, app, data, userId, &chatItem, []string{strings.Join(contents, "")}, nil, nil)
	} else {
		pushMessage(c, ChatEventError, "不支持的响应类型")
	}
}
