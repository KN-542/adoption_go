package model

type ErrorResponse struct {
	// スタータス(ヘッダー)
	Status int `json:"status"`
	// コード
	Code int8 `json:"code"` // 必要な場合のみ
	// エラー
	Error error `json:"error"` // 消す可能性大
}

type ErrorCodeResponse struct {
	// コード
	Code int8 `json:"code"`
}

func ErrorConvert(e ErrorResponse) ErrorCodeResponse {
	return ErrorCodeResponse{Code: e.Code}
}
