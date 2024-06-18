package dto

import (
	"api/src/model/request"
	"time"
)

type SearchApplicant struct {
	request.SearchApplicant
	// ユーザー
	UserIDs []uint64
}

// 予約表サブ
type ReserveTableSub struct {
	// 時間
	Time time.Time `json:"time"`
	// 予約可否
	IsReserve bool `json:"is_reserve"`
}
