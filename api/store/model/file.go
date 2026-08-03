package model

import "time"

type File struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement;not null"`             // 主键ID
	UserId    int       `gorm:"column:user_id;type:int;not null;index;comment:用户 ID"`    // 用户ID
	Name      string    `gorm:"column:name;type:varchar(100);not null;comment:文件名"`      // 文件名称
	ObjKey    string    `gorm:"column:obj_key;type:varchar(100);comment:文件标识"`           // 对象存储键名
	URL       string    `gorm:"column:url;type:varchar(255);not null;comment:文件地址"`      // 文件访问URL
	Ext       string    `gorm:"column:ext;type:varchar(10);not null;comment:文件后缀"`       // 文件扩展名
	Size      int64     `gorm:"column:size;type:bigint;not null;default:0;comment:文件大小"` // 文件大小(字节)
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`   // 创建时间
}
