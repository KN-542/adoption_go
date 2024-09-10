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
	Search(m *dto.SearchApplicant) ([]*entity.SearchApplicant, int64, error)
	// 取得
	Get(m *ddl.Applicant) (*entity.Applicant, error)
	// 種別登録
	InsertType(tx *gorm.DB, m *ddl.ApplicantType) error
	// 種別一覧
	ListType(m *ddl.ApplicantType) ([]entity.ApplicantType, error)
	// 応募者ステータス一覧
	ListStatus(m *ddl.SelectStatus) ([]entity.ApplicantStatus, error)
	// 応募者ステータス削除
	DeleteStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// 応募者ステータス削除_PK
	DeleteStatusByPrimary(tx *gorm.DB, m *ddl.SelectStatus, ids []uint64) error
	// 取得_メールアドレス
	GetByEmail(m *ddl.Applicant) (*entity.Applicant, error)
	// 取得_チームID
	GetByTeamID(m *ddl.Applicant) ([]entity.Applicant, error)
	// 応募者重複チェック_媒体側ID
	CheckDuplByOuterId(m *dto.CheckDuplDownloading) ([]entity.Applicant, error)
	// 選考ステータス更新
	UpdateSelectStatus(tx *gorm.DB, m *ddl.Applicant) error
	// 応募者面接予定紐づけ登録
	InsertApplicantScheduleAssociation(tx *gorm.DB, m *ddl.ApplicantScheduleAssociation) error
	// 応募者面接予定紐づけ更新
	UpdateApplicantScheduleAssociation(tx *gorm.DB, m *ddl.ApplicantScheduleAssociation) error
	// 履歴書登録
	InsertApplicantResumeAssociation(tx *gorm.DB, m *ddl.ApplicantResumeAssociation) error
	// 履歴書削除
	DeleteApplicantResumeAssociation(tx *gorm.DB, m *ddl.ApplicantResumeAssociation) error
	// 職務経歴書登録
	InsertApplicantCurriculumVitaeAssociation(tx *gorm.DB, m *ddl.ApplicantCurriculumVitaeAssociation) error
	// 職務経歴書削除
	DeleteApplicantCurriculumVitaeAssociation(tx *gorm.DB, m *ddl.ApplicantCurriculumVitaeAssociation) error
	// Google Meet URL登録
	InsertApplicantURLAssociation(tx *gorm.DB, m *ddl.ApplicantURLAssociation) error
	// Google Meet URL取得
	GetApplicantURLAssociation(m *ddl.Applicant) ([]entity.ApplicantURLAssociation, error)
	// ユーザー紐づけ一括登録
	InsertsUserAssociation(tx *gorm.DB, m []*ddl.ApplicantUserAssociation) error
	// ユーザー紐づけ取得_ユーザー
	GetUserAssociation(m *ddl.ApplicantUserAssociation) ([]entity.ApplicantUserAssociation, error)
	// ユーザー紐づけ削除
	DeleteUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error
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
		Status:         m.Status,
		TeamID:         m.TeamID,
		NumOfInterview: m.NumOfInterview,
	}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索
