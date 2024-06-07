package dto

import "api/src/model/request"

type UserSearch struct {
	request.UserSearch
	// チームID
	TeamID uint64 `json:"team_id"`
	// TODO
}
