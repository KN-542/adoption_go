package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"api/src/model/static"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *ddl.Applicant) error
	// 一括登録
	Inserts(tx *gorm.DB, m []*ddl.Applicant) error
	// 更新
	Update(tx *gorm.DB, m *ddl.Applicant) error
	// 検索
	Search(m *dto.SearchApplicant) ([]*entity.SearchApplicant, error)
	// 取得
	Get(m *ddl.Applicant) (*entity.Applicant, error)
	// 応募者ステータス一覧
	ListStatus(m *ddl.SelectStatus) ([]entity.ApplicantStatus, error)
	// 応募者ステータス削除
	DeleteStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// 応募者ステータス削除_PK
	DeleteStatusByPrimary(tx *gorm.DB, m *ddl.SelectStatus, ids []uint64) error
	// 取得_メールアドレス
	GetByEmail(m *ddl.Applicant) ([]entity.Applicant, error)
	// 取得_チームID
	GetByTeamID(m *ddl.Applicant) ([]entity.Applicant, error)
	// 応募者重複チェック_媒体側ID
	CheckDuplByOuterId(m *dto.CheckDuplDownloading) ([]entity.Applicant, error)
	// 選考ステータス更新
	UpdateSelectStatus(tx *gorm.DB, m *ddl.Applicant) error
	// 応募者面接予定紐づけ更新
	UpdateApplicantScheduleAssociation(tx *gorm.DB, m *ddl.ApplicantScheduleAssociation) error
}

type ApplicantRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewApplicantRepository(db *gorm.DB, redis *redis.Client) IApplicantRepository {
	return &ApplicantRepository{db, redis}
}

// 登録
func (a *ApplicantRepository) Insert(tx *gorm.DB, m *ddl.Applicant) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 一括登録
func (a *ApplicantRepository) Inserts(tx *gorm.DB, m []*ddl.Applicant) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 更新
func (a *ApplicantRepository) Update(tx *gorm.DB, m *ddl.Applicant) error {
	if err := tx.Where(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: m.HashKey,
		},
	}).Updates(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
		Status:          m.Status,
		Resume:          m.Resume,
		CurriculumVitae: m.CurriculumVitae,
		GoogleMeetURL:   m.GoogleMeetURL,
		TeamID:          m.TeamID,
	}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索
func (a *ApplicantRepository) Search(m *dto.SearchApplicant) ([]*entity.SearchApplicant, error) {
	var applicants []*entity.SearchApplicant

	query := a.db.Table("t_applicant").
		Select(`
			t_applicant.id,
			t_applicant.hash_key,
			t_applicant.name,
			t_applicant.email,
			t_applicant.resume,
			t_applicant.curriculum_vitae,
			t_applicant.google_meet_url,
			t_applicant.created_at,
			t_select_status.status_name as status_name,
			m_site.site_name as site_name,
			t_schedule.hash_key as schedule_hash_key,
			t_schedule.start as start
		`).
		Joins("left join t_select_status on t_applicant.status = t_select_status.id").
		Joins("left join m_site on t_applicant.site_id = m_site.id").
		Joins("left join t_applicant_user_association on t_applicant_user_association.applicant_id = t_applicant.id").
		Joins("left join t_user on t_applicant_user_association.user_id = t_user.id").
		Joins("left join t_applicant_schedule_association on t_applicant_schedule_association.applicant_id = t_applicant.id").
		Joins("left join t_schedule on t_applicant_schedule_association.schedule_id = t_schedule.id").
		Where("t_applicant.team_id = ?", m.TeamID).
		Where("t_applicant.company_id = ?", m.CompanyID)

	if len(m.Sites) > 0 {
		query = query.Where("m_site.hash_key IN ?", m.Sites)
	}

	if len(m.ApplicantStatusList) > 0 {
		query = query.Where("t_select_status.hash_key IN ?", m.ApplicantStatusList)
	}

	if m.ResumeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where("t_applicant.resume != ''")
	} else {
		query = query.Where("t_applicant.resume = ''")
	}
	if m.CurriculumVitaeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where("t_applicant.curriculum_vitae != ''")
	} else {
		query = query.Where("t_applicant.curriculum_vitae = ''")
	}

	if m.Name != "" {
		query = query.Where("t_applicant.name LIKE ?", "%"+m.Name+"%")
	}
	if m.Email != "" {
		query = query.Where("t_applicant.email LIKE ?", "%"+m.Email+"%")
	}

	if len(m.Users) > 0 {
		query = query.Where("t_user.hash_key IN ?", m.Users)
	}

	if m.InterviewerDateFrom.Year() >= 1900 && m.InterviewerDateTo.Year() >= 1900 {
		query = query.Where("t_schedule.start >= ? AND t_schedule.start < ?", m.InterviewerDateFrom, m.InterviewerDateTo.AddDate(0, 0, 1))
	} else if m.InterviewerDateFrom.Year() >= 1900 {
		query = query.Where("t_schedule.start >= ?", m.InterviewerDateFrom)
	} else if m.InterviewerDateTo.Year() >= 1900 {
		query = query.Where("t_schedule.start < ?", m.InterviewerDateTo.AddDate(0, 0, 1))
	}

	if m.CreatedAtFrom.Year() >= 1900 && m.CreatedAtTo.Year() >= 1900 {
		query = query.Where("t_applicant.created_at >= ? AND t_applicant.created_at < ?", m.CreatedAtFrom, m.CreatedAtTo.AddDate(0, 0, 1))
	} else if m.CreatedAtFrom.Year() >= 1900 {
		query = query.Where("t_applicant.created_at >= ?", m.CreatedAtFrom)
	} else if m.CreatedAtTo.Year() >= 1900 {
		query = query.Where("t_applicant.created_at < ?", m.CreatedAtTo.AddDate(0, 0, 1))
	}

	if m.SortKey != "" {
		if m.SortAsc {
			query = query.Order(m.SortKey + " ASC")
		} else {
			query = query.Order(m.SortKey + " DESC")
		}
	}

	if err := query.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "hash_key", "name")
	}).Find(&applicants).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return applicants, nil
}

