package model

import "time"

/*
	OAuth2.0用(削除予定)
*/
type ApplicantResponse struct {
	ID    string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name  string `json:"name" gorm:"notNull;type:varchar(50)"`
	Email string `json:"email" gorm:"notNull;type:varchar(50)"`
}

type ApplicantSearch struct {
	Code            string `json:"code"`
	StartCellRow    int    `json:"start_cell_row"`
	EndCellRow      int    `json:"end_cell_row"`
	StartCellColumn string `json:"start_cell_column"`
	EndCellColumn   string `json:"end_cell_column"`
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
	Applicants []Applicant `json:"applicants"`
}

type Applicant struct {
	// ID
	ID string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	// ハッシュキー
	HashKey string `json:"hash_key" gorm:"unique;type:text"`
	// サイトID
	SiteID int `json:"site_id"`
	// 氏名
	Name string `json:"name" gorm:"type:varchar(50)"`
	// メールアドレス
	Email string `json:"email" gorm:"type:varchar(255);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$'"`
	// TEL
	Tel string `json:"tel" gorm:"type:varchar(20);check:tel ~ '^[0-9]{0,20}$'"`
	// 年齢
	Age int `json:"age" gorm:"check:(age >= 18 AND age <= 100) OR age = -1"`
	// 履歴書
	Resume string `json:"resume" gorm:"type:varchar(255)"`
	// 職務経歴書
	CurriculumVitae string `json:"curriculum_vitae" gorm:"type:varchar(255)"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url" gorm:"type:text"`
	// 希望面接日時
	DesiredAt string `json:"desired_at" gorm:"type:text"`
	// 登録日時
	CreatedAt time.Time `json:"created_at"`
	// 更新日時
	UpdatedAt time.Time `json:"updated_at"`
	// サイト(外部キー)
	Site Site `gorm:"foreignKey:site_id;references:id"`
}

func (t Applicant) TableName() string {
	return "t_applicant"
}

type ApplicantDesired struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 希望面接日時
	DesiredAt []string `json:"desired_at"`
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
