package model

import (
	"time"

	"geekai/core/types"
)

// CreatorScoreLog 创作者积分日志
type CreatorScoreLog struct {
	Id        uint                   `gorm:"column:id;primaryKey;autoIncrement;not null"`                          // 主键ID
	UserId    uint                   `gorm:"column:user_id;type:int;not null;index;comment:用户ID"`                  // 用户ID
	CreatorId uint                   `gorm:"column:creator_id;type:int;not null;index;comment:创作者ID"`              // 创作者ID
	AppId     uint                   `gorm:"column:app_id;type:int;not null;default:0;comment:应用ID"`               // 应用ID
	Type      types.CreatorScoreType `gorm:"column:type;type:char(20);not null;comment:类型（income：收入，withdraw：提现）"` // 积分类型
	Score     int                    `gorm:"column:score;type:int(11);not null;comment:积分数值"`                      // 积分数值
	Balance   int                    `gorm:"column:balance;type:int;not null;comment:余额"`                          // 变动后余额
	Subject   string                 `gorm:"column:subject;type:varchar(50);not null;comment:主题"`                  // 主题名称
	Remark    string                 `gorm:"column:remark;type:varchar(512);not null;comment:备注"`                  // 备注
	Mark      types.ScoreMark        `gorm:"column:mark;type:tinyint(1);not null;comment:资金类型（0：支出，1：收入）"`         // 资金类型
	CreatedAt time.Time              `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`                // 创建时间
}

func (CreatorScoreLog) TableName() string {
	return "geekai_creator_score_logs"
}
