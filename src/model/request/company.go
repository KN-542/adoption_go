package request

import "api/src/model/ddl"

// 登録
type CompanyCreate struct {
	Abstract
	ddl.Company
	// メールアドレス
	Email string `json:"email"`
}
