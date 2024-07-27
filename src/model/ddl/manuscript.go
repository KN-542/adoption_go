package ddl

/*
t_manuscript
原稿
*/
type Manuscript struct {
	AbstractTransactionModel
	// 原稿内容
	Content string `json:"content" gorm:"not null;check:content <> '';type:text;index"`
}

/*
t_manuscript_team_association
原稿チーム紐づけ
*/
type ManuscriptTeamAssociation struct {
	// 原稿ID
	ManuscriptID uint64 `json:"manuscript_id" gorm:"primaryKey"`
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// 原稿(外部キー)
	Manuscript Manuscript `gorm:"foreignKey:manuscript_id;references:id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
}

/*
t_manuscript_site_association
原稿サイト紐づけ
*/
type ManuscriptSiteAssociation struct {
	// 原稿ID
	ManuscriptID uint64 `json:"manuscript_id" gorm:"primaryKey"`
	// サイトID
	SiteID uint64 `json:"site_id" gorm:"primaryKey"`
	// 原稿(外部キー)
	Manuscript Manuscript `gorm:"foreignKey:manuscript_id;references:id"`
	// サイト(外部キー)
	Site Site `gorm:"foreignKey:site_id;references:id"`
}

func (t Manuscript) TableName() string {
	return "t_manuscript"
}
func (t ManuscriptTeamAssociation) TableName() string {
	return "t_manuscript_team_association"
}
func (t ManuscriptSiteAssociation) TableName() string {
	return "t_manuscript_site_association"
}
