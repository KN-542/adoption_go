package ddl

import (
	"fmt"
	"time"
)

/*
t_applicant
応募者
*/
type Applicant struct {
	AbstractTransactionModel
	// 媒体側ID
	OuterID string `json:"outer_id" gorm:"not null;unique;check:outer_id <> '';type:varchar(255)"`
	// サイトID
	SiteID uint `json:"site_id" gorm:"index"`
	// ステータス
	Status uint64 `json:"status" gorm:"index"`
	// 氏名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(50);index"`
	// メールアドレス
	Email string `json:"email" gorm:"not null;type:varchar(255);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// TEL
	Tel string `json:"tel" gorm:"type:varchar(20);check:tel ~ '^[0-9]{0,20}$'"`
	// 年齢
	Age uint `json:"age" gorm:"check:(age >= 18 AND age <= 100) OR age = 0;index"`
	// 履歴書
	Resume string `json:"resume" gorm:"type:varchar(255);index"`
	// 職務経歴書
	CurriculumVitae string `json:"curriculum_vitae" gorm:"type:varchar(255);index"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url" gorm:"type:text"`
	// チームID
	TeamID uint64 `json:"team_id"`
	// カレンダーID
	CalendarID uint64 `json:"calendar_id"`
	// サイト(外部キー)
	Site Site `gorm:"foreignKey:site_id;references:id"`
	// ステータス(外部キー)
	ApplicantStatus SelectStatus `gorm:"foreignKey:status;references:id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
}

/*
t_applicant_user_association
応募者ユーザー紐づけ
*/
type ApplicantUserAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

func (t Applicant) TableName() string {
	year := time.Now().Year()
	month := time.Now().Month()
	return fmt.Sprintf("t_applicant_%d_%02d", year, month)
}
func (t ApplicantUserAssociation) TableName() string {
	year := time.Now().Year()
	month := time.Now().Month()
	return fmt.Sprintf("t_applicant_user_association_%d_%02d", year, month)
}

type GetOauthURLResponse struct {
	Url string `json:"url"`
}

/*
	txt、csvダウンロード用
*/
// 応募者ダウンロード
type ApplicantsDownload struct {
	Values [][]string `json:"values"`
	Site   int        `json:"site"`
}

type ApplicantDesired struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 希望面接日時
	DesiredAt time.Time `json:"desired_at"`
	// タイトル
	Title string `json:"title"`
	// カレンダーハッシュキー
	CalendarHashKey string `json:"calendar_hash_key"`
}

type ApplicantAndUser struct {
	Applicant   Applicant `json:"applicant"`
	UserHashKey string    `json:"user_hash_key"`
	Code        string    `json:"code"`
}

// ファイルアップロード
type FileUpload struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// ファイル拡張子
	Extension string `json:"extension"`
	// ファイル名(Pre)
	NamePre string `json:"name_pre"`
}

// ファイルダウンロード
type FileDownload struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// ファイル名(Pre)
	NamePre string `json:"name_pre"`
}
