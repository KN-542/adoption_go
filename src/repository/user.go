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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

// 登録
func (u *userRepository) Insert(model *model.User) error {
	if err := u.db.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// 全件取得
func (u *userRepository) List() (*[]model.User, error) {
	var l []model.User
	if err := u.db.Find(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}
