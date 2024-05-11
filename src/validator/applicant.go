package validator

import (
	"api/src/model/ddl"
	"api/src/model/enum"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IApplicantValidator interface {
	HashKeyValidate(a *ddl.Applicant) error
	SearchValidator(a *ddl.ApplicantSearchRequest) error
	S3UploadValidator(a *ddl.FileUpload) error
	S3DownloadValidator(a *ddl.FileDownload) error
	InsertDesiredAtValidator(a *ddl.ApplicantDesired) error
}

type ApplicantValidator struct{}

func NewApplicantValidator() IApplicantValidator {
	return &ApplicantValidator{}
}

// uint型の検証
func validateUintRange(min, max uint) validation.RuleFunc {
	return func(value interface{}) error {
		u, ok := value.(uint)
		if !ok {
			return errors.New("invalid type to uint")
		}
		if u < min || u > max {
			return fmt.Errorf("must be between %d and %d", min, max)
		}
		return nil
	}
}

func (v *ApplicantValidator) HashKeyValidate(a *ddl.Applicant) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

func (v *ApplicantValidator) SearchValidator(a *ddl.ApplicantSearchRequest) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Resume,
			validation.By(validateUintRange(0, uint(enum.DOCUMENT_NOT_EXIST))),
		),
		validation.Field(
			&a.CurriculumVitae,
			validation.By(validateUintRange(0, uint(enum.DOCUMENT_NOT_EXIST))),
		),
	)
}

func (v *ApplicantValidator) S3UploadValidator(a *ddl.FileUpload) error {
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
func (v *ApplicantValidator) S3DownloadValidator(a *ddl.FileDownload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

func (v *ApplicantValidator) InsertDesiredAtValidator(a *ddl.ApplicantDesired) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
		validation.Field(
			&a.DesiredAt,
			validation.Required,
		),
	)
}
