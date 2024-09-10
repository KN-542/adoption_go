package request

import "api/src/model/ddl"

// 予定登録
type CreateSchedule struct {
	Abstract
	ddl.Schedule
	// ユーザーリスト
	Users []string `json:"users"`
}

// 予定更新
type UpdateSchedule struct {
	Abstract
	ddl.Schedule
	// ユーザーリスト
	Users []string `json:"users"`
}

// 予定検索
type SearchSchedule struct {
	Abstract
}

// 予定削除
type DeleteSchedule struct {
	Abstract
	ddl.Schedule
}
