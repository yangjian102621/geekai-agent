package vo

import "geekai/core/types"

type Order struct {
	Id          uint              `json:"id"`
	UserId      uint              `json:"user_id"`
	Pid         uint              `json:"pid"`
	Username    string            `json:"username"`
	OrderNo     string            `json:"order_no"`
	TradeNo     string            `json:"trade_no"`
	Subject     string            `json:"subject"`
	Amount      float64           `json:"amount"`
	Status      types.OrderStatus `json:"status"`
	Remark      types.OrderRemark `json:"remark"`
	PayTime     int64             `json:"pay_time"`
	PayWay      string            `json:"pay_way"`
	Channel     string            `json:"channel"`
	ChannelName string            `json:"channel_name"`
	PayName     string            `json:"pay_name"`
	CreatedAt   int64             `json:"created_at"`
	UpdatedAt   int64             `json:"updated_at"`
}
