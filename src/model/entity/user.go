package entity

import "api/src/model/ddl"

// Login
type Login struct {
	ddl.User
}

// User
type User struct {
	ddl.User
}

// Team
type Team struct {
	ddl.Team
}
