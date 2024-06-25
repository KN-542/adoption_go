package response

import "api/src/model/entity"

// 登録
type CreateCompany struct {
	// パスワード
	Password string `json:"password"`
}

// 検索
type SearchCompany struct {
	List []entity.Company `json:"list"`
}
