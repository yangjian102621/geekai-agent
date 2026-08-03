package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OrderStatus int

const (
	OrderNotPaid     = OrderStatus(0)
	OrderPaidSuccess = OrderStatus(1)
	OrderPaidFailed  = OrderStatus(2)
	OrderNotExist    = OrderStatus(3)
)

type OrderRemark struct {
	Credit  int     `json:"credit"`  // 积分点数
	Subject string  `json:"subject"` // 商品名称
	Price   float64 `json:"price"`   // 商品价格
}

// PayChannel 支付渠道
var PayChannel = map[string]string{
	"alipay":  "支付宝商号",
	"wxpay":   "微信商号",
	"geekpay": "易支付",
}

var PayWays = map[string]string{
	"alipay": "支付宝",
	"wxpay":  "微信支付",
}
