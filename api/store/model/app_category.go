package model

import (
	"time"
)

type AppCategory struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`               // 主键ID
	CreatorId uint      `gorm:"column:creator_id;type:int;not null;comment:创作者ID"`         // 创作者ID
	Name      string    `gorm:"column:name;type:varchar(30);not null;comment:分类名称"`        // 分类名称
	Enabled   bool      `gorm:"column:enabled;type:tinyint;not null;default:0;comment:状态"` // 状态
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`     // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`     // 更新时间
}

func (AppCategory) TableName() string {
	return "geekai_app_categories"
}
