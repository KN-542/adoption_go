package repository

import (
	"api/src/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// 登録
	Insert(model *model.User) error
	// 全件取得
	List() (*[]model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// 登録
func (u *UserRepository) Insert(model *model.User) error {
	if err := u.db.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// 全件取得
func (u *UserRepository) List() (*[]model.User, error) {
	var l []model.User
	if err := u.db.Find(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}
