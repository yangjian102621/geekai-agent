package model

import "time"

// CreatorWithdraw 创作者提现申请
type CreatorWithdraw struct {
	Id          uint      `gorm:"column:id;type:int(11);primaryKey;autoIncrement;not null;comment:主键ID"`      // 主键ID
	CreatorId   uint      `gorm:"column:creator_id;type:int(11);not null;comment:创作者ID"`                      // 创作者ID
	Scores      int       `gorm:"column:scores;type:int(11);not null;comment:提现积分"`                           // 提现积分
	TotalMoney  float64   `gorm:"column:total_money;type:decimal(10,2);not null;comment:提现总金额"`               // 提现总金额
	RealMoney   float64   `gorm:"column:real_money;type:decimal(10,2);not null;comment:提现到账金额"`               // 提现到账金额
	Fee         float64   `gorm:"column:fee;type:decimal(10,2);not null;comment:提现手续费"`                       // 提现手续费
	Method      string    `gorm:"column:method;type:varchar(20);not null;comment:收款方式(alipay/wxpay)"`         // 收款方式
	Account     string    `gorm:"column:account;type:varchar(100);not null;comment:收款账号"`                     // 收款账号
	AccountName string    `gorm:"column:account_name;type:varchar(100);not null;comment:收款人姓名"`               // 收款人姓名
	QrCode      string    `gorm:"column:qr_code;type:varchar(255);not null;comment:收款二维码"`                    // 收款二维码(图片)
	Status      string    `gorm:"column:status;type:varchar(20);not null;comment:状态(pending/success/reject)"` // 状态
	Note        string    `gorm:"column:note;type:varchar(255);comment:备注"`                                   // 备注
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`                      // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`                      // 更新时间
}

func (CreatorWithdraw) TableName() string {
	return "geekai_creator_withdraws"
}
