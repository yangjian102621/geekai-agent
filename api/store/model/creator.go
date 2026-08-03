package model

import "time"

// Creator 创作者
type Creator struct {
	Id              uint      `gorm:"column:id;type:int(11);primaryKey;autoIncrement;not null;comment:主键ID"`             // 主键ID
	UserId          uint      `gorm:"column:user_id;type:int(11);not null;unique;comment:关联用户ID"`                        // 关联用户ID
	Username        string    `gorm:"column:username;type:varchar(30);unique;comment:用户名"`                               // 用户名
	Fee             int       `gorm:"column:fee;type:smallint;default:0;comment:提现费率(0-100)"`                            // 提现费率(0-100)
	Name            string    `gorm:"column:name;type:varchar(100);not null;comment:创作者名称"`                              // 创作者名称
	Description     string    `gorm:"column:description;type:varchar(512);comment:创作者简介"`                                // 创作者简介
	Logo            string    `gorm:"column:logo;type:varchar(255);comment:创作者Logo"`                                     // 创作者Logo
	Enabled         bool      `gorm:"column:enabled;type:tinyint(1);default:1;not null;comment:是否启用"`                    // 是否启用
	Scores          int       `gorm:"column:scores;type:int;default:0;comment:积分"`                                       // 积分
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`                             // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`                             // 更新时间
	Check           int8      `gorm:"column:check;type:tinyint(1);default:0;not null;comment:审核状态 0:未审核 1:审核通过 2:审核不通过"` // 审核状态
	CheckNote       string    `gorm:"column:check_note;type:varchar(255);comment:审核备注"`                                  // 审核备注
	WithdrawConfigs string    `gorm:"column:withdraw_configs;type:text;comment:提现配置"`                                    // 提现配置
}
