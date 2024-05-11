package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ICommonValidator interface {
	// サイドバー表示
	Sidebar(c *request.Sidebar) error
	// 使用可能ロール一覧
	Roles(c *request.Roles) error
}

type CommonValidator struct{}

func NewCommonValidator() ICommonValidator {
	return &CommonValidator{}
}

// サイドバー表示
func (v *CommonValidator) Sidebar(c *request.Sidebar) error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.HashKey,
			validation.Required,
		),
	)
}

// 使用可能ロール一覧
func (v *CommonValidator) Roles(c *request.Roles) error {
	return validation.ValidateStruct(
		c,
		validation.Field(
			&c.HashKey,
			validation.Required,
		),
	)
}
