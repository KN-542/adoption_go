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
	SelectSiteByPrimaryKey(key int) (*model.Site, error)
	/*
		m_role
	*/
	// insert
	InsertRole(m *model.Role) error
	// select
	SelectRole() (*[]model.Role, error)
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
func (r *masterRepository) SelectSiteByPrimaryKey(key int) (*model.Site, error) {
	var res model.Site
	if err := r.db.First(&res, key).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

/*
	m_role
*/
// insert
func (r *masterRepository) InsertRole(m *model.Role) error {
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	return nil
}

// select
func (r *masterRepository) SelectRole() (*[]model.Role, error) {
	var res []model.Role
	if err := r.db.Find(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
