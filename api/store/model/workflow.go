package model

import (
	"time"
)

type Workflow struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement;not null"` // 主键ID
	Name       string    `gorm:"column:name;type:varchar(100);comment:名称"`
	WorkflowId string    `gorm:"column:workflow_id;type:varchar(30);comment:workflow_id"`
	Type       string    `gorm:"column:type;type:varchar(10);not null;default:coze;comment:工作流类型:coze,bailian"` // 工作流类型
	Icon       string    `gorm:"column:icon;type:varchar(255);comment:工作流图标"`
	Enabled    bool      `gorm:"column:enabled;type:tinyint(1);comment:是否启用"`
	Params     string    `gorm:"column:params;type:text;not null;comment:工作流参数"`
	AuthConfig string    `gorm:"column:auth_config;type:text;comment:工作流授权配置"` // 工作流授权配置
	Score      int       `gorm:"column:score;type:int;not null;default:0;comment:消耗积分"`
	Summary    string    `gorm:"column:summary;type:varchar(255);comment:工作流简介"`
	LastRunAt  int64     `gorm:"column:last_run_at;type:int(11);not null;default:0;comment:最后运行时间"` // 最后运行时间
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null"`                          // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null"`                          // 更新时间
}

func (Workflow) TableName() string {
	return "geekai_workflows"
}
