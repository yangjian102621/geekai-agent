package model

import (
	"geekai/core/types"
	"time"
)

// ScoreLog 算力消费日志
type ScoreLog struct {
	Id        uint            `gorm:"column:id;primaryKey;autoIncrement;not null"`                     // 主键ID
	UserId    uint            `gorm:"column:user_id;type:int;not null;index;comment:用户ID"`             // 用户ID
	Username  string          `gorm:"column:username;type:varchar(30);not null;comment:用户名"`           // 用户名
	Type      types.ScoreType `gorm:"column:type;type:tinyint(1);not null;comment:类型（1：充值，2：消费，3：退费）"` // 积分类型
	Amount    int             `gorm:"column:amount;type:smallint;not null;comment:算力数值"`               // 变动数量
	Balance   int             `gorm:"column:balance;type:int;not null;comment:余额"`                     // 变动后余额
	Subject   string          `gorm:"column:subject;type:varchar(50);not null;comment:主题"`             // 主题名称
	Remark    string          `gorm:"column:remark;type:varchar(512);not null;comment:备注"`             // 备注
	Mark      types.ScoreMark `gorm:"column:mark;type:tinyint(1);not null;comment:资金类型（0：支出，1：收入）"`    // 资金类型
	CreatedAt time.Time       `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`           // 创建时间
}
