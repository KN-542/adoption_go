package request

// 検索
type SearchManuscript struct {
	Abstract
	// ページ
	Page int `json:"page"`
	// ページサイズ
	PageSize int `json:"page_size"`
	// サイト一覧
	Sites []string `json:"sites"`
}

// 登録
type CreateManuscript struct {
	Abstract
	// 内容
	Content string `json:"content"`
	// 使用可能チーム
	Teams []string `json:"teams"`
	// 使用可能サイト
	Sites []string `json:"sites"`
}

// 応募者紐づけ登録
type CreateApplicantAssociation struct {
	Abstract
	// ハッシュキー
	ManuscriptHash string `json:"manuscript_hash"`
	// 応募者
	Applicants []string `json:"applicants"`
}

// 検索_同一チーム
type SearchManuscriptByTeam struct {
	Abstract
}

// 削除リクエスト
type DeleteManuscriptRequest struct {
	UserHashKey        string   `json:"user_hash_key"`
	ManuscriptHashKeys []string `json:"manuscript_hash_keys"`
}
