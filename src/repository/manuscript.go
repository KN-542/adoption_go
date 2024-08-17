package repository

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"log"

	"gorm.io/gorm"
)

type IManuscriptRepository interface {
	// 紐づけ取得
	GetAssociationByTeamID(m *ddl.ManuscriptTeamAssociation) ([]entity.ManuscriptTeamAssociation, error)
}

type ManuscriptRepository struct {
	db *gorm.DB
}

func NewManuscriptRepository(db *gorm.DB) IManuscriptRepository {
	return &ManuscriptRepository{db}
}

// 紐づけ取得
func (s *ManuscriptRepository) GetAssociationByTeamID(m *ddl.ManuscriptTeamAssociation) ([]entity.ManuscriptTeamAssociation, error) {
	var res []entity.ManuscriptTeamAssociation

	query := s.db.Model(&ddl.ManuscriptTeamAssociation{}).Where(
		&ddl.ManuscriptTeamAssociation{
			TeamID: m.TeamID,
		},
	)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}
