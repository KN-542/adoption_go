package request

import "api/src/model/ddl"

// 登録
type CreateUser struct {
	ddl.User
	// 所属チーム
	Teams []string `json:"teams"`
	// ロールハッシュキー
	RoleHashKey string `json:"role_hash_key"`
}

// 検索
type SearchUser struct {
	ddl.User
}

// 検索_同一企業
type SearchUserByCompany struct {
	ddl.User
}

// 取得
type GetUser struct {
	ddl.User
}

// チーム登録
type CreateTeam struct {
	Abstract
	ddl.Team
	// ユーザーリスト
	Users []string `json:"users"`
}

// チーム更新
type UpdateTeam struct {
	Abstract
	ddl.Team
	// ユーザーリスト
	Users []string `json:"users"`
}

// チーム基本情報更新
type UpdateBasicTeam struct {
	Abstract
	ddl.Team
}

// チーム削除
type DeleteTeam struct {
	Abstract
	ddl.Team
}

// チーム取得
type GetTeam struct {
	Abstract
	ddl.Team
}

// 自チーム取得
type GetOwnTeam struct {
	Abstract
}

// チーム検索
type SearchTeam struct {
	Abstract
	ddl.Team
}

// チーム検索_同一企業
type SearchTeamByCompany struct {
	ddl.User
}

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
	ddl.Schedule
}

// 予定削除
type DeleteSchedule struct {
	Abstract
	ddl.Schedule
}

// チーム毎ステータスイベント取得
type StatusEventsByTeam struct {
	Abstract
}
