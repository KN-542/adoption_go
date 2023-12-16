package validator

import (
	"api/src/model"
	"api/src/model/enum"
	"errors"
	"fmt"

	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IApplicantValidator interface {
	SearchValidator(a *model.ApplicantSearchRequest) error
	S3UploadValidator(a *model.FileUpload) error
	S3DownloadValidator(a *model.FileDownload) error
	InsertDesiredAtValidator(a *model.ApplicantDesired) error
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

func (v *ApplicantValidator) SearchValidator(a *model.ApplicantSearchRequest) error {
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
