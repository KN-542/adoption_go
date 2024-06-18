package dto

import "api/src/model/request"

type SearchUser struct {
	request.SearchUser
	// チームID
	TeamID uint64 `json:"team_id"`
	// TODO
}
