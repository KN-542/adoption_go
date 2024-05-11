package entity

import "api/src/model/ddl"

// ロールマスタ
type Role struct {
	ddl.Role
}

// ロール
type CustomRole struct {
	ddl.CustomRole
}

// 付与ロール
type RoleAssociation struct {
	ddl.RoleAssociation
}
