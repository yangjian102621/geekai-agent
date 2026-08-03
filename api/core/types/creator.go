package types

const (
	WithdrawMethodAlipay = "alipay"
	WithdrawMethodWxPay  = "wxpay"

	// 提现状态
	WithdrawStatusPending = "pending" // 待处理
	WithdrawStatusSuccess = "success" // 已结算
	WithdrawStatusReject  = "reject"  // 已拒绝

	ScoreToRMBRatio = 1000 // 积分兑换比例
)

type CreatorScoreType string

const (
	CreatorScoreTypeIncome   CreatorScoreType = "income"    // 收入
	CreatorScoreTypeWithdraw CreatorScoreType = "withdraw"  // 提现
	CreatorScoreTypeRefund   CreatorScoreType = "refund"    // 退款
	CreatorScoreFineTune     CreatorScoreType = "fine_tune" // 系统调整
)

func (t CreatorScoreType) String() string {
	switch t {
	case CreatorScoreTypeIncome:
		return "收入"
	case CreatorScoreTypeWithdraw:
		return "提现"
	case CreatorScoreTypeRefund:
		return "退款"
	case CreatorScoreFineTune:
		return "系统调整"
	}
	return "其他"
}
