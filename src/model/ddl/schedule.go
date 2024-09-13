package ddl

import (
	"time"
)

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
	Start time.Time `json:"start" gorm:"not null;index"`
	// 終了時刻
	End time.Time `json:"end" gorm:"not null"`
	// チームID
	TeamID uint64 `json:"team_id"`
	// 頻度(外部キー)
	ScheduleFreqStatus ScheduleFreqStatus `gorm:"foreignKey:freq_id;references:id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
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

func (t Schedule) TableName() string {
	return "t_schedule"
}
func (t ScheduleAssociation) TableName() string {
	return "t_schedule_association"
}
