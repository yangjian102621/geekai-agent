package model

import (
	"geekai/core/types"
	"time"
)

type Order struct {
	Id        uint              `gorm:"column:id;primaryKey;autoIncrement;not null"`                                       // 主键ID
	UserId    int               `gorm:"column:user_id;type:int;not null;comment:用户ID"`                                     // 用户ID
	Pid       int               `gorm:"column:pid;type:int;not null;comment:产品ID"`                                         // 产品ID
	Username  string            `gorm:"column:username;type:varchar(30);not null;comment:用户名"`                             // 用户名
	OrderNo   string            `gorm:"column:order_no;type:varchar(30);not null;unique;comment:订单ID"`                     // 订单ID
	TradeNo   string            `gorm:"column:trade_no;type:varchar(60);comment:支付平台交易流水号"`                                // 支付平台交易流水号
	Subject   string            `gorm:"column:subject;type:varchar(100);not null;comment:订单产品"`                            // 订单产品
	Amount    float64           `gorm:"column:amount;type:decimal(10,2);not null;default:0.00;comment:订单金额"`               // 订单金额
	Status    types.OrderStatus `gorm:"column:status;type:tinyint(1);not null;default:0;comment:订单状态（0：待支付，1：已扫码，2：支付成功）"` // 订单状态
	Remark    string            `gorm:"column:remark;type:varchar(255);not null;comment:备注"`                               // 备注
	PayTime   int64             `gorm:"column:pay_time;type:int;comment:支付时间"`                                             // 支付时间
	PayWay    string            `gorm:"column:pay_way;type:varchar(20);not null;comment:支付方式:alipay,wxpay"`                // 支付方式
	Channel   string            `gorm:"column:channel;type:varchar(30);not null;comment:支付类型渠道：支付宝，微信，聚合支付"`               // 支付类型渠道
	CreatedAt time.Time         `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`                             // 创建时间
	UpdatedAt time.Time         `gorm:"column:updated_at;type:datetime;not null;comment:更新时间"`                             // 更新时间
	Checked   bool              `gorm:"column:checked;type:tinyint;not null;default:0;comment:是否已检查"`                      // 是否已检查
}

func (Order) TableName() string {
	return "geekai_orders"
}
