package model

import "time"

// 兑换码

type Redeem struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`               // 主键ID
	UserId     uint      `gorm:"column:user_id;type:int;not null;index;comment:用户 ID"`      // 用户 ID
	Name       string    `gorm:"column:name;type:varchar(30);not null;comment:兑换码名称"`       // 名称
	Amount     uint      `gorm:"column:amount;type:int;not null;comment:额度"`                // 兑换额度
	Code       string    `gorm:"column:code;type:varchar(100);unique;not null;comment:兑换码"` // 兑换码
	Enabled    bool      `gorm:"column:enabled;type:tinyint(1);not null;comment:是否启用"`      // 启用状态
	RedeemedAt int64     `gorm:"column:redeemed_at;type:int;not null;comment:兑换时间"`         // 兑换时间
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`     // 创建时间
}
