package vo

type User struct {
	Id          uint   `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Salt        string `json:"salt"`    // 密码盐
	Scores      int    `json:"scores"`  // 剩余积分
	Enabled     bool   `json:"enabled"` // 当前状态
	ExpiredTime int64  `json:"expired_time"`
	LastLoginAt int64  `json:"last_login_at"` // 最后登录时间
	LastLoginIp string `json:"last_login_ip"` // 最后登录 IP
	Vip         bool   `json:"vip"`
	OpenId      string `json:"openid"`      // 第三方登录 OpenID
	Platform    string `json:"platform"`    // 第三方登录平台
	Invitor     uint   `json:"invitor"`     // 邀请人ID
	InviteCode  string `json:"invite_code"` // 邀请码
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
