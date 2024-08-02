package entity

import "api/src/model/ddl"

// Login
type Login struct {
	ddl.User
	// チームID
	TeamID uint64 `json:"team_id"`
}

// User
type User struct {
	ddl.User
}

// Search
type SearchUser struct {
	ddl.User
	// ロール名
	RoleName string `json:"role_name"`
}

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

// Schedule
type Schedule struct {
	ddl.Schedule
}

// Team Association
type TeamAssociation struct {
	ddl.TeamAssociation
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

// ScheduleAssociation
type ScheduleAssociation struct {
	ddl.ScheduleAssociation
}

// チーム毎イベント
type StatusEventsByTeam struct {
	// イベントハッシュキー
	EventHashKey string `json:"event_hash_key"`
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
	// 選考状況ハッシュキー
	SelectStatusHashKey string `json:"select_status_hash_key"`
	// ステータス名
	StatusName string `json:"status_name"`
}

// チーム面接毎イベント
type InterviewEventsByTeam struct {
	// 面接回数
	NumOfInterview uint `json:"num_of_interview"`
	// 選考状況ハッシュキー
	SelectStatusHashKey string `json:"select_status_hash_key"`
	// ステータス名
	StatusName string `json:"status_name"`
}
