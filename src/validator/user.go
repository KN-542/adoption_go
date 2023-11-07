package validator

import (
	"api/src/model"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IUserValidator interface {
	CreateValidate(u *model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (v *userValidator) CreateValidate(u *model.User) error {
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
			validation.Match(regexp.MustCompile(
				`^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*.)+[a-zA-Z]{2,}$`,
			)),
		),
		validation.Field(
			&u.RoleID,
			validation.Required,
		),
	)
}
