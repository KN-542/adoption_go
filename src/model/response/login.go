package response

import (
	"api/src/model/entity"
)

// Login
type Login struct {
	entity.User
}

// MFA
type MFA struct {
	// 遷移パス
	Path string `json:"path"`
	// パスワード変更_必要性
	IsPasswordChange bool `json:"is_password_change"`
}

// GetLoginType
type GetLoginType struct {
	// ログイン種別
	LoginType uint `json:"login_type"`
}
