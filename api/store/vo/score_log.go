package vo

import "geekai/core/types"

type ScoreLog struct {
	Id        uint            `json:"id"`
	UserId    uint            `json:"user_id"`
	CreatorId uint            `json:"creator_id"`
	Username  string          `json:"username"`
	Type      types.ScoreType `json:"type"`
	TypeStr   string          `json:"type_str"`
	Amount    int             `json:"amount"`
	Mark      types.ScoreMark `json:"mark"`
	Balance   int             `json:"balance"`
	Subject   string          `json:"subject"`
	Remark    string          `json:"remark"`
	CreatedAt int64           `json:"created_at"`
}
