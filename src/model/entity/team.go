package entity

import "api/src/model/ddl"

// Team
type Team struct {
	ddl.Team
	// ルールハッシュ
	RuleHash string `json:"rule_hash"`
	// 所属ユーザー
	Users []*ddl.User `json:"users" gorm:"many2many:t_team_association;foreignKey:id;joinForeignKey:team_id;References:id;joinReferences:user_id"`
}

// Team Search
type SearchTeam struct {
	ddl.Team
	// 所属ユーザー
	Users []*ddl.User `json:"users" gorm:"many2many:t_team_association;foreignKey:id;joinForeignKey:team_id;References:id;joinReferences:user_id"`
}

// Team Association
type TeamAssociation struct {
	ddl.TeamAssociation
}

// Team Event
type TeamEvent struct {
	ddl.TeamEvent
}

// Team Auto Assign Rule
type TeamAutoAssignRule struct {
	ddl.TeamAutoAssignRule
	HashKey string `json:"hash_key"`
}

// Team Assign Priority
type TeamAssignPriority struct {
	ddl.TeamAssignPriority
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 氏名
	Name string `json:"name"`
}

// Team Assign Priority only ID
type TeamAssignPriorityOnly struct {
	ddl.TeamAssignPriority
}

// Team Assign Possible
type TeamAssignPossible struct {
	ddl.TeamAssignPossible
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 氏名
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
}

// Team Per Interview
type TeamPerInterview struct {
	ddl.TeamPerInterview
}

// チーム毎イベント
type StatusEventsByTeam struct {
	// イベントハッシュキー
	EventHashKey string `json:"event_hash_key"`
	// 説明_日本語
	DescJa string `json:"desc_ja"`
	// 説明_英語
	DescEn string `json:"desc_en"`
	// 選考状況ハッシュキー
	SelectStatusHashKey string `json:"select_status_hash_key"`
	// ステータス名
	StatusName string `json:"status_name"`
}

type TeamEventEachInterview struct {
	ddl.TeamEventEachInterview
}

// チーム面接毎イベント
type InterviewEventsByTeam struct {
	// 面接回数
	NumOfInterview uint `json:"num_of_interview"`
	// 過程ハッシュ
	ProcessHash string `json:"process_hash"`
	// 過程
	Processing string `json:"processing"`
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
	// 選考状況ハッシュキー
	SelectStatusHashKey string `json:"select_status_hash_key"`
	// ステータス名
	StatusName string `json:"status_name"`
}

// 面接毎参加可能者予定取得
type AssignPossibleSchedule struct {
	// ユーザーID
	UserID uint64 `json:"user_id"`
	// ユーザーハッシュキー
	UserHashKey string `json:"user_hash_key"`
	// スケジュール
	Schedules []*Schedule `json:"schedules" gorm:"many2many:t_schedule_association;foreignKey:user_id;joinForeignKey:user_id;References:id;joinReferences:schedule_id"`
}
