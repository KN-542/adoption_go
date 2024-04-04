package repository

import (
	model "api/src/model"
	enum "api/src/model/enum"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	// 登録
	Insert(tx *gorm.DB, applicant *model.Applicant) error
	// 検索
	Search(m *model.ApplicantSearchRequest) ([]model.ApplicantWith, error)
	// 取得(Email)
	GetByEmail(applicant *model.Applicant) ([]model.Applicant, error)
	// PK検索(カウント)
	CountByPrimaryKey(key *string) (*int64, error)
	// 応募者取得(ハッシュキー)
	GetByHashKey(m *model.Applicant) (*model.Applicant, error)
	// Google Meet Url 格納
	UpdateGoogleMeet(tx *gorm.DB, m *model.Applicant) error
	// 書類登録状況更新
	UpdateDocument(tx *gorm.DB, m *model.Applicant) error
	// 面接希望日更新
	UpdateDesiredAt(tx *gorm.DB, m *model.Applicant) error
}

type ApplicantRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// 登録
func NewApplicantRepository(db *gorm.DB, redis *redis.Client) IApplicantRepository {
	return &ApplicantRepository{db, redis}
}

func (a *ApplicantRepository) Insert(tx *gorm.DB, applicant *model.Applicant) error {
	if err := tx.Create(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索 TODO 検索仕様追加
func (a *ApplicantRepository) Search(m *model.ApplicantSearchRequest) ([]model.ApplicantWith, error) {
	var applicants []model.ApplicantWith

	query := a.db.Model(&model.Applicant{}).
		Select("t_applicant.*, m_applicant_status.status_name_ja as status_name_ja, m_site.site_name_ja as site_name_ja").
		Joins("left join m_applicant_status on t_applicant.status = m_applicant_status.id").
		Joins("left join m_site on t_applicant.site_id = m_site.id")

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
	if m.Users != "" {
		query = query.Where("users LIKE ?", "%"+m.Users+"%")
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
	if err := a.db.Model(&model.Applicant{}).Where("id = ?", key).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &count, nil
}

// 取得(Email)
func (a *ApplicantRepository) GetByEmail(applicant *model.Applicant) ([]model.Applicant, error) {
	var l []model.Applicant
	if err := a.db.Where(applicant).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// 応募者取得(ハッシュキー)
func (a *ApplicantRepository) GetByHashKey(m *model.Applicant) (*model.Applicant, error) {
	var res model.Applicant
	if err := a.db.Where(
		&model.Applicant{
			HashKey: m.HashKey,
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// Google Meet Url 格納
func (a *ApplicantRepository) UpdateGoogleMeet(tx *gorm.DB, m *model.Applicant) error {
	applicant := model.Applicant{
		GoogleMeetURL: m.GoogleMeetURL,
		UpdatedAt:     time.Now(),
	}
	if err := tx.Model(&model.Applicant{}).Where(
		&model.Applicant{
			HashKey: m.HashKey,
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 書類登録状況更新
func (a *ApplicantRepository) UpdateDocument(tx *gorm.DB, m *model.Applicant) error {
	applicant := model.Applicant{
		Resume:          m.Resume,
		CurriculumVitae: m.CurriculumVitae,
		UpdatedAt:       time.Now(),
	}
	if err := tx.Model(&model.Applicant{}).Where(
		&model.Applicant{
			HashKey: m.HashKey,
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 面接希望日更新
func (a *ApplicantRepository) UpdateDesiredAt(tx *gorm.DB, m *model.Applicant) error {
	applicant := model.Applicant{
		DesiredAt:       m.DesiredAt,
		CalendarHashKey: m.CalendarHashKey,
		Users:           m.Users,
		UpdatedAt:       time.Now(),
	}
	if err := tx.Model(&model.Applicant{}).Where(
		&model.Applicant{
			HashKey: m.HashKey,
		},
	).Updates(applicant).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if applicant.Users == "" {
		if err := tx.Model(&model.Applicant{}).Where(
			&model.Applicant{
				HashKey: m.HashKey,
			},
		).Update("users", "").Error; err != nil {
			log.Printf("%v", err)
			return err
		}
	}

	return nil
}
