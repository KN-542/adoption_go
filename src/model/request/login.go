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

// GetLoginType
type GetLoginType struct {
	ddl.User
}

// LoginApplicant
type LoginApplicant struct {
	ddl.Applicant
	// チームハッシュキー
	TeamHashKey string `json:"team_hash_key"`
}

// CodeGenerateApplicant
type CodeGenerateApplicant struct {
	ddl.Applicant
}

// MFAApplicant
type MFAApplicant struct {
	ddl.Applicant
	// 認証コード
	Code string `json:"code"`
}

// JWTDecodeApplicant
type JWTDecodeApplicant struct {
	ddl.Applicant
}

// LogoutApplicant
type LogoutApplicant struct {
	ddl.Applicant
}
