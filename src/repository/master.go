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
	// select by hash keys
	SelectSiteIDs(m []string) ([]uint, error)
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
		m_schedule_freq_status
	*/
	// insert
	InsertScheduleFreqStatus(tx *gorm.DB, m *ddl.ScheduleFreqStatus) error
	// select
	SelectScheduleFreqStatus() ([]entity.ScheduleFreqStatus, error)
	/*
		m_select_status_event
	*/
	// insert
	InsertSelectStatusEvent(tx *gorm.DB, m *ddl.SelectStatusEvent) error
	// select
	SelectSelectStatusEvent() ([]entity.SelectStatusEvent, error)
	// list
	ListSelectStatusEvent() ([]entity.SelectStatusEvent, error)
	// ハッシュキーから取得
	SelectSelectStatusEventByHashKeys(m []string) ([]entity.SelectStatusEvent, error)
	/*
		m_assign_rule
	*/
	// insert
	InsertAssignRule(tx *gorm.DB, m *ddl.AssignRule) error
	// select
	ListAssignRule() ([]entity.AssignRule, error)
	// select by hash key
	SelectAssignRule(m *ddl.AssignRule) (*entity.AssignRule, error)
	/*
		m_auto_assign_rule
	*/
	// insert
	InsertAutoAssignRule(tx *gorm.DB, m *ddl.AutoAssignRule) error
	// list
	ListAutoAssignRule() ([]entity.AutoAssignRule, error)
	// select by hash key
	SelectAutoAssignRule(m *ddl.AutoAssignRule) (*entity.AutoAssignRule, error)
	/*
		m_document_rule
	*/
	// insert
	InsertDocumentRule(tx *gorm.DB, m *ddl.DocumentRule) error
	// select by hash
	SelectDocumentRuleByHash(m *ddl.DocumentRule) (*entity.DocumentRule, error)
	// list
	ListDocumentRule() ([]entity.DocumentRule, error)
	/*
		m_occupation
	*/
	// insert
	InsertOccupation(tx *gorm.DB, m *ddl.Occupation) error
	// select by hash
	SelectOccupationByHash(m *ddl.Occupation) (*entity.Occupation, error)
	// list
	ListOccupation() ([]entity.Occupation, error)
	/*
		m_interview_processing
	*/
	// insert
	InsertProcessing(tx *gorm.DB, m *ddl.Processing) error
	// list
	ListProcessing() ([]entity.Processing, error)
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

// select by hash key
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

// select by hash keys
func (r *MasterRepository) SelectSiteIDs(m []string) ([]uint, error) {
	var res []entity.Site
	if err := r.db.Model(&ddl.Site{}).
		Select("id").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var IDs []uint
	for _, row := range res {
		IDs = append(IDs, row.ID)
	}

	return IDs, nil
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
		Joins("LEFT JOIN m_login_type ON m_sidebar.func_type = m_login_type.id").
		Joins("LEFT JOIN m_sidebar_role_association ON m_sidebar.id = m_sidebar_role_association.sidebar_id").
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
	m_schedule_freq_status
*/
// insert
func (r *MasterRepository) InsertScheduleFreqStatus(tx *gorm.DB, m *ddl.ScheduleFreqStatus) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectScheduleFreqStatus() ([]entity.ScheduleFreqStatus, error) {
	var res []entity.ScheduleFreqStatus
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_select_status_event
*/
// insert
func (r *MasterRepository) InsertSelectStatusEvent(tx *gorm.DB, m *ddl.SelectStatusEvent) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select
func (r *MasterRepository) SelectSelectStatusEvent() ([]entity.SelectStatusEvent, error) {
	var res []entity.SelectStatusEvent
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// list
func (r *MasterRepository) ListSelectStatusEvent() ([]entity.SelectStatusEvent, error) {
	var res []entity.SelectStatusEvent
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// ハッシュキーから取得
func (r *MasterRepository) SelectSelectStatusEventByHashKeys(m []string) ([]entity.SelectStatusEvent, error) {
	var res []entity.SelectStatusEvent
	if err := r.db.Where("hash_key IN ?", m).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_assign_rule
*/
// insert
func (r *MasterRepository) InsertAssignRule(tx *gorm.DB, m *ddl.AssignRule) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// list
func (r *MasterRepository) ListAssignRule() ([]entity.AssignRule, error) {
	var res []entity.AssignRule
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// select by hash key
func (r *MasterRepository) SelectAssignRule(m *ddl.AssignRule) (*entity.AssignRule, error) {
	var res entity.AssignRule
	if err := r.db.Model(&ddl.AssignRule{}).Where(&ddl.AssignRule{
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
	m_auto_assign_rule
*/
// insert
func (r *MasterRepository) InsertAutoAssignRule(tx *gorm.DB, m *ddl.AutoAssignRule) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// list
func (r *MasterRepository) ListAutoAssignRule() ([]entity.AutoAssignRule, error) {
	var res []entity.AutoAssignRule
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// select by hash key
func (r *MasterRepository) SelectAutoAssignRule(m *ddl.AutoAssignRule) (*entity.AutoAssignRule, error) {
	var res entity.AutoAssignRule
	if err := r.db.Model(&ddl.AutoAssignRule{}).Where(&ddl.AutoAssignRule{
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
	m_document_rule
*/
// insert
func (r *MasterRepository) InsertDocumentRule(tx *gorm.DB, m *ddl.DocumentRule) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select by hash
func (r *MasterRepository) SelectDocumentRuleByHash(m *ddl.DocumentRule) (*entity.DocumentRule, error) {
	var res entity.DocumentRule
	if err := r.db.Model(&ddl.DocumentRule{}).Where(&ddl.DocumentRule{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: m.HashKey,
		},
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

// list
func (r *MasterRepository) ListDocumentRule() ([]entity.DocumentRule, error) {
	var res []entity.DocumentRule
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_occupation
*/
// insert
func (r *MasterRepository) InsertOccupation(tx *gorm.DB, m *ddl.Occupation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// select by hash
func (r *MasterRepository) SelectOccupationByHash(m *ddl.Occupation) (*entity.Occupation, error) {
	var res entity.Occupation
	if err := r.db.Table("m_occupation").Where(&ddl.Occupation{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: m.HashKey,
		},
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

// list
func (r *MasterRepository) ListOccupation() ([]entity.Occupation, error) {
	var res []entity.Occupation
	if err := r.db.Table("m_occupation").Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

/*
	m_interview_processing
*/
// insert
func (r *MasterRepository) InsertProcessing(tx *gorm.DB, m *ddl.Processing) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// list
func (r *MasterRepository) ListProcessing() ([]entity.Processing, error) {
	var res []entity.Processing
	if err := r.db.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}
