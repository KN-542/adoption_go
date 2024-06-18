package response

import "api/src/model/entity"

// 登録
type CreateUser struct {
	entity.User
}

// 検索
type SearchUser struct {
	List []entity.SearchUser `json:"list"`
}

// 取得
type GetUser struct {
	entity.User
}

// チーム検索
type SearchTeam struct {
	List []entity.SearchTeam `json:"list"`
}

// 予定登録種別一覧
type SearchScheduleType struct {
	List []entity.CalendarFreqStatus `json:"list"`
}

// 予定登録
type CreateSchedule struct {
	HashKey string `json:"hash_key"`
}

// 予定検索
type SearchSchedule struct {
	List []entity.UserSchedule `json:"list"`
}
