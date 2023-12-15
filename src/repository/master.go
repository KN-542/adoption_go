package repository

import (
	"api/src/model"
	"log"

	"gorm.io/gorm"
)

type IMasterRepository interface {
	/*
		m_site
	*/
	// insert
	InsertSite(m *model.Site) error
	// select
	SelectSite() (*[]model.Site, error)
	// select by primary key
	SelectSiteByPrimaryKey(key int) (*model.Site, error)
	/*
		m_role
	*/
	// insert
	InsertRole(m *model.Role) error
	// select
	SelectRole() (*[]model.Role, error)
	/*
		m_applicant_status
	*/
	// insert
	InsertApplicantStatus(m *model.ApplicantStatus) error
	// select
	SelectApplicantStatus() (*[]model.ApplicantStatus, error)
}

type MasterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) IMasterRepository {
	return &MasterRepository{db}
}

/*
	m_site
*/
// insert
func (r *MasterRepository) InsertSite(m *model.Site) error {
	if err := r.db.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectSite() (*[]model.Site, error) {
	var res []model.Site
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

// select by primary key
func (r *MasterRepository) SelectSiteByPrimaryKey(key int) (*model.Site, error) {
	var res model.Site
	if err := r.db.First(&res, key).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

/*
	m_role
*/
// insert
func (r *MasterRepository) InsertRole(m *model.Role) error {
	if err := r.db.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectRole() (*[]model.Role, error) {
	var res []model.Role
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

/*
	m_applicant_status
*/
// insert
func (r *MasterRepository) InsertApplicantStatus(m *model.ApplicantStatus) error {
	if err := r.db.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectApplicantStatus() (*[]model.ApplicantStatus, error) {
	var res []model.ApplicantStatus
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}
