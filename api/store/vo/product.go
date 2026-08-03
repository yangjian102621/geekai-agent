package vo

type Product struct {
	Id        uint    `json:"id"`
	Name      string  `json:"name"`     // 名称
	Price     float64 `json:"price"`    // 价格
	Credit    int     `json:"credit"`   // 额度
	Enabled   bool    `json:"enabled"`  // 启用状态
	Sales     int     `json:"sales"`    // 销量
	SortNum   int8    `json:"sort_num"` // 排序
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
