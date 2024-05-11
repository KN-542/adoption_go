package repository

import (
	"api/src/model/ddl"
	"api/src/model/enum"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	// 登録
	Insert(tx *gorm.DB, applicant *ddl.Applicant) error
	// 検索
	Search(m *ddl.ApplicantSearchRequest) ([]ddl.ApplicantWith, error)
	// 取得(Email)
	GetByEmail(applicant *ddl.Applicant) ([]ddl.Applicant, error)
	// PK検索(カウント)
	CountByPrimaryKey(key *string) (*int64, error)
	// 応募者取得(ハッシュキー)
	GetByHashKey(m *ddl.Applicant) (*ddl.Applicant, error)
	// 面接希望日取得
	GetDesiredAt(m *ddl.Applicant) (*ddl.UserSchedule, error)
	// Google Meet Url 格納
	UpdateGoogleMeet(tx *gorm.DB, m *ddl.Applicant) error
	// 書類登録状況更新
	UpdateDocument(tx *gorm.DB, m *ddl.Applicant) error
	// 面接希望日更新
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

// 検索 TODO 検索仕様追加
func (a *ApplicantRepository) Search(m *ddl.ApplicantSearchRequest) ([]ddl.ApplicantWith, error) {
	var applicants []ddl.ApplicantWith

	query := a.db.Model(&ddl.Applicant{}).
		Select(`
			t_applicant.name,
			t_applicant.email,
			t_applicant.age,
			t_applicant.resume,
			t_applicant.curriculum_vitae,
			t_applicant.google_meet_url,
			t_applicant.calendar_id,
			m_applicant_status.status_name_ja as status_name_ja,
			m_site.site_name_ja as site_name_ja,
			t_user_schedule.start as start
		`).
		Joins("left join m_applicant_status on t_applicant.status = m_applicant_status.id").
		Joins("left join m_site on t_applicant.site_id = m_site.id").
		Joins("left join t_user_schedule on t_applicant.calendar_id = t_user_schedule.id")

	if len(m.SiteIDList) > 0 {
		query = query.Where("t_applicant.site_id IN ?", m.SiteIDList)
	}

	if len(m.ApplicantStatusList) > 0 {
		query = query.Where("t_applicant.status IN ?", m.ApplicantStatusList)
	}

	if m.Resume == uint(enum.DOCUMENT_EXIST) {
		query = query.Where("t_applicant.resume != ''")
	} else if m.Resume == uint(enum.DOCUMENT_NOT_EXIST) {
		query = query.Where("t_applicant.resume = ''")
	}
	if m.CurriculumVitae == uint(enum.DOCUMENT_EXIST) {
		query = query.Where("t_applicant.curriculum_vitae != ''")
	} else if m.CurriculumVitae == uint(enum.DOCUMENT_NOT_EXIST) {
		query = query.Where("t_applicant.curriculum_vitae = ''")
	}

	if m.Name != "" {
		query = query.Where("name LIKE ?", "%"+m.Name+"%")
	}
	if m.Email != "" {
		query = query.Where("email LIKE ?", "%"+m.Email+"%")
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
func (a *ApplicantRepository) GetByEmail(applicant *ddl.Applicant) ([]ddl.Applicant, error) {
	var l []ddl.Applicant
	if err := a.db.Where(applicant).Find(&l).Error; err != nil {
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
