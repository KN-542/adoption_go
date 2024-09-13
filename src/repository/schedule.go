package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

type IScheduleRepository interface {
	// 予定登録
	Insert(tx *gorm.DB, m *ddl.Schedule) (*uint64, error)
	// 予定一括登録
	Inserts(tx *gorm.DB, m []*ddl.Schedule) error
	// 予定検索
	Search(m *ddl.Schedule) ([]*entity.Schedule, error)
	// 予定取得
	Get(m *ddl.Schedule) (*entity.Schedule, error)
	// 予定取得_PK
	GetByPrimary(m *ddl.Schedule) (*entity.Schedule, error)
	// 予定取得_チームID
	GetByTeamID(m *ddl.Schedule) ([]entity.Schedule, error)
	// 予定更新
	Update(tx *gorm.DB, m *ddl.Schedule) error
	// 予定更新_PK
	UpdateByPrimary(tx *gorm.DB, m *ddl.Schedule) error
	// 予定削除
	Delete(tx *gorm.DB, m *ddl.Schedule) error
	// 予定一括削除
	Deletes(tx *gorm.DB, m []uint64) error
	// 予定紐づけ登録
	InsertScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error
	// 予定紐づけ一括登録
	InsertsScheduleAssociation(tx *gorm.DB, m []*ddl.ScheduleAssociation) error
	// 予定紐づけ取得
	ListUserScheduleAssociation(m *ddl.ScheduleAssociation) ([]entity.Schedule, error)
	// 予定毎ユーザー紐づけ取得
	SearchScheduleUserAssociation(m *ddl.ScheduleAssociation) ([]entity.ScheduleAssociation, error)
	// ユーザー単位予定取得
	GetScheduleByUser(m *dto.GetScheduleByUser) ([]entity.Schedule2, error)
	// 予定紐づけ削除
	DeleteScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error
}

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) IScheduleRepository {
	return &ScheduleRepository{db}
}

// 予定登録
func (u *ScheduleRepository) Insert(tx *gorm.DB, m *ddl.Schedule) (*uint64, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &m.ID, nil
}

// 予定一括登録
func (u *ScheduleRepository) Inserts(tx *gorm.DB, m []*ddl.Schedule) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 予定検索
func (u *ScheduleRepository) Search(m *ddl.Schedule) ([]*entity.Schedule, error) {
	var res []*entity.Schedule

	query := u.db.Table("t_schedule").
		Select(`
		t_schedule.id,
		t_schedule.hash_key,
		t_schedule.title,
		t_schedule.freq_id,
		t_schedule.interview_flg,
		t_schedule.start,
		t_schedule.end,
		m_schedule_freq_status.freq_name
	`).
		Joins("LEFT JOIN m_schedule_freq_status ON t_schedule.freq_id = m_schedule_freq_status.id").
		Where("t_schedule.team_id = ?", m.TeamID)

	if err := query.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name, email")
	}).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 予定取得_PK
func (u *ScheduleRepository) GetByPrimary(m *ddl.Schedule) (*entity.Schedule, error) {
	var res entity.Schedule
	if err := u.db.Where(
		&ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: m.ID,
			},
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 予定取得_チームID
func (u *ScheduleRepository) GetByTeamID(m *ddl.Schedule) ([]entity.Schedule, error) {
	var res []entity.Schedule
	if err := u.db.Where(
		&ddl.Schedule{
			TeamID: m.TeamID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 予定取得
func (u *ScheduleRepository) Get(m *ddl.Schedule) (*entity.Schedule, error) {
	var res entity.Schedule
	if err := u.db.Where(
		&ddl.Schedule{
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

// 予定更新
func (u *ScheduleRepository) Update(tx *gorm.DB, m *ddl.Schedule) error {
	if err := tx.Model(&ddl.Schedule{}).Where(
		&ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
		InterviewFlg: m.InterviewFlg,
		Title:        m.Title,
		FreqID:       m.FreqID,
		Start:        m.Start,
		End:          m.End,
	}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 予定更新_PK
func (u *ScheduleRepository) UpdateByPrimary(tx *gorm.DB, m *ddl.Schedule) error {
	if err := tx.Model(&ddl.Schedule{}).Where(
		&ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: m.ID,
			},
		},
	).Updates(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
		InterviewFlg: m.InterviewFlg,
		Title:        m.Title,
		FreqID:       m.FreqID,
		Start:        m.Start,
		End:          m.End,
	}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 予定削除
func (u *ScheduleRepository) Delete(tx *gorm.DB, m *ddl.Schedule) error {
	if err := tx.Where(m).Delete(&ddl.Schedule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定一括削除
func (u *ScheduleRepository) Deletes(tx *gorm.DB, m []uint64) error {
	if err := tx.Table("t_schedule").Where("id IN ?", m).Delete(&ddl.Schedule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ登録
func (u *ScheduleRepository) InsertScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ一括登録
func (u *ScheduleRepository) InsertsScheduleAssociation(tx *gorm.DB, m []*ddl.ScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ一覧取得
func (u *ScheduleRepository) ListUserScheduleAssociation(m *ddl.ScheduleAssociation) ([]entity.Schedule, error) {
	var res []entity.Schedule
	if err := u.db.Table("t_schedule").
		Joins(`
			LEFT JOIN
				t_schedule_association
			ON
				t_schedule_association.schedule_id = t_schedule.id
		`).
		Where("t_schedule_association.user_id = ?", m.UserID).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 予定毎ユーザー紐づけ取得
func (u *ScheduleRepository) SearchScheduleUserAssociation(m *ddl.ScheduleAssociation) ([]entity.ScheduleAssociation, error) {
	var res []entity.ScheduleAssociation
	if err := u.db.Where(
		&ddl.ScheduleAssociation{
			ScheduleID: m.ScheduleID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// ユーザー単位予定取得
func (u *ScheduleRepository) GetScheduleByUser(m *dto.GetScheduleByUser) ([]entity.Schedule2, error) {
	var res []entity.Schedule2
	if err := u.db.Table("t_schedule").
		Select(`
			t_schedule.start,
			t_schedule.end,
			t_schedule.freq_id,
			t_schedule.interview_flg
		`).
		Joins(`
			LEFT JOIN
				t_schedule_association
			ON
				t_schedule_association.schedule_id = t_schedule.id
		`).
		Where("t_schedule_association.user_id = ?", m.UserID).
		Where("t_schedule.hash_key NOT IN ?", m.RemoveScheduleHashKeys).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 予定紐づけ削除
func (u *ScheduleRepository) DeleteScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error {
	if err := tx.Table("t_schedule_association").Where(&ddl.ScheduleAssociation{
		ScheduleID: m.ScheduleID,
	}).Delete(&ddl.ScheduleAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}
