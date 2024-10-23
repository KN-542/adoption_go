package request

import "api/src/model/ddl"

// 登録
type CreateUser struct {
	ddl.User
	// 所属チーム
	Teams []string `json:"teams"`
	// ロールハッシュキー
	RoleHashKey string `json:"role_hash_key"`
}

// 検索
type SearchUser struct {
	ddl.User
}

// 検索_同一企業
type SearchUserByCompany struct {
	ddl.User
}

// 取得
type GetUser struct {
	ddl.User
}

// 削除
type DeleteUser struct {
	Abstract
	HashKeys []string `json:"hash_keys"`
}
