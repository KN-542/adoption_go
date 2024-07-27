package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"api/src/model/static"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// ログイン認証
	Login(m *ddl.User) ([]entity.Login, error)
	// 登録
	Insert(tx *gorm.DB, m *ddl.User) (*entity.User, error)
	// 検索
	Search(m *dto.SearchUser) ([]entity.SearchUser, error)
	// 検索_同一企業
	SearchByCompany(m *ddl.User) ([]entity.SearchUser, error)
	// 取得
	Get(m *ddl.User) (*entity.User, error)
	// 取得_PK
	GetByPrimary(m *ddl.User) (*entity.User, error)
	// 更新
	Update(tx *gorm.DB, m *ddl.User) error
	// 削除
	Delete(tx *gorm.DB, m *ddl.User) error
	// チーム登録
	InsertTeam(tx *gorm.DB, m *ddl.Team) (*entity.Team, error)
	// チーム検索
	SearchTeam(m *ddl.Team) ([]*entity.SearchTeam, error)
	// チーム取得
	GetTeam(m *ddl.Team) (*entity.Team, error)
	// チーム取得_PK
	GetTeamByPrimary(m *ddl.Team) (*entity.Team, error)
	// チーム更新
	UpdateTeam(tx *gorm.DB, m *ddl.Team) (*entity.Team, error)
	// チーム削除
	DeleteTeam(tx *gorm.DB, m *ddl.Team) error
	// チーム検索_同一企業
	SearchTeamByCompany(m *ddl.Team) ([]entity.SearchTeam, error)
	// 予定登録
	InsertSchedule(tx *gorm.DB, m *ddl.Schedule) (*uint64, error)
	// 予定検索
	SearchSchedule(m *ddl.Schedule) ([]*entity.Schedule, error)
	// 予定取得
	GetSchedule(m *ddl.Schedule) (*entity.Schedule, error)
	// 予定取得_PK
	GetScheduleByPrimary(m *ddl.Schedule) (*entity.Schedule, error)
	// 予定更新
	UpdateSchedule(tx *gorm.DB, m *ddl.Schedule) (*uint64, error)
	// 予定削除
	DeleteSchedule(tx *gorm.DB, m *ddl.Schedule) error
	// チーム紐づけ登録
	InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error
	// チーム紐づけ一括登録
	InsertsTeamAssociation(tx *gorm.DB, m []*ddl.TeamAssociation) error
	// チーム紐づけ一覧取得
	ListTeamAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error)
	// ユーザー紐づけ一覧取得
	ListUserAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error)
	// 所属チーム一覧取得
	ListBelongTeam(m *ddl.TeamAssociation) ([]entity.Team, error)
	// チーム紐づけ削除
	DeleteTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error
	// チーム毎ステータスイベント取得
	StatusEventsByTeam(m *ddl.Team) ([]entity.StatusEventsByTeam, error)
	// チーム面接毎イベント取得
	InterviewEventsByTeam(m *ddl.Team) ([]entity.InterviewEventsByTeam, error)
	// 予定紐づけ登録
	InsertScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error
	// 予定紐づけ一括登録
	InsertsScheduleAssociation(tx *gorm.DB, m []*ddl.ScheduleAssociation) error
	// 予定紐づけ取得
	ListUserScheduleAssociation(m *ddl.ScheduleAssociation) ([]entity.Schedule, error)
	// 予定毎ユーザー紐づけ取得
	SearchScheduleUserAssociation(m *ddl.ScheduleAssociation) ([]entity.ScheduleAssociation, error)
	// 予定紐づけ削除
	DeleteScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error
	// 選考状況登録
	InsertSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// 選考状況一括登録
	InsertsSelectStatus(tx *gorm.DB, m []*ddl.SelectStatus) (*entity.ApplicantStatusList, error)
	// 選考状況削除
	DeleteSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// ユーザー紐づけ登録
	InsertUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error
	// ユーザー紐づけ取得
	GetUserAssociation(m []uint64) ([]entity.ApplicantUserAssociation, error)
	// ユーザー紐づけ削除
	DeleteUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error
	// イベント一括登録
	InsertsEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error
	// イベント削除
	DeleteEventAssociation(tx *gorm.DB, m *ddl.TeamEvent) error
	// 面接毎イベント一括登録
	InsertsEventEachInterviewAssociation(tx *gorm.DB, m []*ddl.TeamEventEachInterview) error
	// 面接毎イベント削除
	DeleteEventEachInterviewAssociation(tx *gorm.DB, m *ddl.TeamEventEachInterview) error
	// メールアドレス重複チェック
	EmailDuplCheck(m *ddl.User) error
	// メールアドレス重複チェック_管理者
	EmailDuplCheckManagement(m *ddl.User, teams []uint64) error
	// ID取得
	GetIDs(m []string) ([]uint64, error)
	// チームID取得
	GetTeamIDs(m []string) ([]uint64, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// ログイン認証
func (u *UserRepository) Login(m *ddl.User) ([]entity.Login, error) {
	var res []entity.Login

	if err := u.db.Model(&entity.Login{}).
		Select(`
			t_user.hash_key,
			t_user.name,
			t_user.company_id,
			t_user.password,
			t_user.init_password,
			t_user.role_id,
			t_user.user_type
		`).
		Where(
			&ddl.User{
				Email: m.Email,
			},
		).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 登録
func (u *UserRepository) Insert(tx *gorm.DB, m *ddl.User) (*entity.User, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &entity.User{
		User: *m,
	}, nil
}

// 検索
func (u *UserRepository) Search(m *dto.SearchUser) ([]entity.SearchUser, error) {
	var l []entity.SearchUser

	// TODO
	query := u.db.Model(&entity.SearchUser{}).
		Select(`
			t_user.hash_key,
			t_user.name,
			t_user.email,
			t_role.name as role_name
		`).
		Joins("left join t_role on t_role.id = t_user.role_id").
		Where("t_user.company_id = ?", m.CompanyID)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// 検索_同一企業
func (u *UserRepository) SearchByCompany(m *ddl.User) ([]entity.SearchUser, error) {
	var l []entity.SearchUser

	if err := u.db.Model(&entity.SearchUser{}).
		Select("hash_key, name, email").
		Where("company_id = ?", m.CompanyID).
		Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// 取得
func (u *UserRepository) Get(m *ddl.User) (*entity.User, error) {
	var res entity.User
	if err := u.db.Where(
		&ddl.User{
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

// 取得_PK
func (u *UserRepository) GetByPrimary(m *ddl.User) (*entity.User, error) {
	var res entity.User
	if err := u.db.Where(
		&ddl.User{
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

// 更新
func (u *UserRepository) Update(tx *gorm.DB, m *ddl.User) error {
	user := ddl.User{
		Name:         m.Name,
		Email:        m.Email,
		Password:     m.Password,
		RoleID:       m.RoleID,
		RefreshToken: m.RefreshToken,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.User{}).Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(user).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 削除
func (u *UserRepository) Delete(tx *gorm.DB, m *ddl.User) error {
	if err := tx.Where(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: m.HashKey,
		},
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム登録
func (u *UserRepository) InsertTeam(tx *gorm.DB, m *ddl.Team) (*entity.Team, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.Team{
		Team: *m,
	}, nil
}

// チーム検索
func (u *UserRepository) SearchTeam(m *ddl.Team) ([]*entity.SearchTeam, error) {
	var l []*entity.SearchTeam

	query := u.db.Table("t_team").
		Select(`
			t_team.id,
			t_team.hash_key,
			t_team.name
		`).
		Where("t_team.company_id = ?", m.CompanyID)

	if err := query.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name")
	}).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// チーム取得
func (u *UserRepository) GetTeam(m *ddl.Team) (*entity.Team, error) {
	var res entity.Team
	if err := u.db.Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name, email")
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// チーム取得_PK
func (u *UserRepository) GetTeamByPrimary(m *ddl.Team) (*entity.Team, error) {
	var res entity.Team
	if err := u.db.Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: m.ID,
			},
		},
	).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name, email")
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// チーム更新
func (u *UserRepository) UpdateTeam(tx *gorm.DB, m *ddl.Team) (*entity.Team, error) {
	team := ddl.Team{
		Name:           m.Name,
		NumOfInterview: m.NumOfInterview,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.Team{}).Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(team).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var updatedTeam ddl.Team
	if err := tx.Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&updatedTeam).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.Team{
		Team: updatedTeam,
	}, nil
}

// チーム削除
func (u *UserRepository) DeleteTeam(tx *gorm.DB, m *ddl.Team) error {
	if err := tx.Where(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: m.HashKey,
		},
	}).Delete(&ddl.Team{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム検索_同一企業
func (u *UserRepository) SearchTeamByCompany(m *ddl.Team) ([]entity.SearchTeam, error) {
	var l []entity.SearchTeam

	if err := u.db.Model(&entity.SearchTeam{}).
		Select("hash_key, name").
		Where("company_id = ?", m.CompanyID).
		Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// チーム毎ステータスイベント取得
func (u *UserRepository) StatusEventsByTeam(m *ddl.Team) ([]entity.StatusEventsByTeam, error) {
	var res []entity.StatusEventsByTeam

	query := u.db.Table("t_team_event").
		Select(`
			m_select_status_event.hash_key as event_hash_key,
			m_select_status_event.desc_ja as desc_ja,
			m_select_status_event.desc_en as desc_en,
			t_select_status.hash_key as select_status_hash_key,
			t_select_status.status_name as status_name
		`).
		Joins("left join t_select_status on t_team_event.status_id = t_select_status.id").
		Joins("left join m_select_status_event on t_team_event.event_id = m_select_status_event.id").
		Where("t_team_event.team_id = ?", m.ID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// チーム面接毎イベント取得
func (u *UserRepository) InterviewEventsByTeam(m *ddl.Team) ([]entity.InterviewEventsByTeam, error) {
	var res []entity.InterviewEventsByTeam

	query := u.db.Table("t_team_event_each_interview").
		Select(`
			t_team_event_each_interview.num_of_interview as num_of_interview,
			t_select_status.hash_key as select_status_hash_key,
			t_select_status.status_name as status_name
		`).
		Joins("left join t_select_status on t_team_event_each_interview.status_id = t_select_status.id").
		Where("t_team_event_each_interview.team_id = ?", m.ID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 予定登録
func (u *UserRepository) InsertSchedule(tx *gorm.DB, m *ddl.Schedule) (*uint64, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &m.ID, nil
}

// 予定検索
func (u *UserRepository) SearchSchedule(m *ddl.Schedule) ([]*entity.Schedule, error) {
	var res []*entity.Schedule

	query := u.db.Table("t_schedule").
		Select("t_schedule.*, m_schedule_freq_status.freq").
		Joins("left join m_schedule_freq_status on t_schedule.freq_id = m_schedule_freq_status.id").
		Where("t_schedule.company_id = ?", m.CompanyID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 予定取得_PK
func (u *UserRepository) GetScheduleByPrimary(m *ddl.Schedule) (*entity.Schedule, error) {
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

// 予定取得
func (u *UserRepository) GetSchedule(m *ddl.Schedule) (*entity.Schedule, error) {
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
func (u *UserRepository) UpdateSchedule(tx *gorm.DB, m *ddl.Schedule) (*uint64, error) {
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
		Title:  m.Title,
		FreqID: m.FreqID,
		Start:  m.Start,
		End:    m.End,
	}).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &m.ID, nil
}

// 予定削除
func (u *UserRepository) DeleteSchedule(tx *gorm.DB, m *ddl.Schedule) error {
	if err := tx.Where(m).Delete(&ddl.Schedule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム紐づけ登録
func (u *UserRepository) InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム紐づけ一括登録
func (u *UserRepository) InsertsTeamAssociation(tx *gorm.DB, m []*ddl.TeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 所属チーム一覧取得
func (u *UserRepository) ListBelongTeam(m *ddl.TeamAssociation) ([]entity.Team, error) {
	var res []entity.Team
	if err := u.db.
		Joins("left join t_team_association on t_team_association.team_id = t_team.id").
		Where("t_team_association.user_id = ?", m.UserID).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("team not exist")
	}

	return res, nil
}

// チーム紐づけ一覧取得
func (u *UserRepository) ListTeamAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error) {
	var res []entity.TeamAssociation
	if err := u.db.Where(
		&ddl.TeamAssociation{
			UserID: m.UserID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// ユーザー紐づけ一覧取得
func (u *UserRepository) ListUserAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error) {
	var res []entity.TeamAssociation
	if err := u.db.Where(
		&ddl.TeamAssociation{
			TeamID: m.TeamID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// チーム紐づけ削除
func (u *UserRepository) DeleteTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error {
	if err := tx.Where(m).Delete(&ddl.TeamAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ登録
func (u *UserRepository) InsertScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ一括登録
func (u *UserRepository) InsertsScheduleAssociation(tx *gorm.DB, m []*ddl.ScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 予定紐づけ一覧取得
func (u *UserRepository) ListUserScheduleAssociation(m *ddl.ScheduleAssociation) ([]entity.Schedule, error) {
	var res []entity.Schedule
	if err := u.db.Table("t_schedule").
		Joins(`
			left join
				t_schedule_association
			on
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
func (u *UserRepository) SearchScheduleUserAssociation(m *ddl.ScheduleAssociation) ([]entity.ScheduleAssociation, error) {
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

// 予定紐づけ削除
func (u *UserRepository) DeleteScheduleAssociation(tx *gorm.DB, m *ddl.ScheduleAssociation) error {
	if err := tx.Where(&ddl.ScheduleAssociation{
		ScheduleID: m.ScheduleID,
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント登録
func (u *UserRepository) InsertEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント一括登録
func (u *UserRepository) InsertsEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント削除
func (u *UserRepository) DeleteEventAssociation(tx *gorm.DB, m *ddl.TeamEvent) error {
	if err := tx.Model(&ddl.TeamEvent{}).Where(&ddl.TeamEvent{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamEvent{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎イベント一括登録
func (u *UserRepository) InsertsEventEachInterviewAssociation(tx *gorm.DB, m []*ddl.TeamEventEachInterview) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎イベント削除
func (u *UserRepository) DeleteEventEachInterviewAssociation(tx *gorm.DB, m *ddl.TeamEventEachInterview) error {
	if err := tx.Model(&ddl.TeamEventEachInterview{}).Where(&ddl.TeamEventEachInterview{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamEventEachInterview{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 選考状況登録
func (u *UserRepository) InsertSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 選考状況一括登録
func (u *UserRepository) InsertsSelectStatus(tx *gorm.DB, m []*ddl.SelectStatus) (*entity.ApplicantStatusList, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.ApplicantStatusList{
		List: m,
	}, nil
}

// 選考状況削除
func (u *UserRepository) DeleteSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error {
	if err := tx.Model(&ddl.SelectStatus{}).Where(&ddl.SelectStatus{
		TeamID: m.TeamID,
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー紐づけ登録
func (u *UserRepository) InsertUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー紐づけ取得
func (a *UserRepository) GetUserAssociation(m []uint64) ([]entity.ApplicantUserAssociation, error) {
	var res []entity.ApplicantUserAssociation

	query := a.db.Table("t_applicant_user_association").
		Where("user_id IN ?", m)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// ユーザー紐づけ削除
func (u *UserRepository) DeleteUserAssociation(tx *gorm.DB, m *ddl.ApplicantUserAssociation) error {
	if err := tx.Where(&ddl.ApplicantUserAssociation{
		ApplicantID: m.ApplicantID,
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// メールアドレス重複チェック
func (u *UserRepository) EmailDuplCheck(m *ddl.User) error {
	var count int64
	if err := u.db.Model(&ddl.User{}).Where(
		&ddl.User{
			Email:    m.Email,
			UserType: static.LOGIN_TYPE_ADMIN,
		},
	).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("duplicate Email Address")
	}

	return nil
}

// メールアドレス重複チェック_管理者
func (u *UserRepository) EmailDuplCheckManagement(m *ddl.User, teams []uint64) error {
	if len(teams) == 0 {
		return fmt.Errorf("not exist teams")
	}

	var count int64

	if err := u.db.Model(&ddl.User{}).
		Joins(`left join t_team_association on t_team_association.user_id = t_user.id`).
		Where(
			&ddl.User{
				Email:    m.Email,
				UserType: static.LOGIN_TYPE_MANAGEMENT,
			},
		).
		Where("t_team_association.team_id IN ?", teams).
		Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("duplicate Email")
	}

	return nil
}

// ID取得
func (u *UserRepository) GetIDs(m []string) ([]uint64, error) {
	var res []entity.User
	if err := u.db.Table("t_user").
		Select("id").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var IDs []uint64
	for _, row := range res {
		IDs = append(IDs, row.ID)
	}

	return IDs, nil
}

// チームID取得
func (u *UserRepository) GetTeamIDs(m []string) ([]uint64, error) {
	var res []entity.Team
	if err := u.db.Model(&ddl.Team{}).
		Select("id").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var IDs []uint64
	for _, row := range res {
		IDs = append(IDs, row.ID)
	}

	return IDs, nil
}
