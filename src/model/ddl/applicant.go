package ddl

/*
t_applicant
応募者
*/
type Applicant struct {
	AbstractTransactionModel
	// 媒体側ID
	OuterID string `json:"outer_id" gorm:"not null;check:outer_id <> '';type:varchar(255);index"`
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
	// コミットID
	CommitID string `json:"commit_id" gorm:"not null;check:commit_id <> '';type:text;index"`
	// 面接回数
	NumOfInterview uint `json:"num_of_interview"`
	// 書類通過フラグ
	DocumentPassFlg uint `json:"document_pass_flg"`
	// チームID
	TeamID uint64 `json:"team_id" gorm:"index"`
	// サイト(外部キー)
	Sites Site `gorm:"foreignKey:site_id;references:id"`
	// ステータス(外部キー)
	ApplicantStatus SelectStatus `gorm:"foreignKey:status;references:id"`
	// チーム(外部キー)
	Teams Team `gorm:"foreignKey:team_id;references:id"`
}

/*
t_applicant_user_association
応募者ユーザー紐づけ
*/
type ApplicantUserAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey;index"`
	// 表示フラグ
	DisplayFlg uint `json:"display_flg" gorm:"index"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_applicant_type
応募者種別
*/
type ApplicantType struct {
	AbstractTransactionModel
	// チームID
	TeamID uint64 `json:"team_id"`
	// 書類提出ルールID
	RuleID uint `json:"rule_id"`
	// 職種ID
	OccupationID uint `json:"occupation_id"`
	// 種別名
	Name string `json:"name" gorm:"not null;type:varchar(40)"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// 書類提出ルール(外部キー)
	DocumentRule DocumentRule `gorm:"foreignKey:rule_id;references:id"`
	// 職種(外部キー)
	Occupation Occupation `gorm:"foreignKey:occupation_id;references:id"`
}

/*
t_applicant_type_association
応募者種別紐づけ
*/
type ApplicantTypeAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// 種別ID
	TypeID uint64 `json:"type_id" gorm:"index"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// 種別(外部キー)
	Type ApplicantType `gorm:"foreignKey:type_id;references:id"`
}

/*
t_applicant_schedule_association
応募者面接予定紐づけ
*/
type ApplicantScheduleAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// 予定ID
	ScheduleID uint64 `json:"schedule_id" gorm:"index"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// ユーザー(外部キー)
	Schedule Schedule `gorm:"foreignKey:schedule_id;references:id"`
}

/*
t_applicant_resume_association
応募者履歴書紐づけ
*/
type ApplicantResumeAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// 拡張子
	Extension string `json:"extension" gorm:"not null;check:extension <> '';type:varchar(30)"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
}

/*
t_applicant_curriculum_vitae_association
応募者職務経歴書紐づけ
*/
type ApplicantCurriculumVitaeAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// 拡張子
	Extension string `json:"extension" gorm:"not null;check:extension <> '';type:varchar(30)"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
}

/*
t_applicant_url_association
応募者面接用URL紐づけ
*/
type ApplicantURLAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey;index"`
	// URL
	URL string `json:"url" gorm:"not null;check:url <> '';type:text"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
}

func (t Applicant) TableName() string {
	return "t_applicant"
}
func (t ApplicantUserAssociation) TableName() string {
	return "t_applicant_user_association"
}
func (t ApplicantType) TableName() string {
	return "t_applicant_type"
}
func (t ApplicantTypeAssociation) TableName() string {
	return "t_applicant_type_association"
}
func (t ApplicantScheduleAssociation) TableName() string {
	return "t_applicant_schedule_association"
}
func (t ApplicantResumeAssociation) TableName() string {
	return "t_applicant_resume_association"
}
func (t ApplicantCurriculumVitaeAssociation) TableName() string {
	return "t_applicant_curriculum_vitae_association"
}
func (t ApplicantURLAssociation) TableName() string {
	return "t_applicant_url_association"
}
