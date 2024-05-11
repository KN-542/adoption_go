package request

import "api/src/model/ddl"

// Login
type Login struct {
	ddl.User
}

// CodeGenerate
type CodeGenerate struct {
	ddl.User
}

// MFA
type MFA struct {
	ddl.User
	// 認証コード
	Code string `json:"code"`
}

// JWTDecode
type JWTDecode struct {
	ddl.User
}

// PasswordChange
type PasswordChange struct {
	ddl.User
}

// Logout
type Logout struct {
	ddl.User
}
