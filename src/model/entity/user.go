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
	// 所属チーム
	Teams []*ddl.Team `json:"teams" gorm:"many2many:t_team_association;foreignKey:id;joinForeignKey:user_id;References:id;joinReferences:team_id"`
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
