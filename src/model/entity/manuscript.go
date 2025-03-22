package entity

import "api/src/model/ddl"

type Manuscript struct {
	ddl.Manuscript
	// 原稿と紐づくチーム
	Teams []*ddl.Team `json:"teams"  gorm:"many2many:t_manuscript_team_association;foreignKey:id;joinForeignKey:manuscript_id;References:id;joinReferences:team_id"`
	//　原稿と紐づくサイト
	Sites []*ddl.Site `json:"sites"  gorm:"many2many:t_manuscript_site_association;foreignKey:id;joinForeignKey:manuscript_id;References:id;joinReferences:site_id"`
}

type SearchManuscript struct {
	ddl.Manuscript
	// サイト
	Sites []ddl.Site `json:"sites"  gorm:"many2many:t_manuscript_site_association;foreignKey:id;joinForeignKey:manuscript_id;References:id;joinReferences:site_id"`
}

type ManuscriptTeamAssociation struct {
	ddl.ManuscriptTeamAssociation
}
type ManuscriptSiteAssociation struct {
	ddl.ManuscriptSiteAssociation
}
