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