// 応募者取得(ハッシュキー)
func (a *ApplicantRepository) Get(m *ddl.Applicant) (*entity.Applicant, error) {
	var res entity.Applicant
	if err := a.db.Model(&ddl.Applicant{}).
		Select("t_applicant.*, t_applicant_schedule_association.schedule_id").
		Joins("left join t_applicant_schedule_association on t_applicant_schedule_association.applicant_id = t_applicant.id").
		Where(
			&ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: m.HashKey,
				},
			},
		).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 応募者ステータス一覧
func (a *ApplicantRepository) ListStatus(m *ddl.SelectStatus) ([]entity.ApplicantStatus, error) {
	var res []entity.ApplicantStatus

	if err := a.db.Where(
		&ddl.SelectStatus{
			TeamID: m.TeamID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 応募者ステータス削除
func (a *ApplicantRepository) DeleteStatus(tx *gorm.DB, m *ddl.SelectStatus) error {
	if err := tx.Where(&ddl.SelectStatus{
		TeamID: m.TeamID,
	}).Delete(&ddl.SelectStatus{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 応募者ステータス削除_PK
func (a *ApplicantRepository) DeleteStatusByPrimary(tx *gorm.DB, m *ddl.SelectStatus, ids []uint64) error {
	if err := tx.Where(&ddl.SelectStatus{
		TeamID: m.TeamID,
	}).Where("id IN ?", ids).Delete(&ddl.SelectStatus{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// PK検索(カウント)
func (a *ApplicantRepository) CountByPrimaryKey(key *string) (*int64, error) {
	var count int64
	if err := a.db.Model(&ddl.Applicant{}).Where("id = ?", key).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &count, nil
}

// 取得_メールアドレス
func (a *ApplicantRepository) GetByEmail(m *ddl.Applicant) ([]entity.Applicant, error) {
	var l []entity.Applicant
	if err := a.db.Where(&ddl.Applicant{
		Email: m.Email,
	}).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// 取得_チームID
func (a *ApplicantRepository) GetByTeamID(m *ddl.Applicant) ([]entity.Applicant, error) {
	var l []entity.Applicant
	if err := a.db.Where(&ddl.Applicant{
		TeamID: m.TeamID,
	}).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// 応募者重複チェック_媒体側ID
func (a *ApplicantRepository) CheckDuplByOuterId(m *dto.CheckDuplDownloading) ([]entity.Applicant, error) {
	var res []entity.Applicant
	if err := a.db.Table("t_applicant").
		Where(&ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: m.CompanyID,
			},
			TeamID: m.TeamID,
		}).
		Where("outer_id IN ?", m.List).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接希望日取得
func (a *ApplicantRepository) GetDesiredAt(m *ddl.Applicant) (*ddl.Schedule, error) {
	var l ddl.Schedule
	if err := a.db.Model(&ddl.Schedule{}).
		Select("t_schedule.start, t_schedule.end").
		Joins("left join t_applicant on t_applicant.schedule_id = t_schedule.id").
		Where(
			&ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: m.HashKey,
				},
			},
		).First(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &l, nil
}

// 選考ステータス更新
func (a *ApplicantRepository) UpdateSelectStatus(tx *gorm.DB, m *ddl.Applicant) error {
	sql := `
		UPDATE t_applicant
		SET status = ?, updated_at = ?
		FROM t_select_status
		WHERE t_applicant.status = t_select_status.id AND t_select_status.hash_key = ?
	`

	if err := tx.Exec(sql, m.Status, time.Now(), m.HashKey).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 応募者面接予定紐づけ更新
func (a *ApplicantRepository) UpdateApplicantScheduleAssociation(tx *gorm.DB, m *ddl.ApplicantScheduleAssociation) error {
	if err := tx.Model(&ddl.ApplicantScheduleAssociation{}).
		Where(&ddl.ApplicantScheduleAssociation{
			ApplicantID: m.ApplicantID,
		}).
		Updates(&ddl.ApplicantScheduleAssociation{
			ScheduleID: m.ScheduleID,
		}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}
