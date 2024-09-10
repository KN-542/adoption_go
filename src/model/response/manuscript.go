package response

import "api/src/model/entity"

// 検索
type SearchManuscript struct {
	List []entity.SearchManuscript `json:"list"`
	// 総数
	Num uint64 `json:"num"`
}
