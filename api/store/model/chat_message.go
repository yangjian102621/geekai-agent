package model

import "time"

type ChatMessage struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`                     // 主键ID
	ChatId     string    `gorm:"column:chat_id;type:char(40);not null;index;comment:会话 ID"`       // 会话 ID
	UserId     uint      `gorm:"column:user_id;type:int;not null;index;comment:用户 ID"`            // 用户 ID
	AppId      uint      `gorm:"column:app_id;type:int;not null;index;comment:智能体ID"`             // 智能体 ID
	Role       string    `gorm:"column:role;type:varchar(10);not null;comment:user or ai"`        // ai/user
	Tokens     int       `gorm:"column:tokens;type:smallint;not null;comment:耗费 token 数量"`        // 消息token数量
	Content    string    `gorm:"column:content;type:text;not null;comment:聊天内容"`                  // 消息内容
	UseContext bool      `gorm:"column:use_context;type:tinyint(1);not null;comment:是否允许作为上下文语料"` // 是否可以作为聊天上下文
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null"`                        // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null"`                        // 更新时间
}
