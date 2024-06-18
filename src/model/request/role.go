package request

// ロールチェック
type CheckRole struct {
	Abstract
	// 該当ロールID
	ID uint `json:"id"`
}
