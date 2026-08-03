package model

import "time"

type AdminUser struct {
	Id          uint      `gorm:"column:id;primaryKey;autoIncrement;not null;comment:主键ID"`     // 主键ID
	Username    string    `gorm:"column:username;type:varchar(30);not null;unique;comment:用户名"` // 用户名
	Password    string    `gorm:"column:password;type:char(64);not null;comment:密码"`            // 密码
	Salt        string    `gorm:"column:salt;type:char(12);not null;comment:密码盐"`               // 密码盐
	Status      bool      `gorm:"column:status;type:tinyint(1);not null;comment:当前状态"`          // 当前状态
	LastLoginAt int64     `gorm:"column:last_login_at;type:int;not null;comment:最后登录时间"`        // 最后登录时间
	LastLoginIp string    `gorm:"column:last_login_ip;type:char(16);not null;comment:最后登录 IP"`  // 最后登录 IP
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`        // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`        // 更新时间
}
