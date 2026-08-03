package model

import "time"

type User struct {
	Id          uint      `gorm:"column:id;primaryKey;autoIncrement;not null;comment:主键ID"`       // 主键ID
	Username    string    `gorm:"column:username;type:varchar(30);not null;unique;comment:用户名"`   // 用户名
	Nickname    string    `gorm:"column:nickname;type:varchar(30);not null;comment:昵称"`           // 昵称
	Password    string    `gorm:"column:password;type:char(64);not null;comment:密码"`              // 密码
	Avatar      string    `gorm:"column:avatar;type:varchar(255);not null;comment:头像"`            // 头像
	Salt        string    `gorm:"column:salt;type:char(12);not null;comment:密码盐"`                 // 密码盐
	Scores      int       `gorm:"column:scores;type:int;default:0;comment:剩余积分"`                  // 剩余积分
	ExpiredTime int64     `gorm:"column:expired_time;type:int;not null;comment:账户到期时间"`           // 账户到期时间
	Enabled     bool      `gorm:"column:enabled;type:tinyint(1);default:1;not null;comment:当前状态"` // 当前状态
	LastLoginAt int64     `gorm:"column:last_login_at;type:int;not null;comment:最后登录时间"`          // 最后登录时间
	LastLoginIp string    `gorm:"column:last_login_ip;type:char(16);not null;comment:最后登录IP"`     // 最后登录 IP
	OpenId      string    `gorm:"column:openid;type:varchar(100);comment:第三方平台OpenID"`            // 第三方平台OpenID
	Platform    string    `gorm:"column:platform;type:varchar(30);comment:第三方平台类型"`               // 第三方平台类型
	Vip         bool      `gorm:"column:vip;type:tinyint(1);default:0;comment:是否VIP会员"`           // 是否 VIP 会员
	Invitor     uint      `gorm:"column:invitor;type:int;default:0;comment:邀请人ID"`                // 邀请人ID
	InviteCode  string    `gorm:"column:invite_code;type:varchar(100);comment:邀请码"`               // 邀请码
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`          // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`          // 更新时间
}
