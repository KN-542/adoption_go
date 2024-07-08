package response

import "api/src/model/entity"

// ロール検索_企業ID
type SearchRoleByComapny struct {
	List []entity.CustomRole `json:"list"`
}
