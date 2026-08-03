package model

import (
	"time"
)

type Product struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`                    // 主键ID
	Name      string    `gorm:"column:name;type:varchar(100);not null;comment:产品名称"`            // 产品名称
	Price     float64   `gorm:"column:price;type:decimal(10,2);not null;comment:产品价格"`          // 产品价格
	Credit    int       `gorm:"column:credit;type:int;not null;comment:积分额度"`                   // 积分额度
	Enabled   bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1;comment:启用状态"` // 启用状态
	Sales     int       `gorm:"column:sales;type:int;not null;default:0;comment:销量"`            // 销量
	SortNum   int       `gorm:"column:sort_num;type:tinyint;not null;default:0;comment:排序"`     // 排序
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`          // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`          // 更新时间
}

func (p *Product) TableName() string {
	return "geekai_products"
}
