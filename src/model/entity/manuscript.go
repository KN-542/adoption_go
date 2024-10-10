package entity

import "api/src/model/ddl"

type Manuscript struct {
	ddl.Manuscript
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
