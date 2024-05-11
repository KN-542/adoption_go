package validator

import (
	"api/src/model/ddl"
	"api/src/model/enum"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// TODO 名前の末尾からValidateを除去
type IUserValidator interface {
	CreateValidate(u *ddl.User) error
	CreateTeamValidate(u *ddl.TeamRequest) error
	CreateScheduleValidate(u *ddl.UserScheduleRequest) error
	ScheduleHashKeyValidate(u *ddl.UserSchedule) error
	LoginApplicantValidate(u *ddl.Applicant) error
	HashKeyValidateApplicant(u *ddl.Applicant) error
	HashKeyValidate(u *ddl.User) error
}

type UserValidator struct{}

func NewUserValidator() IUserValidator {
	return &UserValidator{}
}

func (v *UserValidator) CreateValidate(u *ddl.User) error {
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
func (v *UserValidator) CreateTeamValidate(u *ddl.TeamRequest) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30*3),
		),
	)
}
func (v *UserValidator) CreateScheduleValidate(u *ddl.UserScheduleRequest) error {
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

func (v *UserValidator) ScheduleHashKeyValidate(u *ddl.UserSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *UserValidator) HashKeyValidateApplicant(u *ddl.Applicant) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *UserValidator) LoginApplicantValidate(u *ddl.Applicant) error {
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

func (v *UserValidator) HashKeyValidate(u *ddl.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}
