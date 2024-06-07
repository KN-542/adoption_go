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
type UserSearch struct {
	ddl.User
	// TODO
}

// Team
type Team struct {
	ddl.Team
}

// Team Association
type TeamAssociation struct {
	ddl.TeamAssociation
}
