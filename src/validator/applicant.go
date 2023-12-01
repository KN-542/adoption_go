package validator

import (
	"api/src/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IApplicantValidator interface {
	S3UploadValidator(u *model.FileUpload) error
}

type ApplicantValidator struct{}

func NewApplicantValidator() IApplicantValidator {
	return &ApplicantValidator{}
}

func (v *ApplicantValidator) S3UploadValidator(a *model.FileUpload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
		validation.Field(
			&a.Name,
			validation.Required,
		),
		validation.Field(
			&a.Extension,
			validation.Required,
		),
		validation.Field(
			&a.NamePre,
			validation.Required,
		),
	)
}
