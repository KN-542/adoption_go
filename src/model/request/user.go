package request

import "api/src/model/ddl"

// 登録
type UserCreate struct {
	ddl.User
	// 所属チーム
	Teams []string `json:"teams"`
	// ロールハッシュキー
	RoleHashKey string `json:"role_hash_key"`
}

// チーム所属ユーザー一覧
type UserSearch struct {
	ddl.User
}
