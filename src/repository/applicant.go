package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"api/src/model/static"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	// 登録
	Insert(tx *gorm.DB, applicant *ddl.Applicant) error
	// 検索
	Search(m *dto.ApplicantSearch) ([]*entity.ApplicantSearch, error)
	// 応募者ステータス一覧
	ListStatus(m *ddl.SelectStatus) ([]entity.ApplicantStatus, error)
	// 取得(Email)
	GetByEmail(applicant *ddl.Applicant) ([]entity.Applicant, error)
	// 応募者重複チェック(outer_id)
	CheckDuplByOuterId(m *ddl.Applicant) (*int64, error)
	// PK検索(カウント)*
	CountByPrimaryKey(key *string) (*int64, error)
	// 応募者取得(ハッシュキー)*
	GetByHashKey(m *ddl.Applicant) (*ddl.Applicant, error)
	// 面接希望日取得*
	GetDesiredAt(m *ddl.Applicant) (*ddl.UserSchedule, error)
	// Google Meet Url 格納*
	UpdateGoogleMeet(tx *gorm.DB, m *ddl.Applicant) error
	// 書類登録状況更新*
	UpdateDocument(tx *gorm.DB, m *ddl.Applicant) error
	// 面接希望日更新*
	UpdateDesiredAt(tx *gorm.DB, m *ddl.Applicant) error
}

type ApplicantRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewApplicantRepository(db *gorm.DB, redis *redis.Client) IApplicantRepository {
	return &ApplicantRepository{db, redis}
}

// 登録
func (a *ApplicantRepository) Insert(tx *gorm.DB, applicant *ddl.Applicant) error {
	if err := tx.Create(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索
func (a *ApplicantRepository) Search(m *dto.ApplicantSearch) ([]*entity.ApplicantSearch, error) {
	var applicants []*entity.ApplicantSearch

	year := time.Now().Year()
	month := time.Now().Month()
	tApplicant := fmt.Sprintf("t_applicant_%d_%02d", year, month)
	tApplicantUserAssociation := fmt.Sprintf("t_applicant_user_association_%d_%02d", year, month)

	query := a.db.Table(tApplicant).
		Select(fmt.Sprintf(`
				%s.id,
				%s.hash_key,
				%s.name,
				%s.email,
				%s.resume,
				%s.curriculum_vitae,
				%s.google_meet_url,
				%s.calendar_id,
				t_select_status.status_name as status_name,
				m_site.site_name as site_name,
				t_user_schedule.hash_key as calendar_hash_key,
				t_user_schedule.start as start
			`,
			tApplicant,
			tApplicant,
			tApplicant,
			tApplicant,
			tApplicant,
			tApplicant,
			tApplicant,
			tApplicant,
		)).
		Joins(fmt.Sprintf("left join t_select_status on %s.status = t_select_status.id", tApplicant)).
		Joins(fmt.Sprintf("left join %s on %s.id = %s.applicant_id", tApplicantUserAssociation, tApplicant, tApplicantUserAssociation)).
		Joins(fmt.Sprintf("left join m_site on %s.site_id = m_site.id", tApplicant)).
		Joins(fmt.Sprintf("left join t_user_schedule on %s.calendar_id = t_user_schedule.id", tApplicant)).
		Where(fmt.Sprintf("%s.team_id = ?", tApplicant), m.TeamID)

	if len(m.Sites) > 0 {
		query = query.Where("m_site.hash_key IN ?", m.Sites)
	}

	if len(m.ApplicantStatusList) > 0 {
		query = query.Where("t_select_status.hash_key IN ?", m.ApplicantStatusList)
	}

	if m.ResumeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where(fmt.Sprintf("%s.resume != ''", tApplicant))
	} else if m.ResumeFlg == uint(static.DOCUMENT_NOT_EXIST) {
		query = query.Where(fmt.Sprintf("%s.resume = ''", tApplicant))
	}
	if m.CurriculumVitaeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where(fmt.Sprintf("%s.curriculum_vitae != ''", tApplicant))
	} else if m.CurriculumVitaeFlg == uint(static.DOCUMENT_NOT_EXIST) {
		query = query.Where(fmt.Sprintf("%s.curriculum_vitae = ''", tApplicant))
	}

	if m.Name != "" {
		query = query.Where(fmt.Sprintf("%s.name LIKE ?", tApplicant), "%"+m.Name+"%")
	}
	if m.Email != "" {
		query = query.Where(fmt.Sprintf("%s.email LIKE ?", tApplicant), "%"+m.Email+"%")
	}

	if len(m.UserIDs) > 0 {
		query = query.Where(fmt.Sprintf("%s.user_id IN ?", tApplicantUserAssociation), m.UserIDs)
	}

	if m.SortKey != "" {
		if m.SortAsc {
			query = query.Order(m.SortKey + " ASC")
		} else {
			query = query.Order(m.SortKey + " DESC")
		}
	}

	if err := query.Find(&applicants).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return applicants, nil
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

// PK検索(カウント)
func (a *ApplicantRepository) CountByPrimaryKey(key *string) (*int64, error) {
	var count int64
	if err := a.db.Model(&ddl.Applicant{}).Where("id = ?", key).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &count, nil
}

// 取得(Email)
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

// 応募者取得(ハッシュキー)
func (a *ApplicantRepository) GetByHashKey(m *ddl.Applicant) (*ddl.Applicant, error) {
	var res ddl.Applicant
	if err := a.db.Where(
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

// 応募者重複チェック(outer_id)
func (a *ApplicantRepository) CheckDuplByOuterId(m *ddl.Applicant) (*int64, error) {
	var count int64
	if err := a.db.Model(&ddl.Applicant{}).Where("outer_id = ?", m.OuterID).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &count, nil
}

// 面接希望日取得
func (a *ApplicantRepository) GetDesiredAt(m *ddl.Applicant) (*ddl.UserSchedule, error) {
	var l ddl.UserSchedule
	if err := a.db.Model(&ddl.UserSchedule{}).
		Select("t_user_schedule.start, t_user_schedule.end").
		Joins("left join t_applicant on t_applicant.calendar_id = t_user_schedule.id").
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

// Google Meet Url 格納
func (a *ApplicantRepository) UpdateGoogleMeet(tx *gorm.DB, m *ddl.Applicant) error {
	applicant := ddl.Applicant{
		GoogleMeetURL: m.GoogleMeetURL,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.Applicant{}).Where(
		&ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 書類登録状況更新
func (a *ApplicantRepository) UpdateDocument(tx *gorm.DB, m *ddl.Applicant) error {
	applicant := ddl.Applicant{
		Resume:          m.Resume,
		CurriculumVitae: m.CurriculumVitae,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.Applicant{}).Where(
		&ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 面接希望日更新
func (a *ApplicantRepository) UpdateDesiredAt(tx *gorm.DB, m *ddl.Applicant) error {
	applicant := ddl.Applicant{
		CalendarID: m.CalendarID,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.Applicant{}).Where(
		&ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}
