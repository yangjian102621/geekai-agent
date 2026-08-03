package vo

type ChatItem struct {
	Id             uint   `json:"id"`
	UserId         uint   `json:"user_id"`
	Icon           string `json:"icon"`
	AppId          uint   `json:"app_id"`
	ChatId         string `json:"chat_id"`
	Title          string `json:"title"`
	ConversationId string `json:"conversation_id"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
}
