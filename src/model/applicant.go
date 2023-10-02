package model

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
	// サイトID
	SiteID int `json:"site_id" gorm:"check:site_id >= 1 AND site_id <= 10"`
	// 氏名
	Name string `json:"name" gorm:"type:varchar(50)"`
	// メールアドレス
	Email string `json:"email" gorm:"type:varchar(255);check:email ~ '^[a-zA-Z0-9_+-]+(\\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\\.)+[a-zA-Z]{2,}$'"`
	// TEL
	Tel string `json:"tel" gorm:"type:varchar(20);check:tel ~ '^[0-9]{0,20}$'"`
	// 年齢
	Age int `json:"age" gorm:"check:(age >= 18 AND age <= 100) OR age = -1"`
	// TODO リレーション
	// Site Site `gorm:"foreignKey:site_id"`
}

func (t Applicant) TableName() string {
	return "t_applicant"
}
