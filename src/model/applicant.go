package model

import "time"

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
	Status uint `json:"status" gorm:"index"`
	// 氏名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(50);index"`
	// メールアドレス
	Email string `json:"email" gorm:"not null;unique;type:varchar(255);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$';index"`
	// TEL
	Tel string `json:"tel" gorm:"type:varchar(20);check:tel ~ '^[0-9]{0,20}$'"`
	// 年齢
	Age int `json:"age" gorm:"check:(age >= 18 AND age <= 100) OR age = -1;index"`
	// 履歴書
	Resume string `json:"resume" gorm:"type:varchar(255);index"`
	// 職務経歴書
	CurriculumVitae string `json:"curriculum_vitae" gorm:"type:varchar(255);index"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url" gorm:"type:text"`
	// カレンダーID
	CalendarID uint `json:"calendar_id"`
	// サイト(外部キー)
	Site Site `gorm:"foreignKey:site_id;references:id"`
	// ステータス(外部キー)
	ApplicantStatus ApplicantStatus `gorm:"foreignKey:status;references:id"`
	// スケジュール(外部キー)
	Schedule UserSchedule `gorm:"foreignKey:calendar_id;references:id"`
}

func (t Applicant) TableName() string {
	return "t_applicant"
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

// 応募者ダウンロード Response
type ApplicantsDownloadResponse struct {
	Applicants []ApplicantWith `json:"applicants"`
}

type ApplicantWith struct {
	Applicant
	StatusNameJa string    `json:"status_name_ja"`
	SiteNameJa   string    `json:"site_name_ja"`
	UserNames    string    `json:"user_names"`
	Start        time.Time `json:"start"`
}

type ApplicantSearchRequest struct {
	// サイトID
	SiteIDList []uint `json:"site_id_list"`
	// 応募者ステータス
	ApplicantStatusList []uint `json:"applicant_status_list"`
	// 履歴書
	Resume uint `json:"resume"`
	// 職務経歴書
	CurriculumVitae uint `json:"curriculum_vitae"`
	// 氏名
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
	// 面接官
	Users string `json:"users"`
	// ソート(key)
	SortKey string `json:"sort_key"`
	// ソート(向き)
	SortAsc bool `json:"sort_asc"`
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
