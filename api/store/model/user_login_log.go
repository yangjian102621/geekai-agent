package model

import "time"

type UserLoginLog struct {
	Id           uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`                 // 主键ID
	UserId       uint      `gorm:"column:user_id;type:int;not null;index;comment:用户ID"`         // 用户ID
	Username     string    `gorm:"column:username;type:varchar(30);not null;comment:用户名"`       // 用户名
	LoginIp      string    `gorm:"column:login_ip;type:char(16);not null;comment:登录IP"`         // 登录IP
	LoginAddress string    `gorm:"column:login_address;type:varchar(30);not null;comment:登录地址"` // 登录地址
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null"`                    // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null"`                    // 更新时间
}
