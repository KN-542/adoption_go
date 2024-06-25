package request

import "api/src/model/ddl"

// 登録
type CreateCompany struct {
	Abstract
	ddl.Company
	// メールアドレス
	Email string `json:"email"`
}

// 検索
type SearchCompany struct {
	Abstract
	ddl.Company
}
