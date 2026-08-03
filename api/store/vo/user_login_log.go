package vo

type UserLoginLog struct {
	Id           uint   `json:"id"`
	UserId       uint   `json:"user_id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	LoginIp      string `json:"login_ip"`
	LoginAddress string `json:"login_address"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