func (a *ApplicantRepository) Search(m *dto.SearchApplicant) ([]*entity.SearchApplicant, int64, error) {
	var applicants []*entity.SearchApplicant
	var totalCount int64

	query := a.db.Table("t_applicant").
		Joins(`
			LEFT JOIN
				t_select_status
			ON
				t_applicant.status = t_select_status.id
		`).
		Joins("LEFT JOIN m_site ON t_applicant.site_id = m_site.id").
		Joins(`
			LEFT JOIN
				t_applicant_schedule_association
			ON
				t_applicant_schedule_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_schedule
			ON
				t_applicant_schedule_association.schedule_id = t_schedule.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_resume_association
			ON
				t_applicant_resume_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_curriculum_vitae_association
			ON
				t_applicant_curriculum_vitae_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_url_association
			ON
				t_applicant_url_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_manuscript_applicant_association
			ON
				t_manuscript_applicant_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_manuscript
			ON
				t_manuscript_applicant_association.manuscript_id = t_manuscript.id
		`).
		Where("t_applicant.team_id = ? AND t_applicant.company_id = ?", m.TeamID, m.CompanyID)

	if len(m.Users) > 0 {
		query = query.Joins("INNER JOIN t_applicant_user_association ON t_applicant_user_association.applicant_id = t_applicant.id").
			Joins("INNER JOIN t_user ON t_applicant_user_association.user_id = t_user.id").
			Where("t_user.hash_key IN ?", m.Users)
	}

	if len(m.Sites) > 0 {
		query = query.Where("m_site.hash_key IN ?", m.Sites)
	}

	if len(m.ApplicantStatusList) > 0 {
		query = query.Where("t_select_status.hash_key IN ?", m.ApplicantStatusList)
	}

	if len(m.Manuscripts) > 0 {
		query = query.Where("t_manuscript.hash_key IN ?", m.Manuscripts)
	}

	if m.ResumeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where("t_applicant_resume_association.applicant_id IS NOT NULL")
	} else if m.ResumeFlg == uint(static.DOCUMENT_NOT_EXIST) {
		query = query.Where("t_applicant_resume_association.applicant_id IS NULL")
	}

	if m.CurriculumVitaeFlg == uint(static.DOCUMENT_EXIST) {
		query = query.Where("t_applicant_curriculum_vitae_association.applicant_id IS NOT NULL")
	} else if m.CurriculumVitaeFlg == uint(static.DOCUMENT_NOT_EXIST) {
		query = query.Where("t_applicant_curriculum_vitae_association.applicant_id IS NULL")
	}

	if m.Name != "" {
		query = query.Where("t_applicant.name LIKE ?", "%"+m.Name+"%")
	}
	if m.Email != "" {
		query = query.Where("t_applicant.email LIKE ?", "%"+m.Email+"%")
	}
	if m.CommitID != "" {
		query = query.Where("t_applicant.commit_id LIKE ?", "%"+m.CommitID+"%")
	}

	if !m.InterviewerDateFrom.IsZero() && !m.InterviewerDateTo.IsZero() {
		query = query.Where("t_schedule.start >= ? AND t_schedule.start < ?", m.InterviewerDateFrom, m.InterviewerDateTo.AddDate(0, 0, 1))
	} else if !m.InterviewerDateFrom.IsZero() {
		query = query.Where("t_schedule.start >= ?", m.InterviewerDateFrom)
	} else if !m.InterviewerDateTo.IsZero() {
		query = query.Where("t_schedule.start < ?", m.InterviewerDateTo.AddDate(0, 0, 1))
	}

	if !m.CreatedAtFrom.IsZero() && !m.CreatedAtTo.IsZero() {
		query = query.Where("t_applicant.created_at >= ? AND t_applicant.created_at < ?", m.CreatedAtFrom, m.CreatedAtTo.AddDate(0, 0, 1))
	} else if !m.CreatedAtFrom.IsZero() {
		query = query.Where("t_applicant.created_at >= ?", m.CreatedAtFrom)
	} else if !m.CreatedAtTo.IsZero() {
		query = query.Where("t_applicant.created_at < ?", m.CreatedAtTo.AddDate(0, 0, 1))
	}

	if m.SortKey != "" {
		if m.SortAsc {
			query = query.Order(m.SortKey + " ASC")
		} else {
			query = query.Order(m.SortKey + " DESC")
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("%v", err)
		return nil, 0, err
	}

	offset := (m.Page - 1) * m.PageSize
	err := query.Select(`
		t_applicant.hash_key,
		t_applicant.outer_id,
		t_applicant.site_id,
		t_applicant.status,
		t_applicant.name,
		t_applicant.email,
		t_applicant.age,
		t_applicant.commit_id,
		t_applicant.created_at,
		t_select_status.status_name,
		m_site.site_name,
		t_schedule.hash_key as schedule_hash_key,
		t_schedule.start,
		t_applicant_resume_association.extension as resume_extension,
		t_applicant_curriculum_vitae_association.extension as curriculum_vitae_extension,
		t_applicant_url_association.url as google_meet_url,
		t_manuscript.content as content
	`).
		Offset(offset).
		Limit(m.PageSize).
		Find(&applicants).
		Error

	if err != nil {
		log.Printf("%v", err)
		return nil, 0, err
	}

	if len(applicants) > 0 {
		var applicantIDs []uint64
		for _, app := range applicants {
			applicantIDs = append(applicantIDs, app.ID)
		}

		var userAssociations []struct {
			ApplicantID uint64
			UserID      uint64
			HashKey     string
			Name        string
		}

		err := a.db.Table("t_applicant_user_association").
			Select("t_applicant_user_association.applicant_id, t_user.id as user_id, t_user.hash_key, t_user.name").
			Joins("INNER JOIN t_user ON t_applicant_user_association.user_id = t_user.id").
			Where("t_applicant_user_association.applicant_id IN ?", applicantIDs).
			Find(&userAssociations).Error

		if err != nil {
			log.Printf("%v", err)
			return nil, 0, err
		}

		userMap := make(map[uint64][]*ddl.User)
		for _, assoc := range userAssociations {
			user := &ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					ID:      assoc.UserID,
					HashKey: assoc.HashKey,
				},
				Name: assoc.Name,
			}
			userMap[assoc.ApplicantID] = append(userMap[assoc.ApplicantID], user)
		}

		for _, app := range applicants {
			app.Users = userMap[app.ID]
		}
	}

	return applicants, totalCount, nil
}

