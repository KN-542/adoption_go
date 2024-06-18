package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ICompanyValidator interface {
	// 登録
	Create(c *request.CreateCompany) error
}

type CompanyValidator struct{}

func NewCompanyValidator() ICompanyValidator {
	return &CompanyValidator{}
}

// 登録
func (v *CompanyValidator) Create(c *request.CreateCompany) error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.Email,
			validation.Required,
			validation.Length(1, 50),
			is.Email,
		),
		validation.Field(
			&c.UserHashKey,
			validation.Required,
		),
		validation.Field(
			&c.Name,
			validation.Required,
			validation.Length(1, 30),
		),
	)
}
