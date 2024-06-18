package entity

import "api/src/model/ddl"

// Login
type Login struct {
	ddl.User
	// チームID
	TeamID uint64 `json:"team_id"`
}

// User
type User struct {
	ddl.User
}

// Search
type SearchUser struct {
	ddl.User
	// TODO
}

// Team
type Team struct {
	ddl.Team
}

// Team Search
type SearchTeam struct {
	ddl.Team
	// 所属ユーザー
	Users []string `json:"users"`
}

// UserSchedule
type UserSchedule struct {
	ddl.UserSchedule
}

// Team Association
type TeamAssociation struct {
	ddl.TeamAssociation
}

// UserScheduleAssociation
type UserScheduleAssociation struct {
	ddl.UserScheduleAssociation
}
