package model

type Config struct {
	Id    uint   `gorm:"column:id;primaryKey;autoIncrement;not null"`               // 主键ID
	Name  string `gorm:"column:name;type:varchar(20);unique;not null;comment:配置名称"` // 配置名称
	Value string `gorm:"column:value;type:text;not null"`                           // 配置值
}
