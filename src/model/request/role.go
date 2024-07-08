package request

// ロールチェック
type CheckRole struct {
	Abstract
	// 該当ロールID
	ID uint `json:"id"`
}

// ロール検索_企業ID
type SearchRoleByComapny struct {
	Abstract
}
