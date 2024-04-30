package repository

import (
	"api/src/model"
	enum "api/src/model/enum"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *model.User) error
	// 検索
	List() ([]model.UserResponse, error)
	// 取得
	Get(m *model.User) (*model.User, error)
	// 更新
	Update(tx *gorm.DB, m *model.User) error
	// 削除
	Delete(tx *gorm.DB, m *model.User) error
	// ユーザーグループ紐づけ登録
	InsertGroupAssociation(tx *gorm.DB, m *model.UserGroupAssociation) error
	// ユーザー予定紐づけ登録
	InsertScheduleAssociation(tx *gorm.DB, m *model.UserScheduleAssociation) error
	// ユーザー予定紐づけ削除
	DeleteScheduleAssociation(tx *gorm.DB, m *model.UserScheduleAssociation) error
	// ユーザー予定紐づけ取得_ユーザー予定ID
	GetUserScheduleAssociationByScheduleID(m *model.UserScheduleAssociation) ([]model.UserScheduleAssociationWithName, error)
	// ユーザー基本情報取得
	GetUserBasicByHashKeys(hashKeys []string) ([]model.CommonModel, error)
	// ユーザー存在確認
	ConfirmUserByHashKeys(hashKeys []string) ([]model.UserResponse, error)
	// ログイン認証
	Login(m *model.User) ([]model.User, error)
	// パスワード変更
	PasswordChange(tx *gorm.DB, m *model.User) error
	// 初期パスワード一致確認
	ConfirmInitPassword(m *model.User) (*int8, error)
	// 初期パスワード一致確認2
	ConfirmInitPassword2(m *model.User) (*string, error)
	// メールアドレス重複チェック
	EmailDuplCheck(m *model.User) error
	// 検索(ユーザーグループ)
	SearchGroup() ([]model.UserGroupResponse, error)
	// グループ登録
	InsertGroup(tx *gorm.DB, m *model.UserGroup) (*uint, error)
	// スケジュール登録
	InsertSchedule(tx *gorm.DB, m *model.UserSchedule) (*uint, error)
	// スケジュール更新
	UpdateSchedule(tx *gorm.DB, m *model.UserSchedule) (*uint, error)
	// スケジュール一覧
	ListSchedule() ([]model.UserScheduleResponse, error)
	// スケジュール取得
	GetSchedule(m *model.UserSchedule) (*model.UserSchedule, error)
	// スケジュール削除
	DeleteSchedule(tx *gorm.DB, m *model.UserSchedule) error
	// 過去のスケジュールを更新
	UpdatePastSchedule(tx *gorm.DB, m *model.UserSchedule) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// 登録