// 応募者取得(ハッシュキー)
func (a *ApplicantRepository) Get(m *ddl.Applicant) (*entity.Applicant, error) {
	var res entity.Applicant
	if err := a.db.Model(&ddl.Applicant{}).
		Select(`
			t_applicant.*,
			t_applicant_schedule_association.schedule_id,
			t_applicant_resume_association.extension as resume_extension,
			t_applicant_curriculum_vitae_association.extension as curriculum_vitae_extension,
			t_applicant_url_association.url as google_meet_url
		`).
		Joins(`
			LEFT JOIN
				t_applicant_schedule_association
			ON
				t_applicant_schedule_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_resume_association
			ON
				t_applicant_resume_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_curriculum_vitae_association
			ON
				t_applicant_curriculum_vitae_association.applicant_id = t_applicant.id
		`).
		Joins(`
			LEFT JOIN
				t_applicant_url_association
			ON
				t_applicant_url_association.applicant_id = t_applicant.id
		`).
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

// 種別登録
func (a *ApplicantRepository) InsertType(tx *gorm.DB, m *ddl.ApplicantType) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 種別一覧
func (a *ApplicantRepository) ListType(m *ddl.ApplicantType) ([]entity.ApplicantType, error) {
	var res []entity.ApplicantType

	if err := a.db.Table("t_applicant_type").Select(`
			t_applicant_type.hash_key,
			t_applicant_type.name,
			m_document_rule.rule_ja as rule_ja,
			m_document_rule.rule_en as rule_en,
			m_occupation.name_ja as name_ja,
			m_occupation.name_en as name_en
		`).
		Joins(`
			LEFT JOIN
				m_document_rule
			ON
				m_document_rule.id = t_applicant_type.rule_id
		`).
		Joins(`
			LEFT JOIN
				m_occupation
			ON
				m_occupation.id = t_applicant_type.occupation_id
		`).
		Where(
			&ddl.ApplicantType{
				TeamID: m.TeamID,
			},
		).Limit(static.APPLICANT_TYPE_SIZE).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
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
func (a *ApplicantRepository) GetByEmail(m *ddl.Applicant) (*entity.Applicant, error) {
	var l entity.Applicant
	if err := a.db.Where(&ddl.Applicant{
		Email:  m.Email,
		TeamID: m.TeamID,
	}).First(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &l, nil
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
		Joins("LEFT JOIN t_applicant ON t_applicant.schedule_id = t_schedule.id").
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

// 応募者面接予定紐づけ登録
func (a *ApplicantRepository) InsertApplicantScheduleAssociation(tx *gorm.DB, m *ddl.ApplicantScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
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

// 履歴書登録
func (a *ApplicantRepository) InsertApplicantResumeAssociation(tx *gorm.DB, m *ddl.ApplicantResumeAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 履歴書削除
func (a *ApplicantRepository) DeleteApplicantResumeAssociation(tx *gorm.DB, m *ddl.ApplicantResumeAssociation) error {
	if err := tx.Where(&ddl.ApplicantResumeAssociation{
		ApplicantID: m.ApplicantID,
	}).Delete(&ddl.ApplicantResumeAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 職務経歴書登録
func (a *ApplicantRepository) InsertApplicantCurriculumVitaeAssociation(tx *gorm.DB, m *ddl.ApplicantCurriculumVitaeAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 職務経歴書削除
func (a *ApplicantRepository) DeleteApplicantCurriculumVitaeAssociation(tx *gorm.DB, m *ddl.ApplicantCurriculumVitaeAssociation) error {
	if err := tx.Where(&ddl.ApplicantCurriculumVitaeAssociation{
		ApplicantID: m.ApplicantID,
	}).Delete(&ddl.ApplicantCurriculumVitaeAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// Google Meet URL登録
func (a *ApplicantRepository) InsertApplicantURLAssociation(tx *gorm.DB, m *ddl.ApplicantURLAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// Google Meet URL取得
func (a *ApplicantRepository) GetApplicantURLAssociation(m *ddl.Applicant) ([]entity.ApplicantURLAssociation, error) {
	var l []entity.ApplicantURLAssociation
	if err := a.db.Model(&ddl.ApplicantURLAssociation{}).
		Joins("LEFT JOIN t_applicant ON t_applicant.id = t_applicant_url_association.applicant_id").
		Where("t_applicant.hash_key = ?", m.HashKey).
		Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// ユーザー紐づけ一括登録
func (u *ApplicantRepository) InsertsUserAssociation(tx *gorm.DB, m []*ddl.ApplicantUserAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー紐づけ取得
func (u *ApplicantRepository) GetUserAssociation(m *ddl.ApplicantUserAssociation) ([]entity.ApplicantUserAssociation, error) {
	var res []entity.ApplicantUserAssociation

	query := u.db.Table("t_applicant_user_association").
		Where(&ddl.ApplicantUserAssociation{
			ApplicantID: m.ApplicantID,
		})

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// ユーザー紐づけ削除
func (u *ApplicantRepository) DeleteUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error {
	if err := tx.Where(&ddl.ApplicantUserAssociation{
		ApplicantID: m.ApplicantID,
	}).Delete(&ddl.ApplicantUserAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}
