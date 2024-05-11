package response

import "api/src/model/entity"

// Sidebar
type Sidebar struct {
	Sidebars []entity.Sidebar `json:"sidebars"`
}

// Roles
type Roles struct {
	Map map[string]bool `json:"map"`
}
