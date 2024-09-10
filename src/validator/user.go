package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IUserValidator interface {
	// 登録
	Create(u *request.CreateUser) error
	// 登録_管理者
	CreateManagement(u *request.CreateUser) error
	// 検索
	Search(u *request.SearchUser) error
	// 取得
	Get(u *request.GetUser) error
}

type UserValidator struct{}

func NewUserValidator() IUserValidator {
	return &UserValidator{}
}

// 登録
func (v *UserValidator) Create(u *request.CreateUser) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30),
		),
		validation.Field(
			&u.Email,
			validation.Required,
			validation.Length(1, 50),
			is.Email,
		),
		validation.Field(
			&u.RoleHashKey,
			validation.Required,
		),
	)
}

// 登録_管理者
func (v *UserValidator) CreateManagement(u *request.CreateUser) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Teams,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
	)
}

// 検索
func (v *UserValidator) Search(u *request.SearchUser) error {
	return validation.ValidateStruct(
		u,
	)
}

// 取得
func (v *UserValidator) Get(u *request.GetUser) error {
	return validation.ValidateStruct(
		u,
	)
}
