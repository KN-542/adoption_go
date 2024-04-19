package validator

import (
	"api/src/model"
	"api/src/model/enum"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IUserValidator interface {
	CreateValidate(u *model.User) error
	CreateGroupValidate(u *model.UserGroupRequest) error
	CreateScheduleValidate(u *model.UserScheduleRequest) error
	LoginValidate(u *model.User) error
	MFAValidate(u *model.UserMFA) error
	PasswordChangeValidate(u *model.User) error
	HashKeyValidate(u *model.User) error
	ScheduleHashKeyValidate(u *model.UserSchedule) error
	LoginApplicantValidate(u *model.Applicant) error
	HashKeyValidateApplicant(u *model.Applicant) error
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
			is.Email,
		),
		validation.Field(
			&u.RoleID,
			validation.Required,
		),
	)
}
func (v *userValidator) CreateGroupValidate(u *model.UserGroupRequest) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30*3),
		),
	)
}
func (v *userValidator) CreateScheduleValidate(u *model.UserScheduleRequest) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.FreqID,
			validation.By(validateUintRange(0, uint(enum.FREQ_NONE))),
		),
		validation.Field(
			&u.Start,
			validation.Required,
		),
		validation.Field(
			&u.End,
			validation.Required,
		),
		validation.Field(
			&u.Title,
			validation.Required,
			validation.Length(1, 30),
		),
	)
}
func (v *userValidator) LoginValidate(u *model.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Email,
			validation.Required,
			validation.Length(1, 50),
			is.Email,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
		),
	)
}
func (v *userValidator) MFAValidate(u *model.UserMFA) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
		validation.Field(
			&u.Code,
			validation.Required,
			validation.Length(6, 6),
			is.UTFNumeric,
		),
	)
}

func (v *userValidator) PasswordChangeValidate(u *model.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
		),
		validation.Field(
			&u.InitPassword,
			validation.Required,
			validation.Length(8, 16),
			is.Alphanumeric,
			validation.By(func(value interface{}) error {
				initPassword, _ := value.(string)
				if initPassword == u.Password {
					return errors.New("password cannot be the same as the initial password")
				}
				return nil
			}),
		),
	)
}

func (v *userValidator) HashKeyValidate(u *model.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *userValidator) ScheduleHashKeyValidate(u *model.UserSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *userValidator) HashKeyValidateApplicant(u *model.Applicant) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *userValidator) LoginApplicantValidate(u *model.Applicant) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Email,
			validation.Required,
			validation.Length(1, 50),
			is.Email,
		),
	)
}
