package vo

type AdminUser struct {
	Id          uint   `json:"id"`
	Username    string `json:"username"`
	Status      bool   `json:"status"`        // 当前状态
	LastLoginAt int64  `json:"last_login_at"` // 最后登录时间
	LastLoginIp string `json:"last_login_ip"` // 最后登录 IP
	RoleIds     any    `json:"role_ids"`      //角色ids
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
