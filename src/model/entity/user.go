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
	// ロール名
	RoleName string `json:"role_name"`
}

// User Refresh Token Association
type UserRefreshTokenAssociation struct {
	ddl.UserRefreshTokenAssociation
}
