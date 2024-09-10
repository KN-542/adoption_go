package response

import "api/src/model/entity"

// 登録
type CreateUser struct {
	entity.User
}

// 検索
type SearchUser struct {
	List []entity.SearchUser `json:"list"`
}

// 検索_同一企業
type SearchUserByCompany struct {
	List []entity.SearchUser `json:"list"`
}

// 取得
type GetUser struct {
	entity.User
}

// ステータスイベントマスタ一覧
type ListStatusEvent struct {
	List []entity.SelectStatusEvent `json:"list"`
}

// アサイン関連マスタ取得
type AssignMaster struct {
	Rule     []entity.AssignRule     `json:"rule"`
	AutoRule []entity.AutoAssignRule `json:"auto_rule"`
}

// 書類提出ルールマスタ取得
type DocumentRule struct {
	List []entity.DocumentRule `json:"list"`
}

// 職種マスタ取得
type Occupation struct {
	List []entity.Occupation `json:"list"`
}
