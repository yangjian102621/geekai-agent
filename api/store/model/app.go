package model

import (
	"geekai/core/types"
	"time"
)

type App struct {
	Id            uint              `gorm:"column:id;primaryKey;autoIncrement;not null"`                     // 主键ID
	CreatorId     uint              `gorm:"column:creator_id;type:int(11);not null;default:0;comment:创作者ID"` // 创作者ID
	Name          string            `gorm:"column:name;type:varchar(30);comment:名称"`
	Type          types.AppType     `gorm:"column:type;type:varchar(10);not null;default:openai;comment:openai,dify,coze"`
	Cid           uint              `gorm:"column:cid;type:int(6);not null;comment:分类ID"`
	BotId         string            `gorm:"column:bot_id;type:varchar(30);not null;comment:机器人ID（coze 专用）"`
	Enabled       bool              `gorm:"column:enabled;type:tinyint(1);comment:是否启用"`
	Configs       string            `gorm:"column:configs;type:text;not null;comment:智能体配置参数"`
	Params        string            `gorm:"column:params;type:text;comment:应用参数"`
	Score         int               `gorm:"column:score;type:int(11);not null;default:0;comment:单次对话消耗积分"`
	BillingMode   types.BillingMode `gorm:"column:billing_mode;type:varchar(20);not null;default:immediate;comment:扣费模式"`
	BillingConfig string            `gorm:"column:billing_config;type:text;comment:扣费配置JSON"`
	Icon          string            `gorm:"column:icon;type:varchar(255);comment:应用图标"`
	Summary       string            `gorm:"column:summary;type:varchar(512);comment:应用简介"`
	Check         int8              `gorm:"column:check;type:tinyint;not null;default:0;comment:审核状态 0:未审核 1:审核通过 -1:审核不通过"` // 审核状态
	CheckNote     string            `gorm:"column:check_note;type:varchar(255);comment:审核备注"`                                // 审核备注
	CreatedAt     time.Time         `gorm:"column:created_at;type:datetime;not null"`                                        // 创建时间
	UpdatedAt     time.Time         `gorm:"column:updated_at;type:datetime;not null"`                                        // 更新时间
	IsHot         bool              `gorm:"column:is_hot;type:tinyint(1);not null;default:0;comment:是否热门"`                   // 是否热门
	UseCount      int               `gorm:"column:use_count;type:int(11);not null;default:0;comment:使用次数"`                   // 使用次数
}