func (u *UserRepository) Insert(tx *gorm.DB, m *model.User) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索
func (u *UserRepository) List() ([]model.UserResponse, error) {
	var l []model.UserResponse

	query := u.db.Model(&model.User{}).
		Select("t_user.hash_key, t_user.name, t_user.email, t_user.role_id, m_role.name_ja as role_name_ja").
		Joins("left join m_role on t_user.role_id = m_role.id")

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// 取得
func (u *UserRepository) Get(m *model.User) (*model.User, error) {
	var res model.User
	if err := u.db.Where(
		&model.User{
			AbstractTransactionModel: model.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 更新
func (u *UserRepository) Update(tx *gorm.DB, m *model.User) error {
	user := model.User{
		Name:         m.Name,
		Email:        m.Email,
		Password:     m.Password,
		RoleID:       m.RoleID,
		RefreshToken: m.RefreshToken,
		AbstractTransactionModel: model.AbstractTransactionModel{
			UpdatedAt: m.UpdatedAt,
		},
	}
	if err := tx.Model(&model.User{}).Where(
		&model.User{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) Delete(tx *gorm.DB, m *model.User) error {
	if err := tx.Where(&model.User{
		AbstractTransactionModel: model.AbstractTransactionModel{
			HashKey: m.HashKey,
		},
	}).Delete(&model.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザーグループ紐づけ登録
func (u *UserRepository) InsertGroupAssociation(tx *gorm.DB, m *model.UserGroupAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ登録
func (u *UserRepository) InsertScheduleAssociation(tx *gorm.DB, m *model.UserScheduleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ削除
func (u *UserRepository) DeleteScheduleAssociation(tx *gorm.DB, m *model.UserScheduleAssociation) error {
	if err := tx.Where(&model.UserScheduleAssociation{
		UserScheduleID: m.UserScheduleID,
	}).Delete(&model.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ユーザー予定紐づけ取得_ユーザー予定ID
func (u *UserRepository) GetUserScheduleAssociationByScheduleID(m *model.UserScheduleAssociation) ([]model.UserScheduleAssociationWithName, error) {
	var l []model.UserScheduleAssociationWithName
	if err := u.db.Model(&model.UserScheduleAssociationWithName{}).
		Select("t_user.name as name").
		Joins("left join t_user on t_user_schedule_association.user_id = t_user.name").
		Where(
			&model.UserScheduleAssociation{
				UserScheduleID: m.UserScheduleID,
			},
		).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// ユーザー基本情報取得
func (u *UserRepository) GetUserBasicByHashKeys(hashKeys []string) ([]model.CommonModel, error) {
	if len(hashKeys) == 0 {
		return nil, nil
	}
	var l []model.CommonModel

	query := u.db.Model(&model.User{}).
		Select("t_user.id, t_user.hash_key")

	query = query.Where("t_user.hash_key IN ?", hashKeys)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// ユーザー存在確認
func (u *UserRepository) ConfirmUserByHashKeys(hashKeys []string) ([]model.UserResponse, error) {
	if len(hashKeys) == 0 {
		return nil, nil
	}
	var l []model.UserResponse

	query := u.db.Model(&model.User{}).
		Select("t_user.hash_key, t_user.name, t_user.email, t_user.role_id, m_role.name_ja as role_name_ja").
		Joins("left join m_role on t_user.role_id = m_role.id")

	query = query.Where("t_user.hash_key IN ?", hashKeys)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// ログイン認証
func (u *UserRepository) Login(m *model.User) ([]model.User, error) {
	var res []model.User
	if err := u.db.Where(
		&model.User{
			Email: m.Email,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// パスワード変更
func (u *UserRepository) PasswordChange(tx *gorm.DB, m *model.User) error {
	user := model.User{Password: m.Password}
	if err := tx.Model(&model.User{}).Where(
		&model.User{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) ConfirmInitPassword(m *model.User) (*int8, error) {
	var confirm model.User
	if err := u.db.Where(
		&model.User{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) ConfirmInitPassword2(m *model.User) (*string, error) {
	var confirm model.User
	if err := u.db.Where(
		&model.User{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) EmailDuplCheck(m *model.User) error {
	var count int64
	if err := u.db.Model(&model.User{}).Where(
		&model.User{
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

// 検索(ユーザーグループ)
func (u *UserRepository) SearchGroup() ([]model.UserGroupResponse, error) {
	var l []model.UserGroupResponse

	query := u.db.Model(&model.UserGroup{}).
		Select("t_user_group.hash_key, t_user_group.name, t_user_group.users")

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// グループ登録
func (u *UserRepository) InsertGroup(tx *gorm.DB, m *model.UserGroup) (*uint, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	id := uint(m.ID)
	return &id, nil
}

// スケジュール登録
func (u *UserRepository) InsertSchedule(tx *gorm.DB, m *model.UserSchedule) (*uint, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	id := uint(m.ID)
	return &id, nil
}

// スケジュール更新
func (u *UserRepository) UpdateSchedule(tx *gorm.DB, m *model.UserSchedule) (*uint, error) {
	if err := tx.Model(&model.UserSchedule{}).Where(
		&model.UserSchedule{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) ListSchedule() ([]model.UserScheduleResponse, error) {
	var res []model.UserScheduleResponse

	query := u.db.Model(&model.UserSchedule{}).
		Select("t_user_schedule.*, m_calendar_freq_status.freq").
		Joins("left join m_calendar_freq_status on t_user_schedule.freq_id = m_calendar_freq_status.id")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// スケジュール取得
func (u *UserRepository) GetSchedule(m *model.UserSchedule) (*model.UserSchedule, error) {
	var res model.UserSchedule
	if err := u.db.Where(
		&model.UserSchedule{
			AbstractTransactionModel: model.AbstractTransactionModel{
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
func (u *UserRepository) DeleteSchedule(tx *gorm.DB, m *model.UserSchedule) error {
	if err := tx.Where(m).Delete(&model.UserSchedule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 過去のスケジュールを更新
func (u *UserRepository) UpdatePastSchedule(tx *gorm.DB, m *model.UserSchedule) error {
	schedule := model.UserSchedule{
		Start: m.Start,
		End:   m.End,
		AbstractTransactionModel: model.AbstractTransactionModel{
			UpdatedAt: m.UpdatedAt,
		},
	}
	if err := tx.Model(&model.UserSchedule{}).Where(
		&model.UserSchedule{
			AbstractTransactionModel: model.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(schedule).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}
