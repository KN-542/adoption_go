package request

import "api/src/model/ddl"

// Sidebar
type Sidebar struct {
	ddl.User
}

// Roles
type Roles struct {
	ddl.User
}

// Teams
type TeamsBelong struct {
	ddl.User
}
