package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 智能体会话信息，取消会话的时候需要用到
type ChatSession struct {
	ChatId         string             // 智能体业务会话ID
	ConversationId string             // 智能体会话ID
	BotId          string             // 智能体ID
	CancelFunc     context.CancelFunc // 取消会话的函数
}

type ChatHandler struct {
	BaseHandler
	userService    *service.UserService
	redis          *redis.Client
	ChatSessions   map[string]ChatSession     // 业务 ChatId -> ChatSession
	ChatContexts   *types.LMap[string, []any] // 聊天上下文 Map [chatId] => []Message
	liceseService  *service.LicenseService
	cozeService    *service.CozeService
	creatorService *service.CreatorService
}

var (
	// for chat
	ChatEventMessageCompleted = "message.completed"
	ChatEventMessageDelta     = "message.delta"
	ChatEventTool             = "tool"
	ChatEventStart            = "start"
	ChatEventEnd              = "end"
	ChatEventError            = "error"
	ChatEventTitle            = "title"
	ChatEventInput            = "input"     // 表单输入节点
	ChatEventAnswer           = "answer"    // 问答节点
	ChatEventFollowUp         = "follow_up" // 跟随问题操作
)

func NewChatHandler(
	app *core.AppServer,
	db *gorm.DB,
	userService *service.UserService,
	redis *redis.Client,
	licenseService *service.LicenseService,
	cozeService *service.CozeService,
	creatorService *service.CreatorService,
) *ChatHandler {
	return &ChatHandler{
		BaseHandler:    BaseHandler{App: app, DB: db},
		userService:    userService,
		redis:          redis,
		ChatSessions:   make(map[string]ChatSession),
		ChatContexts:   types.NewLMap[string, []any](),
		liceseService:  licenseService,
		cozeService:    cozeService,
		creatorService: creatorService,
	}
}

// RegisterRoutes 注册路由
func (h *ChatHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/chat/")

	// 需要用户登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		// 每2秒限制1次请求
		rg.POST("message", middleware.RateLimitEvery(h.App.Redis, 2*time.Second), h.Chat)
		rg.GET("list", h.List)
		rg.POST("update", h.Update)
		rg.GET("remove", h.Remove)
		rg.GET("clear", h.Clear)
		rg.GET("messages", h.Messages)
		rg.GET("last-message", h.GetLastMessage)
		rg.POST("cancel", h.CancelSession)
		rg.GET("detail", h.Detail)
	}
}

type ChatInput struct {
	ChatId    string           `json:"chat_id,omitempty"`
	AppId     int              `json:"app_id,omitempty"`
	Prompt    string           `json:"prompt,omitempty"`
	Files     []vo.MessageFile `json:"files,omitempty"`
	LastMsgId uint             `json:"last_msg_id,omitempty"`
}

