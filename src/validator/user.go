package validator

import (
	"api/src/model/request"
	"api/src/model/static"

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
	// チーム検索
	SearchTeam(u *request.SearchTeam) error
	// チーム登録
	CreateTeam(u *request.CreateTeam) error
	// チーム更新
	UpdateTeam(u *request.UpdateTeam) error
	// チーム基本情報更新
	UpdateBasicTeam(u *request.UpdateBasicTeam) error
	// チーム削除
	DeleteTeam(u *request.DeleteTeam) error
	// チーム取得
	GetTeam(u *request.GetTeam) error
	// 予定登録
	CreateSchedule(u *request.CreateSchedule) error
	// 予定更新
	UpdateSchedule(u *request.UpdateSchedule) error
	// 予定検索
	SearchSchedule(u *request.SearchSchedule) error
	// 予定削除
	DeleteSchedule(u *request.DeleteSchedule) error
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

// チーム検索
func (v *UserValidator) SearchTeam(u *request.SearchTeam) error {
	// TODO
	return validation.ValidateStruct(
		u,
	)
}

// チーム登録
func (v *UserValidator) CreateTeam(u *request.CreateTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30*3),
		),
	)
}

// チーム更新
func (v *UserValidator) UpdateTeam(u *request.UpdateTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30*3),
		),
	)
}

// チーム基本情報更新
func (v *UserValidator) UpdateBasicTeam(u *request.UpdateBasicTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.NumOfInterview,
			validation.Min(uint(1)),
			validation.Max(uint(30)),
		),
	)
}

// チーム削除
func (v *UserValidator) DeleteTeam(u *request.DeleteTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

// チーム取得
func (v *UserValidator) GetTeam(u *request.GetTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

func (v *UserValidator) CreateSchedule(u *request.CreateSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.FreqID,
			validation.Min(0),
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
func (v *UserValidator) UpdateSchedule(u *request.UpdateSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.FreqID,
			validation.Min(0),
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

// 予定検索
func (v *UserValidator) SearchSchedule(u *request.SearchSchedule) error {
	return validation.ValidateStruct(
		u,
	)
}

// 予定削除
func (v *UserValidator) DeleteSchedule(u *request.DeleteSchedule) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}
