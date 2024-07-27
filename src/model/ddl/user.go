package ddl

import (
	"time"
)

/*
t_user
ユーザー
*/
type User struct {
	AbstractTransactionModel
	// 氏名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(75);index"`
	// メールアドレス
	Email string `json:"email" gorm:"not null;type:varchar(100);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// パスワード(ハッシュ化)
	Password string `json:"password" gorm:"not null;check:password <> ''"`
	// 初回パスワード(ハッシュ化)
	InitPassword string `json:"init_password" gorm:"not null;check:init_password <> ''"`
	// ロールID
	RoleID uint64 `json:"role_id"`
	// ユーザー種別
	UserType uint `json:"user_type"`
	// リフレッシュトークン
	RefreshToken string `json:"refresh_token" gorm:"type:text"`
	// ロール(外部キー)
	Role CustomRole `gorm:"foreignKey:role_id;references:id"`
	// ログイン種別(外部キー)
	LoginType LoginType `gorm:"foreignKey:user_type;references:id"`
}

/*
t_team
チーム
*/
type Team struct {
	AbstractTransactionModel
	// チーム名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(30);index"`
	// 最大面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"check:num_of_interview >= 1 AND num_of_interview <= 30"`
}

/*
t_team_event
チームイベント
*/
type TeamEvent struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// イベントID
	EventID uint `json:"event_id" gorm:"primaryKey"`
	// ステータスID
	StatusID uint64 `json:"status_id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// イベント(外部キー)
	Event SelectStatusEvent `gorm:"foreignKey:event_id;references:id"`
	// ステータス(外部キー)
	Status SelectStatus `gorm:"foreignKey:status_id;references:id"`
}

/*
t_team_event_each_interview
チーム面接毎イベント
*/
type TeamEventEachInterview struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// 面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"primaryKey"`
	// ステータスID
	StatusID uint64 `json:"status_id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ステータス(外部キー)
	Status SelectStatus `gorm:"foreignKey:status_id;references:id"`
}

/*
t_team_association
チーム紐づけ
*/
type TeamAssociation struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_select_status
選考状況
*/
type SelectStatus struct {
	AbstractTransactionModel
	// チームID
	TeamID uint64 `json:"team_id"`
	// ステータス名
	StatusName string `json:"status_name" gorm:"not null;type:varchar(50)"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
}

/*
t_schedule
予定
*/
type Schedule struct {
	AbstractTransactionModel
	// タイトル
	Title string `json:"title" gorm:"not null;check:title <> '';type:varchar(30)"`
	// 頻度ID
	FreqID uint `json:"freq_id"`
	// 面接フラグ
	InterviewFlg uint `json:"interview_flg"`
	// 開始時刻
	Start time.Time `json:"start" gorm:"not null"`
	// 終了時刻
	End time.Time `json:"end" gorm:"not null"`
	// 頻度(外部キー)
	ScheduleFreqStatus ScheduleFreqStatus `gorm:"foreignKey:freq_id;references:id"`
}

/*
t_schedule_association
予定紐づけ
*/
type ScheduleAssociation struct {
	// 予定ID
	ScheduleID uint64 `json:"schedule_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// 予定(外部キー)
	Schedule Schedule `gorm:"foreignKey:schedule_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t User) TableName() string {
	return "t_user"
}
func (t Team) TableName() string {
	return "t_team"
}
func (t TeamEvent) TableName() string {
	return "t_team_event"
}
func (t TeamEventEachInterview) TableName() string {
	return "t_team_event_each_interview"
}
func (t TeamAssociation) TableName() string {
	return "t_team_association"
}
func (t SelectStatus) TableName() string {
	return "t_select_status"
}
func (t Schedule) TableName() string {
	return "t_schedule"
}
func (t ScheduleAssociation) TableName() string {
	return "t_schedule_association"
}
