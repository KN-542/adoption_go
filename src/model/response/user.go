package response

import "api/src/model/entity"

// 登録
type UserCreate struct {
	entity.User
}

// 検索
type UserSearch struct {
	List []entity.UserSearch `json:"list"`
}
