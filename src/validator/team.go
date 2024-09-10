package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ITeamValidator interface {
	// 検索
	Search(u *request.SearchTeam) error
	// 登録
	Create(u *request.CreateTeam) error
	// 更新
	Update(u *request.UpdateTeam) error
	// 基本情報更新
	UpdateBasic(u *request.UpdateBasicTeam) error
	// 削除
	Delete(u *request.DeleteTeam) error
	// 取得
	Get(u *request.GetTeam) error
	// 面接官割り振り方法更新
	UpdateAssignMethod(u *request.UpdateAssignMethod) error
	// 面接官割り振り方法更新2
	UpdateAssignMethod2(u *request.UpdateAssignMethod) error
	// 面接官割り振り方法更新3
	UpdateAssignMethod3(u *request.UpdateAssignMethod) error
	// 面接官割り振り方法更新4
	UpdateAssignMethod4(u *request.UpdateAssignMethodSub) error
}

type TeamValidator struct{}

func NewTeamValidator() ITeamValidator {
	return &TeamValidator{}
}

// 検索
func (v *TeamValidator) Search(u *request.SearchTeam) error {
	// TODO
	return validation.ValidateStruct(
		u,
	)
}

// 登録
func (v *TeamValidator) Create(u *request.CreateTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Name,
			validation.Required,
			validation.Length(1, 30*3),
		),
	)
}

// 更新
func (v *TeamValidator) Update(u *request.UpdateTeam) error {
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

// 基本情報更新
func (v *TeamValidator) UpdateBasic(u *request.UpdateBasicTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.NumOfInterview,
			validation.Min(uint(1)),
			validation.Max(uint(30)),
		),
	)
}

// 削除
func (v *TeamValidator) Delete(u *request.DeleteTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

// 取得
func (v *TeamValidator) Get(u *request.GetTeam) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKey,
			validation.Required,
		),
	)
}

// 面接官割り振り方法更新
func (v *TeamValidator) UpdateAssignMethod(u *request.UpdateAssignMethod) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.RuleHash,
			validation.Required,
		),
		validation.Field(
			&u.PossibleList,
			validation.Required,
			validation.Length(1, 0),
		),
	)
}

// 面接官割り振り方法更新2
func (v *TeamValidator) UpdateAssignMethod2(u *request.UpdateAssignMethod) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.AutoRuleHash,
			validation.Required,
		),
	)
}

// 面接官割り振り方法更新3
func (v *TeamValidator) UpdateAssignMethod3(u *request.UpdateAssignMethod) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Priority,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&u.PossibleList,
			validation.Required,
			validation.Length(1, 0),
		),
	)
}

// 面接官割り振り方法更新4
func (v *TeamValidator) UpdateAssignMethod4(u *request.UpdateAssignMethodSub) error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.HashKeys,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&u.NumOfInterview,
			validation.Min(uint(1)),
		),
	)
}
