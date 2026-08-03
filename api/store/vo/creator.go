package vo

// Creator 创作者VO
type Creator struct {
	Id          uint   `json:"id"`          // 主键ID
	UserId      uint   `json:"user_id"`     // 关联用户ID
	Username    string `json:"username"`    // 关联用户名
	Fee         int    `json:"fee"`         // 提现费率(0-100)
	Name        string `json:"name"`        // 创作者名称
	Description string `json:"description"` // 创作者简介
	Logo        string `json:"logo"`        // 创作者Logo
	Enabled     bool   `json:"enabled"`     // 是否启用
	Scores      int    `json:"scores"`      // 积分
	CreatedAt   int64  `json:"created_at"`  // 创建时间
	UpdatedAt   int64  `json:"updated_at"`  // 更新时间
	Check       int8   `json:"check"`       // 审核状态
	CheckNote   string `json:"check_note"`  // 审核备注

	AppCount        int            `json:"app_count"`        // 应用数量
	TotalEarnings   int            `json:"total_earnings"`   // 总收益
	TodayEarnings   int            `json:"today_earnings"`   // 今日收益
	WithdrawConfigs WithdrawConfig `json:"withdraw_configs"` // 提现配置
}

// WithdrawConfig 提现配置
type WithdrawConfig struct {
	Name            string `json:"name"`               // 真实姓名
	Mobile          string `json:"mobile"`             // 联系手机
	Method          string `json:"method"`             // 收款方式(alipay,wx_pay)
	Account         string `json:"account"`            // 收款账号
	Qrcode          string `json:"qrcode"`             // 收款二维码
	ScoreToRMBRatio int    `json:"score_to_rmb_ratio"` // 积分兑换比例
}
