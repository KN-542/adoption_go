package model

// m_site
type Site struct {
	// ID
	SiteID int `json:"site_id" gorm:"primaryKey;check:site_id >= 1 AND site_id <= 10"`
	// 媒体名_日本語
	SiteNameJa string `json:"site_name" gorm:"type:varchar(20)"`
}

func (m Site) TableName() string {
	return "m_site"
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
