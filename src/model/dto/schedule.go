package dto

import "api/src/model/ddl"

// ユーザー単位予定取得
type GetScheduleByUser struct {
	ddl.ScheduleAssociation
	// 除外予定ハッシュリスト
	RemoveScheduleHashKeys []string `json:"remove_schedule_hash_keys"`
}
