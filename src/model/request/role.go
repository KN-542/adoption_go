package request

// ロールチェック
type RoleCheck struct {
	Abstract
	// 該当ロールID
	ID uint `json:"id"`
}
