package model

import "time"

type ChatItem struct {
	Id             uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`                      // 主键ID
	ChatId         string    `gorm:"column:chat_id;type:char(40);unique;not null;comment:会话 ID"`       // 会话 ID
	UserId         uint      `gorm:"column:user_id;type:int;not null;index;comment:用户 ID"`             // 用户 ID
	AppId          uint      `gorm:"column:app_id;type:int;not null;index;comment:智能体ID"`              // 智能体 ID
	Title          string    `gorm:"column:title;type:varchar(100);not null;comment:会话标题"`             // 会话标题
	Icon           string    `gorm:"column:icon;type:varchar(255);not null;comment:图标地址"`              // 会话图标
	ConversationId string    `gorm:"column:conversation_id;type:varchar(100);comment:会话ID(coze/dify)"` // 会话 ID(coze/dify)
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null"`                         // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null"`                         // 更新时间
}
