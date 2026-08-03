package vo

type AppCategory struct {
	Id        uint   `json:"id"`
	CreatorId uint   `json:"creator_id"` // 创作者ID
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
