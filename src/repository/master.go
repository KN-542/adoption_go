package repository

import (
	"api/src/model"

	"gorm.io/gorm"
)

type IMasterRepository interface {
	/*
		m_site
	*/
	// insert
	InsertSite(m *model.Site) error
	// select by primary key
	SelectByPrimaryKey(key int) (*model.Site, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) IMasterRepository {
	return &masterRepository{db}
}

/*
	m_site
*/
// insert
func (r *masterRepository) InsertSite(m *model.Site) error {
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	return nil
}

// select by primary key
func (r *masterRepository) SelectByPrimaryKey(key int) (*model.Site, error) {
	var res model.Site
	if err := r.db.First(&res, key).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
