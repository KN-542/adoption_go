package repository

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"api/src/model/enum"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// ログイン認証
	Login(m *ddl.User) ([]entity.Login, error)
	// 取得
	Get(m *ddl.User) (*entity.User, error)
	// 登録
	Insert(tx *gorm.DB, m *ddl.User) (*entity.User, error)
	// 検索
	List() ([]ddl.UserResponse, error)
	// 更新
	Update(tx *gorm.DB, m *ddl.User) error
	// 削除
	Delete(tx *gorm.DB, m *ddl.User) error
	// チーム登録
	InsertTeam(tx *gorm.DB, m *ddl.Team) (*entity.Team, error)
	// チーム紐づけ登録
	InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error
	// ユーザー予定紐づけ登録
	InsertScheduleAssociation(tx *gorm.DB, m *ddl.UserScheduleAssociation) error
	// ユーザー予定紐づけ削除
	DeleteScheduleAssociation(tx *gorm.DB, m *ddl.UserScheduleAssociation) error
	// ユーザー予定紐づけ取得_ユーザー予定ID
	GetUserScheduleAssociationByScheduleID(m *ddl.UserScheduleAssociation) ([]ddl.UserScheduleAssociationWithName, error)
	// ユーザー基本情報取得
	GetUserBasicByHashKeys(hashKeys []string) ([]ddl.CommonModel, error)
	// ユーザー存在確認
	ConfirmUserByHashKeys(hashKeys []string) ([]ddl.UserResponse, error)
	// パスワード変更
	PasswordChange(tx *gorm.DB, m *ddl.User) error
	// 初期パスワード一致確認
	ConfirmInitPassword(m *ddl.User) (*int8, error)
	// 初期パスワード一致確認2
	ConfirmInitPassword2(m *ddl.User) (*string, error)
	// メールアドレス重複チェック
	EmailDuplCheck(m *ddl.User) error
	// 検索(チーム)
	SearchTeam() ([]ddl.TeamResponse, error)
	// スケジュール登録
	InsertSchedule(tx *gorm.DB, m *ddl.UserSchedule) (*uint, error)
	// スケジュール更新
	UpdateSchedule(tx *gorm.DB, m *ddl.UserSchedule) (*uint, error)
	// スケジュール一覧
	ListSchedule() ([]ddl.UserScheduleResponse, error)
	// スケジュール取得
	GetSchedule(m *ddl.UserSchedule) (*ddl.UserSchedule, error)
	// スケジュール削除
	DeleteSchedule(tx *gorm.DB, m *ddl.UserSchedule) error
	// 過去のスケジュールを更新
	UpdatePastSchedule(tx *gorm.DB, m *ddl.UserSchedule) error
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
func (u *UserRepository) List() ([]ddl.UserResponse, error) {
	var l []ddl.UserResponse

	query := u.db.Model(&ddl.User{}).
		Select("t_user.hash_key, t_user.name, t_user.email, t_user.role_id, m_role.name_ja as role_name_ja").
		Joins("left join m_role on t_user.role_id = m_role.id")

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
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
			UpdatedAt: m.UpdatedAt,
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

// チーム紐づけ登録
func (u *UserRepository) InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ登録
func (u *UserRepository) InsertScheduleAssociation(tx *gorm.DB, m *ddl.UserScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ削除
func (u *UserRepository) DeleteScheduleAssociation(tx *gorm.DB, m *ddl.UserScheduleAssociation) error {
	if err := tx.Where(&ddl.UserScheduleAssociation{
		UserScheduleID: m.UserScheduleID,
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ取得_ユーザー予定ID
func (u *UserRepository) GetUserScheduleAssociationByScheduleID(m *ddl.UserScheduleAssociation) ([]ddl.UserScheduleAssociationWithName, error) {
	var l []ddl.UserScheduleAssociationWithName
	if err := u.db.Model(&ddl.UserScheduleAssociationWithName{}).
		Select("t_user.name as name").
		Joins("left join t_user on t_user_schedule_association.user_id = t_user.name").
		Where(
			&ddl.UserScheduleAssociation{
				UserScheduleID: m.UserScheduleID,
			},
		).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// ユーザー基本情報取得
func (u *UserRepository) GetUserBasicByHashKeys(hashKeys []string) ([]ddl.CommonModel, error) {
	if len(hashKeys) == 0 {
		return nil, nil
	}
	var l []ddl.CommonModel

	query := u.db.Model(&ddl.User{}).
		Select("t_user.id, t_user.hash_key")

	query = query.Where("t_user.hash_key IN ?", hashKeys)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// ユーザー存在確認
func (u *UserRepository) ConfirmUserByHashKeys(hashKeys []string) ([]ddl.UserResponse, error) {
	if len(hashKeys) == 0 {
		return nil, nil
	}
	var l []ddl.UserResponse

	query := u.db.Model(&ddl.User{}).
		Select("t_user.hash_key, t_user.name, t_user.email, t_user.role_id, m_role.name_ja as role_name_ja").
		Joins("left join m_role on t_user.role_id = m_role.id")

	query = query.Where("t_user.hash_key IN ?", hashKeys)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// パスワード変更
func (u *UserRepository) PasswordChange(tx *gorm.DB, m *ddl.User) error {
	user := ddl.User{Password: m.Password}
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

// 初期パスワード一致確認
func (u *UserRepository) ConfirmInitPassword(m *ddl.User) (*int8, error) {
	var confirm ddl.User
	if err := u.db.Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&confirm).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if confirm.Password == confirm.InitPassword {
		res := int8(enum.PASSWORD_CHANGE_REQUIRED)
		return &res, nil
	}
	res := int8(enum.PASSWORD_CHANGE_UNREQUIRED)
	return &res, nil
}

// 初期パスワード一致確認2
func (u *UserRepository) ConfirmInitPassword2(m *ddl.User) (*string, error) {
	var confirm ddl.User
	if err := u.db.Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&confirm).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &confirm.InitPassword, nil
}

// メールアドレス重複チェック
func (u *UserRepository) EmailDuplCheck(m *ddl.User) error {
	var count int64
	if err := u.db.Model(&ddl.User{}).Where(
		&ddl.User{
			Email: m.Email,
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

// 検索(チーム)
func (u *UserRepository) SearchTeam() ([]ddl.TeamResponse, error) {
	var l []ddl.TeamResponse

	query := u.db.Model(&ddl.Team{}).
		Select("t_team.hash_key, t_team.name, t_team.users")

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// スケジュール登録
func (u *UserRepository) InsertSchedule(tx *gorm.DB, m *ddl.UserSchedule) (*uint, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	id := uint(m.ID)
	return &id, nil
}

// スケジュール更新
func (u *UserRepository) UpdateSchedule(tx *gorm.DB, m *ddl.UserSchedule) (*uint, error) {
	if err := tx.Model(&ddl.UserSchedule{}).Where(
		&ddl.UserSchedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	id := uint(m.ID)
	return &id, nil
}

// スケジュール一覧
func (u *UserRepository) ListSchedule() ([]ddl.UserScheduleResponse, error) {
	var res []ddl.UserScheduleResponse

	query := u.db.Model(&ddl.UserSchedule{}).
		Select("t_user_schedule.*, m_calendar_freq_status.freq").
		Joins("left join m_calendar_freq_status on t_user_schedule.freq_id = m_calendar_freq_status.id")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// スケジュール取得
func (u *UserRepository) GetSchedule(m *ddl.UserSchedule) (*ddl.UserSchedule, error) {
	var res ddl.UserSchedule
	if err := u.db.Where(
		&ddl.UserSchedule{
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

// スケジュール削除
func (u *UserRepository) DeleteSchedule(tx *gorm.DB, m *ddl.UserSchedule) error {
	if err := tx.Where(m).Delete(&ddl.UserSchedule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 過去のスケジュールを更新
func (u *UserRepository) UpdatePastSchedule(tx *gorm.DB, m *ddl.UserSchedule) error {
	schedule := ddl.UserSchedule{
		Start: m.Start,
		End:   m.End,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: m.UpdatedAt,
		},
	}
	if err := tx.Model(&ddl.UserSchedule{}).Where(
		&ddl.UserSchedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(schedule).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}
