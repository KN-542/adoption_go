package response

import "api/src/model/entity"

// 検索
type SearchManuscript struct {
	List []entity.SearchManuscript `json:"list"`
	// 総数
	Num uint64 `json:"num"`
}

// 検索_同一チーム
type SearchManuscriptByTeam struct {
	List []entity.Manuscript `json:"list"`
}
