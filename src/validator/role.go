package validator

import (
	"api/src/model/request"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IRoleValidator interface {
	// ロールチェック
	Check(req *request.CheckRole) error
}

type RoleValidator struct{}

func NewRoleValidator() IRoleValidator {
	return &RoleValidator{}
}

// ロールチェック
func (v *RoleValidator) Check(req *request.CheckRole) error {
	return validation.ValidateStruct(
		req,
		validation.Field(
			&req.UserHashKey,
			validation.Required,
		),
	)
}
