package validator

import (
	"api/src/model"

	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IApplicantValidator interface {
	S3UploadValidator(a *model.FileUpload) error
	S3DownloadValidator(a *model.FileDownload) error
	InsertDesiredAtValidator(a *model.ApplicantDesired) error
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
			&a.Extension,
			validation.Required,
		),
		validation.Field(
			&a.NamePre,
			validation.Required,
		),
	)
}
func (v *ApplicantValidator) S3DownloadValidator(a *model.FileDownload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

func (v *ApplicantValidator) InsertDesiredAtValidator(a *model.ApplicantDesired) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
		validation.Field(
			&a.DesiredAt,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`)),
		),
	)
}
