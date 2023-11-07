package validator

import (
	"api/src/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IUserValidator interface {
	UserValidate(u model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (v *userValidator) UserValidate(u model.User) error {
	return validation.ValidateStruct(&u, validation.Field(&u.Name, validation.Required))
}
