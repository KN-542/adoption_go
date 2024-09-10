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
