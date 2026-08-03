package model

import (
	"geekai/store/vo"
	"time"
)

// 工作流任务
type WorkflowTask struct {
	Id           uint                  `gorm:"column:id;primaryKey;autoIncrement;not null"` // 主键ID
	UserId       uint                  `gorm:"column:user_id;type:int;not null;comment:用户ID"`
	TaskId       string                `gorm:"column:task_id;type:varchar(20);unique;comment:任务ID"`
	WorkflowId   string                `gorm:"column:workflow_id;type:varchar(30);comment:工作流 ID"`
	WorkflowName string                `gorm:"column:workflow_name;type:varchar(100);comment:工作流名称"`
	Status       vo.WorkflowTaskStatus `gorm:"column:status;type:varchar(20);comment:状态"`
	Progress     int                   `gorm:"column:progress;type:int;default:0;comment:任务进度"`
	Params       string                `gorm:"column:params;type:text;not null;comment:工作流参数"`
	Score        int                   `gorm:"column:score;type:int;not null;default:0;comment:消耗积分"`
	Output       string                `gorm:"column:output;type:text;comment:输出结果"`
	Error        string                `gorm:"column:error;type:text;comment:错误信息"`
	CreatedAt    time.Time             `gorm:"column:created_at;type:datetime;not null"` // 创建时间
	UpdatedAt    time.Time             `gorm:"column:updated_at;type:datetime;not null"` // 更新时间
}

func (WorkflowTask) TableName() string {
	return "geekai_workflow_tasks"
}
