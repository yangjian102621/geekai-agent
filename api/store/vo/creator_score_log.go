package vo

import (
	"geekai/core/types"
)

// CreatorScoreLog 创作者VO
type CreatorScoreLog struct {
	Id        uint                   `json:"id"`         // 主键ID
	UserId    uint                   `json:"user_id"`    // 关联用户ID
	CreatorId uint                   `json:"creator_id"` // 关联创作者ID
	AppId     uint                   `json:"app_id"`     // 应用ID
	Type      types.CreatorScoreType `json:"type"`       // 类型
	Score     int                    `json:"score"`      // 积分
	Balance   int                    `json:"balance"`    // 余额
	Subject   string                 `json:"subject"`    // 主题
	Remark    string                 `json:"remark"`     // 备注
	Mark      types.ScoreMark        `json:"mark"`       // 资金类型
	CreatedAt int64                  `json:"created_at"` // 创建时间
}