// Chat 处理聊天请求
func (h *ChatHandler) Chat(c *gin.Context) {
	var data ChatInput
	userId := h.GetLoginUserId(c)

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	if err := c.ShouldBindJSON(&data); err != nil {
		pushMessage(c, ChatEventError, "参数解析失败："+err.Error())
		return
	}

	// 从数据库获取应用配置
	var app model.App
	if err := h.DB.Where("id", data.AppId).First(&app).Error; err != nil {
		pushMessage(c, ChatEventError, "获取应用配置失败，请联系管理员去添加应用。")
		return
	}

	if !app.Enabled {
		pushMessage(c, ChatEventError, "当前应用已禁用，请联系管理员！")
		return
	}

	// 检查用户积分是否足够
	var user model.User
	if err := h.DB.Where("id", userId).First(&user).Error; err != nil {
		pushMessage(c, ChatEventError, "获取用户信息失败")
		return
	}
	if !user.Enabled {
		pushMessage(c, ChatEventError, "当前用户已被禁用，请联系管理员！")
		return
	}

	if user.Scores < app.Score {
		pushMessage(c, ChatEventError, "当前用户积分不足以完成本次对话，请先购买积分！")
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	// 检查是否存在会话
	if session, ok := h.ChatSessions[data.ChatId]; ok {
		session.CancelFunc = cancel
		h.ChatSessions[data.ChatId] = session
	} else {
		h.ChatSessions[data.ChatId] = ChatSession{
			CancelFunc: cancel,
		}
	}

	switch app.Type {
	case types.AppOpenAI:
		h.openAIMessage(c, ctx, app, data, userId)
	case types.AppCoze:
		h.cozeMessage(c, app, data, userId)
	case types.AppDify:
		h.difyMessage(c, app, data, userId)
	case types.AppBailian:
		h.bailianMessage(c, app, data, userId)
	default:
		pushMessage(c, ChatEventError, "不支持的应用类型")
	}
	pushMessage(c, ChatEventEnd, "对话结束")
}

// 取消会话
func (h *ChatHandler) CancelSession(c *gin.Context) {
	var data ChatInput
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if data.ChatId == "" || data.AppId == 0 {
		resp.ERROR(c, "app_id 和 chat_id 不能为空")
		return
	}
	// 从数据库获取应用配置
	var app model.App
	if err := h.DB.Where("id", data.AppId).First(&app).Error; err != nil {
		resp.ERROR(c, "找不到应用")
		return
	}
	switch app.Type {
	case "openai":
		if session, ok := h.ChatSessions[data.ChatId]; ok {
			session.CancelFunc()
		}
	case "coze":
		h.cozeCancelSession(app, data)
	case "dify":
		h.difyCancelSession(app, data)
	case "bailian":
		h.bailianCancelSession(app, data)
	}
	resp.SUCCESS(c, types.OkMsg)
}

// List 获取会话列表
func (h *ChatHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	var items = make([]vo.ChatItem, 0)
	var chats []model.ChatItem
	h.DB.Where("user_id", userId).Order("id DESC").Find(&chats)
	if len(chats) == 0 {
		resp.SUCCESS(c, items)
		return
	}

	appIds := make([]uint, 0)
	for _, chat := range chats {
		appIds = append(appIds, chat.AppId)
	}
	var apps []model.App
	h.DB.Where("id IN ?", appIds).Find(&apps)
	appMap := make(map[uint]model.App)
	for _, app := range apps {
		appMap[app.Id] = app
	}

	for _, chat := range chats {
		var item vo.ChatItem
		err := utils.CopyObject(chat, &item)
		if err != nil {
			continue
		}
		item.Icon = appMap[chat.AppId].Icon
		if item.Icon == "" {
			item.Icon = "/images/app-placeholder.png"
		}
		items = append(items, item)
	}
	resp.SUCCESS(c, items)
}

// Update 更新会话标题
func (h *ChatHandler) Update(c *gin.Context) {
	var data struct {
		ChatId string `json:"chat_id"`
		Title  string `json:"title"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	res := h.DB.Model(&model.ChatItem{}).Where("chat_id = ?", data.ChatId).UpdateColumn("title", data.Title)
	if res.Error != nil {
		resp.ERROR(c, "Failed to update database")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
}

// Clear 清空所有聊天记录
func (h *ChatHandler) Clear(c *gin.Context) {
	// 获取当前登录用户所有的聊天会话
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	var chats []model.ChatItem
	res := h.DB.Where("user_id = ?", user.Id).Find(&chats)
	if res.Error != nil {
		resp.ERROR(c, "No chats found")
		return
	}

	var chatIds = make([]string, 0)
	for _, chat := range chats {
		chatIds = append(chatIds, chat.ChatId)
	}
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		res := h.DB.Where("user_id", user.Id).Delete(&model.ChatItem{})
		if res.Error != nil {
			return res.Error
		}

		res = h.DB.Where("user_id", user.Id).Where("chat_id IN ?", chatIds).Delete(&model.ChatMessage{})
		if res.Error != nil {
			return res.Error
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error with delete chats: %+v", err)
		resp.ERROR(c, "Failed to remove chat from database.")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
}

// Messages 获取聊天历史记录
func (h *ChatHandler) Messages(c *gin.Context) {
	chatId := c.Query("chat_id") // 会话 ID
	var items []model.ChatMessage
	var messages = make([]vo.ChatMessage, 0)
	err := h.DB.Where("chat_id = ?", chatId).Find(&items).Error
	if err != nil {
		resp.ERROR(c, "No history message")
		return
	}

	// 获取应用图标
	var appIds = make([]uint, 0)
	for _, item := range items {
		appIds = append(appIds, item.AppId)
	}
	var apps []model.App
	h.DB.Where("id IN ?", appIds).Find(&apps)
	appMap := make(map[uint]model.App)
	for _, app := range apps {
		appMap[app.Id] = app
	}

	for _, item := range items {
		var v vo.ChatMessage
		err := utils.CopyObject(item, &v)
		if err != nil {
			continue
		}
		v.Id = item.Id
		v.CreatedAt = item.CreatedAt.Unix()
		v.UpdatedAt = item.UpdatedAt.Unix()
		err = utils.JsonDecode(item.Content, &v.Content)
		if err != nil {
			v.Content = vo.MessageContent{
				Texts: []string{item.Content},
				Tools: nil,
			}
		}
		v.Icon = appMap[item.AppId].Icon
		if v.Icon == "" {
			v.Icon = "/images/app-placeholder.png"
		}
		messages = append(messages, v)
	}

	resp.SUCCESS(c, messages)
}

// Remove 删除会话
func (h *ChatHandler) Remove(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if chatId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	res := h.DB.Where("user_id", user.Id).Where("chat_id", chatId).Delete(&model.ChatItem{})
	if res.Error != nil {
		resp.ERROR(c, "Failed to update database")
		return
	}

	// 删除当前会话的聊天记录
	res = h.DB.Where("user_id", user.Id).Where("chat_id", chatId).Delete(&model.ChatMessage{})
	if res.Error != nil {
		resp.ERROR(c, "Failed to remove chat from database.")
		return
	}

	resp.SUCCESS(c, types.OkMsg)
}

// Detail 对话详情，用户导出对话
func (h *ChatHandler) Detail(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if utils.IsEmptyValue(chatId) {
		resp.ERROR(c, "Invalid chatId")
		return
	}

	var chatItem model.ChatItem
	res := h.DB.Where("chat_id = ?", chatId).First(&chatItem)
	if res.Error != nil {
		resp.ERROR(c, "No chat found")
		return
	}

	var chatItemVo vo.ChatItem
	err := utils.CopyObject(chatItem, &chatItemVo)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, chatItemVo)
}

// 保存对话信息
func (h *ChatHandler) saveChat(c *gin.Context,
	app model.App,
	data ChatInput,
	userId uint,
	chatItem *model.ChatItem,
	contents []string,
	tools []vo.ToolCall,
	actionResult *ActionContentResult,
) {
	if chatItem.Id == 0 {
		chatItem = &model.ChatItem{
			ChatId:         data.ChatId,
			UserId:         userId,
			AppId:          app.Id,
			ConversationId: chatItem.ConversationId,
			Title:          utils.TruncateString(data.Prompt, 30),
			Icon:           app.Icon,
		}
		if err := h.DB.Create(chatItem).Error; err != nil {
			logger.Error("保存聊天记录失败:", err)
		} else {
			pushMessage(c, ChatEventTitle, chatItem.Title)
		}
	}
	messageType := "text"
	if len(data.Files) > 0 {
		messageType = "file"
	}

	// 保存对话消息
	message := model.ChatMessage{
		ChatId: data.ChatId,
		UserId: userId,
		AppId:  app.Id,
		Role:   "user",
		Content: utils.JsonEncode(vo.MessageContent{
			Files: data.Files,
			Texts: []string{data.Prompt},
			Type:  messageType,
		}),
	}
	if err := h.DB.Create(&message).Error; err != nil {
		logger.Error("保存聊天记录失败:", err)
		return
	}

	replyContent := vo.MessageContent{
		Texts: contents,
		Tools: tools,
	}

	if actionResult != nil {
		switch actionResult.EventType {
		case ChatEventInput:
			replyContent.Inputs = actionResult.Content
		case ChatEventAnswer:
			replyContent.Answers = actionResult.Content
		}
		messageType = actionResult.EventType
	}
	replyContent.Type = messageType
	message = model.ChatMessage{
		ChatId:  data.ChatId,
		UserId:  userId,
		AppId:   app.Id,
		Role:    "assistant",
		Content: utils.JsonEncode(replyContent),
	}
	if err := h.DB.Create(&message).Error; err != nil {
		logger.Error("保存聊天记录失败:", err)
	}

	// 条件扣费：只有积分大于0时才检查扣费条件
	if app.Score > 0 {
		triggerReason := h.shouldDeductScore(app, replyContent)
		if triggerReason != "" {
			h.subUserPower(userId, app, data.Prompt, triggerReason)
			// 增加应用使用次数
			h.DB.Model(&model.App{}).Where("id", app.Id).Update("use_count", gorm.Expr("use_count + 1"))
		}
	}
}

func (h *ChatHandler) subUserPower(userId uint, app model.App, prompt string, triggerReason string) {
	power := 1
	if app.Score > 0 {
		power = app.Score
	}

	remark := fmt.Sprintf("应用名称：%s, 提问:%s, 扣费方式：%s", app.Name, utils.TruncateString(prompt, 20), triggerReason)
	err := h.userService.DecreaseScores(userId, power, model.ScoreLog{
		Type:    types.ScoreConsume,
		Subject: app.Name,
		Remark:  remark,
	})

	if app.CreatorId > 0 {
		err = h.creatorService.IncreaseScores(app.CreatorId, power, model.CreatorScoreLog{
			Type:    types.CreatorScoreTypeIncome,
			Subject: app.Name,
			Remark:  remark,
		})
	}
	if err != nil {
		logger.Error(err)
	}
}

// shouldDeductScore 判断是否应该扣费，返回扣费触发原因，空字符串表示不扣费
func (h *ChatHandler) shouldDeductScore(app model.App, content vo.MessageContent) string {
	// 根据扣费模式判断
	switch app.BillingMode {
	case types.BillingModeImmediate:
		return "立即扣费"
	case types.BillingModeFileSuffix:
		return h.checkFileSuffix(app, content)
	case types.BillingModeStringMarker:
		return h.checkStringMarker(app, content)
	default:
		return "立即扣费" // 默认立即扣费
	}
}

// checkFileSuffix 检查文件后缀是否匹配
func (h *ChatHandler) checkFileSuffix(app model.App, content vo.MessageContent) string {
	// 解析扣费配置
	var billingConfig vo.BillingConfig
	if err := utils.JsonDecode(app.BillingConfig, &billingConfig); err != nil {
		logger.Errorf("解析扣费配置失败: %v", err)
		return ""
	}

	// 如果没有配置后缀，不扣费
	if len(billingConfig.Suffixes) == 0 {
		return ""
	}

	logger.Infof("文件后缀配置: %+v", billingConfig)
	logger.Infof("文本内容: %+v", content.Texts)
	// 遍历文本内容，提取URL并检查后缀
	for _, text := range content.Texts {
		urls := utils.ExtractURLsFromText(text)
		for _, url := range urls {
			suffix := utils.ExtractFileSuffixFromURL(url)
			// 检查后缀是否在配置列表中
			for _, configSuffix := range billingConfig.Suffixes {
				if strings.ToLower(configSuffix) == suffix {
					return fmt.Sprintf("文件后缀触发：%s", suffix)
				}
			}
		}
	}

	return "" // 未匹配到文件后缀，不扣费
}

// checkStringMarker 检查字符串标记是否匹配
func (h *ChatHandler) checkStringMarker(app model.App, content vo.MessageContent) string {
	// 解析扣费配置
	var billingConfig vo.BillingConfig
	if err := utils.JsonDecode(app.BillingConfig, &billingConfig); err != nil {
		logger.Errorf("解析扣费配置失败: %v", err)
		return ""
	}

	// 如果没有配置标记，不扣费
	if billingConfig.Marker == "" {
		return ""
	}

	// 遍历文本内容，检查是否包含标记
	for _, text := range content.Texts {
		matched := strings.Contains(strings.ToLower(text), strings.ToLower(billingConfig.Marker))
		if matched {
			return fmt.Sprintf("字符串标记触发：%s", billingConfig.Marker)
		}
	}

	return "" // 未匹配到标记，不扣费
}

// 获取对话的最后一条回复消息
func (h *ChatHandler) GetLastMessage(c *gin.Context) {
	chatId := h.GetTrim(c, "chat_id")
	if utils.IsEmptyValue(chatId) {
		resp.ERROR(c, "Invalid chatId")
		return
	}

	var message model.ChatMessage
	res := h.DB.Where("chat_id = ? AND role = ?", chatId, "assistant").Order("id DESC").First(&message)
	if res.Error != nil {
		resp.ERROR(c, "No message found")
		return
	}
	var v vo.ChatMessage
	err := utils.CopyObject(message, &v)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, v)
}
