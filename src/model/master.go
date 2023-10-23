package model

// m_site
type Site struct {
	SiteID   int    `json:"site_id" gorm:"primaryKey;check:site_id >= 1 AND site_id <= 10"`
	SiteName string `json:"site_name" gorm:"type:varchar(20)"`
}

func (m Site) TableName() string {
	return "m_site"
}

// m_role
type Role struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// ロール名
	Name string `json:"name" gorm:"unique;type:varchar(20)"`
}

func (m Role) TableName() string {
	return "m_role"
}
