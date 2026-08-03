package vo

// CreatorWithdraw 创作者提现VO
type CreatorWithdraw struct {
	Id          uint    `json:"id"`           // 主键ID
	CreatorId   uint    `json:"creator_id"`   // 创作者ID
	Scores      int     `json:"scores"`       // 提现积分
	Fee         float64 `json:"fee"`          // 提现手续费
	Method      string  `json:"method"`       // 提现方式(alipay/wxpay)
	Account     string  `json:"account"`      // 收款账号
	AccountName string  `json:"account_name"` // 收款人姓名
	QrCode      string  `json:"qr_code"`      // 收款二维码(图片)
	Status      string  `json:"status"`       // 状态(pending/success/reject)
	Note        string  `json:"note"`         // 备注
	CreatedAt   int64   `json:"created_at"`   // 创建时间
	UpdatedAt   int64   `json:"updated_at"`   // 更新时间
	TotalMoney  float64 `json:"total_money"`  // 总金额
	RealMoney   float64 `json:"real_money"`   // 实际到账金额
}
