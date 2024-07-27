package ddl

/*
t_applicant
応募者
*/
type Applicant struct {
	AbstractTransactionModel
	// 媒体側ID
	OuterID string `json:"outer_id" gorm:"not null;check:outer_id <> '';type:varchar(255)"`
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
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_applicant_schedule_association
応募者面接予定紐づけ
*/
type ApplicantScheduleAssociation struct {
	// 応募者ID
	ApplicantID uint64 `json:"applicant_id" gorm:"primaryKey"`
	// 予定ID
	ScheduleID uint64 `json:"schedule_id"`
	// 応募者(外部キー)
	Applicant Applicant `gorm:"foreignKey:applicant_id;references:id"`
	// ユーザー(外部キー)
	Schedule Schedule `gorm:"foreignKey:schedule_id;references:id"`
}

func (t Applicant) TableName() string {
	return "t_applicant"
}
func (t ApplicantUserAssociation) TableName() string {
	return "t_applicant_user_association"
}
func (t ApplicantScheduleAssociation) TableName() string {
	return "t_applicant_schedule_association"
}
