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

// 検索_同一企業
type SearchUserByCompany struct {
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

// チーム取得
type GetTeam struct {
	entity.Team
}

// チーム取得
type GetOwnTeam struct {
	Team   entity.Team                    `json:"team"`
	Events []entity.InterviewEventsByTeam `json:"events"`
}

// チーム検索_同一企業
type SearchTeamByCompany struct {
	List []entity.SearchTeam `json:"list"`
}

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

// ステータスイベントマスタ一覧
type ListStatusEvent struct {
	List []entity.SelectStatusEvent `json:"list"`
}

// チーム毎ステータスイベント取得
type StatusEventsByTeam struct {
	List []entity.StatusEventsByTeam `json:"list"`
}
