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

// Team
type Team struct {
	ddl.Team
	// 所属ユーザー
	Users []*ddl.User `json:"users" gorm:"many2many:t_team_association;foreignKey:id;joinForeignKey:team_id;References:id;joinReferences:user_id"`
}

// Team Search
type SearchTeam struct {
	ddl.Team
	// 所属ユーザー
	Users []*ddl.User `json:"users" gorm:"many2many:t_team_association;foreignKey:id;joinForeignKey:team_id;References:id;joinReferences:user_id"`
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
