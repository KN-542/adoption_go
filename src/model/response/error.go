package response

type Error struct {
	// スタータス(ヘッダー)
	Status int `json:"status"`
	// コード
	Code uint `json:"code"` // 必要な場合のみ
	// エラー
	Error error `json:"error"` // 消す可能性大
}

type ErrorCode struct {
	// コード
	Code uint `json:"code"`
}

func ErrorConvert(e Error) ErrorCode {
	return ErrorCode{Code: e.Code}
}
