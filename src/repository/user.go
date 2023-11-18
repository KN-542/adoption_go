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
	Insert(m *model.User) error
	// 全件取得
	List() (*[]model.User, error)
	// ログイン認証
	Login(m *model.User) ([]model.User, error)
	// パスワード変更
	PasswordChange(m *model.User) error
	// 初期パスワード一致確認
	ConfirmInitPassword(m *model.User) (*int8, error)
	// 初期パスワード一致確認2
	ConfirmInitPassword2(m *model.User) (*string, error)
	// メールアドレス重複チェック
	EmailDuplCheck(m *model.User) error
	// ユーザー存在確認
	UserCheck(m *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// 登録
func (u *UserRepository) Insert(m *model.User) error {
	if err := u.db.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 全件取得
func (u *UserRepository) List() (*[]model.User, error) {
	var l []model.User
	if err := u.db.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &l, nil
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
func (u *UserRepository) PasswordChange(m *model.User) error {
	user := model.User{Password: m.Password}
	if err := u.db.Model(&model.User{}).Where(
		&model.User{
			HashKey: m.HashKey,
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
			HashKey: m.HashKey,
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
			HashKey: m.HashKey,
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
		return fmt.Errorf("Duplicate Email Address")
	}

	return nil
}

// ユーザー存在確認
func (u *UserRepository) UserCheck(m *model.User) error {
	var count int64
	if err := u.db.Model(&model.User{}).Where(
		&model.User{
			HashKey: m.HashKey,
		},
	).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if count == 0 {
		return fmt.Errorf("User Not Exist")
	} else if count > 1 {
		return fmt.Errorf("User Duplicate Exist")
	}

	return nil
}
