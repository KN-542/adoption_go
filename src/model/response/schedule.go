package response

import "api/src/model/entity"

// 予定登録種別一覧
type SearchScheduleType struct {
	List []entity.ScheduleFreqStatus `json:"list"`
}

// 予定登録
type CreateSchedule struct {
	HashKey string `json:"hash_key"`
}

// 予定検索
type SearchSchedule struct {
	List []entity.Schedule `json:"list"`
}
