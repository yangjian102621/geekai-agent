package types

type PaymentConfig struct {
	AlipayConfig AlipayConfig `json:"alipay"` // 支付宝支付渠道配置
	EPayConfig   EPayConfig   `json:"epay"`   // 易支付配置
	WxPayConfig  WxPayConfig  `json:"wxpay"`  // 微信支付渠道配置
}

// AlipayConfig 支付宝支付配置
type AlipayConfig struct {
	Enabled         bool   `json:"enabled"`           // 是否启用该支付通道
	AppId           string `json:"app_id"`            // 应用 ID
	PrivateKey      string `json:"private_key"`       // 应用私钥
	AlipayPublicKey string `json:"alipay_public_key"` // 支付宝公钥
	Domain          string `json:"domain"`            // 支付回调域名
}

// WxPayConfig 微信支付配置
type WxPayConfig struct {
	Enabled    bool   `json:"enabled"`     // 是否启用该支付通道
	AppId      string `json:"app_id"`      // 公众号的APPID,如：wxd678efh567hg6787
	MchId      string `json:"mch_id"`      // 直连商户的商户号，由微信支付生成并下发
	SerialNo   string `json:"serial_no"`   // 商户证书的证书序列号
	PrivateKey string `json:"private_key"` // 商户证书私钥
	ApiV3Key   string `json:"api_v3_key"`  // API V3 秘钥
	Domain     string `json:"domain"`      // 支付回调域名
}

// EPayConfig 易支付配置
type EPayConfig struct {
	Enabled    bool   `json:"enabled"`     // 是否启用该支付通道
	AppId      string `json:"app_id"`      // 商户 ID
	PrivateKey string `json:"private_key"` // 私钥
	ApiURL     string `json:"api_url"`     // z支付 API 网关
	Domain     string `json:"domain"`      // 支付回调域名
}
