package repository

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"log"

	"gorm.io/gorm"
)

type IMasterRepository interface {
	/*
		m_login_type
	*/
	// insert
	InsertLoginType(tx *gorm.DB, m *ddl.LoginType) error
	/*
		m_site
	*/
	// list
	InsertSite(tx *gorm.DB, m *ddl.Site) error
	// select
	ListSite() ([]entity.Site, error)
	// select by hash key
	SelectSite(m *ddl.Site) (*entity.Site, error)
	/*
		m_role
	*/
	// insert
	InsertRole(tx *gorm.DB, m *ddl.Role) error
	// select
	ListRole(m *ddl.Role) ([]entity.Role, error)
	/*
		m_sidebar
	*/
	// insert
	InsertSidebar(tx *gorm.DB, m *ddl.Sidebar) error
	// filter list
	ListSidebar(roles []ddl.Role, loginType *ddl.LoginType) ([]entity.Sidebar, error)
	/*
		m_sidebar_role_association
	*/
	// insert
	InsertSidebarRoleAssociation(tx *gorm.DB, m *ddl.SidebarRoleAssociation) error
	/*
		m_hash_key_pre
	*/
	// insert
	InsertHashKeyPre(tx *gorm.DB, m *ddl.HashKeyPre) error
	/*
		m_calendar_freq_status
	*/
	// insert
	InsertCalendarFreqStatus(tx *gorm.DB, m *ddl.CalendarFreqStatus) error
	// select
	SelectCalendarFreqStatus() ([]entity.CalendarFreqStatus, error)
}

type MasterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) IMasterRepository {
	return &MasterRepository{db}
}

/*
	m_login_type
*/
// insert
func (r *MasterRepository) InsertLoginType(tx *gorm.DB, m *ddl.LoginType) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

/*
	m_site
*/
// insert
func (r *MasterRepository) InsertSite(tx *gorm.DB, m *ddl.Site) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// list
func (r *MasterRepository) ListSite() ([]entity.Site, error) {
	var res []entity.Site
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// select by primary key
func (r *MasterRepository) SelectSite(m *ddl.Site) (*entity.Site, error) {
	var res entity.Site
	if err := r.db.Model(&ddl.Site{}).Where(&ddl.Site{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: m.HashKey,
		},
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

/*
	m_role
*/
// insert
func (r *MasterRepository) InsertRole(tx *gorm.DB, m *ddl.Role) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) ListRole(m *ddl.Role) ([]entity.Role, error) {
	var res []entity.Role
	if err := r.db.Where(m).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_sidebar
*/
// insert
func (r *MasterRepository) InsertSidebar(tx *gorm.DB, m *ddl.Sidebar) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// filter list
func (r *MasterRepository) ListSidebar(roles []ddl.Role, loginType *ddl.LoginType) ([]entity.Sidebar, error) {
	var res []entity.Sidebar

	var roleIDs []uint
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	query := r.db.Model(&entity.Sidebar{}).
		Select(`
			m_sidebar.name_ja,
			m_sidebar.name_en,
			m_sidebar.path
		`).
		Joins("left join m_login_type on m_sidebar.func_type = m_login_type.id").
		Joins("left join m_sidebar_role_association on m_sidebar.id = m_sidebar_role_association.sidebar_id").
		Where("m_sidebar.func_type = ?", loginType.ID)

	if len(roles) > 0 {
		query = query.Where("m_sidebar_role_association.role_id IN ?", roleIDs)
	}

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_sidebar_role_association
*/
// insert
func (r *MasterRepository) InsertSidebarRoleAssociation(tx *gorm.DB, m *ddl.SidebarRoleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

/*
	m_hash_key_pre
*/
// insert
func (r *MasterRepository) InsertHashKeyPre(tx *gorm.DB, m *ddl.HashKeyPre) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

/*
	m_calendar_freq_status
*/
// insert
func (r *MasterRepository) InsertCalendarFreqStatus(tx *gorm.DB, m *ddl.CalendarFreqStatus) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectCalendarFreqStatus() ([]entity.CalendarFreqStatus, error) {
	var res []entity.CalendarFreqStatus
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}
