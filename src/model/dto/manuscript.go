package dto

import "api/src/model/request"

// 検索
type SearchManuscript struct {
	request.SearchManuscript
	// チームID
	TeamID uint64 `json:"team_id"`
}
