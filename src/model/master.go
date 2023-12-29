package model

// m_site
type Site struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 媒体名_日本語
	SiteNameJa string `json:"site_name_ja" gorm:"type:varchar(20)"`
}

func (m Site) TableName() string {
	return "m_site"
}

type Sites struct {
	List []Site `json:"list"`
}

// m_role
type Role struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// ロール名_日本語
	NameJa string `json:"name_ja" gorm:"unique;type:varchar(20)"`
}

func (m Role) TableName() string {
	return "m_role"
}

// m_applicant_status
type ApplicantStatus struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// ステータス名_日本語
	StatusNameJa string `json:"status_name_ja" gorm:"unique;type:varchar(20)"`
}

func (m ApplicantStatus) TableName() string {
	return "m_applicant_status"
}

type ApplicantStatusList struct {
	List []ApplicantStatus `json:"list"`
}

// m_calendar_freq_status
type CalendarFreqStatus struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 頻度
	Freq string `json:"freq" gorm:"unique;type:varchar(10)"`
}

func (m CalendarFreqStatus) TableName() string {
	return "m_calendar_freq_status"
}
