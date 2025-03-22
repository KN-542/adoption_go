package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IManuscriptValidator interface {
	// 検索
	Search(m *request.SearchManuscript) error
	// 登録
	Create(m *request.CreateManuscript) error
	// 応募者紐づけ登録
	CreateApplicantAssociation(m *request.CreateApplicantAssociation) error
	// 取得
	Get(m *request.GetManuscript) error
	// 更新
	Update(m *request.UpdateManuscript) error
}

type ManuscriptValidator struct{}

func NewManuscriptValidator() IManuscriptValidator {
	return &ManuscriptValidator{}
}

// 検索
func (v *ManuscriptValidator) Search(m *request.SearchManuscript) error {
	return validation.ValidateStruct(
		m,
	)
}

// 取得
func (v *ManuscriptValidator) Get(m *request.GetManuscript) error {
	return validation.ValidateStruct(
		m,
	)
}

// 登録
func (v *ManuscriptValidator) Create(m *request.CreateManuscript) error {
	return validation.ValidateStruct(
		m,
		validation.Field(
			&m.Content,
			validation.Required,
		),
		validation.Field(
			&m.Teams,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&m.Sites,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
	)
}

// 応募者紐づけ登録
func (v *ManuscriptValidator) CreateApplicantAssociation(m *request.CreateApplicantAssociation) error {
	return validation.ValidateStruct(
		m,
		validation.Field(
			&m.ManuscriptHash,
			validation.Required,
		),
		validation.Field(
			&m.Applicants,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
			UniqueValidator{},
		),
	)
}

// 更新
func (v *ManuscriptValidator) Update(m *request.UpdateManuscript) error {
	return validation.ValidateStruct(
		m,
		validation.Field(
			&m.HashKey,
			validation.Required,
		),
		validation.Field(
			&m.Content,
			validation.Required,
		),
		validation.Field(
			&m.Teams,
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&m.Sites,
			validation.Each(validation.Required),
			UniqueValidator{},
		),
	)
}
