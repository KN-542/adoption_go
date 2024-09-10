package validator

import (
	"api/src/model/request"
	"api/src/model/static"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IScheduleValidator interface {
	// 予定登録
	Create(u *request.CreateSchedule) error
	// 予定更新
	Update(u *request.UpdateSchedule) error
	// 予定削除
	Delete(u *request.DeleteSchedule) error
}

type ScheduleValidator struct{}

func NewScheduleValidator() IScheduleValidator {
	return &ScheduleValidator{}
}

// 予定登録
func (v *ScheduleValidator) Create(u *request.CreateSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.FreqID,
			validation.Min(uint(0)),
			validation.Max(uint(static.FREQ_NONE)),
			IsUintValidator{},
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

// 予定更新
func (v *ScheduleValidator) Update(u *request.UpdateSchedule) error {
	return validation.ValidateStruct(
		u,
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

// 予定削除
func (v *ScheduleValidator) Delete(u *request.DeleteSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}
