package validator

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/static"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IApplicantValidator interface {
	Search(a *request.ApplicantSearch) error
	Download(a *request.ApplicantDownload) error
	DownloadSub(a *request.ApplicantDownloadSub) error
	GetStatusList(a *request.ApplicantStatusList) error
	HashKeyValidate(a *ddl.Applicant) error
	S3UploadValidator(a *ddl.FileUpload) error
	S3DownloadValidator(a *ddl.FileDownload) error
	InsertDesiredAtValidator(a *ddl.ApplicantDesired) error
}

type ApplicantValidator struct{}

func NewApplicantValidator() IApplicantValidator {
	return &ApplicantValidator{}
}

func (v *ApplicantValidator) Search(a *request.ApplicantSearch) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Sites,
			validation.Each(validation.Required),
			validation.Each(UniqueValidator{}),
		),
		validation.Field(
			&a.ApplicantStatusList,
			validation.Each(validation.Required),
			validation.Each(UniqueValidator{}),
		),
		validation.Field(
			&a.ResumeFlg,
			validation.Min(0),
			validation.Max(uint(static.DOCUMENT_EXIST)),
			IsUintValidator{},
		),
		validation.Field(
			&a.CurriculumVitaeFlg,
			validation.Min(0),
			validation.Max(uint(static.DOCUMENT_EXIST)),
			IsUintValidator{},
		),
		validation.Field(
			&a.Users,
			validation.Each(validation.Required),
			validation.Each(UniqueValidator{}),
		),
	)
}

func (v *ApplicantValidator) Download(a *request.ApplicantDownload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.SiteHashKey,
			validation.Required,
		),
	)
}

func (v *ApplicantValidator) DownloadSub(a *request.ApplicantDownloadSub) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.OuterID,
			validation.Required,
			validation.Length(1, 100),
		),
		validation.Field(
			&a.Name,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(
			&a.Email,
			validation.Required,
			validation.Length(1, 100),
			is.Email,
		),
		validation.Field(
			&a.Tel,
			validation.Required,
			validation.Length(1, 100),
			is.UTFNumeric,
		),
		validation.Field(
			&a.Age,
			validation.Min(12),
			validation.Max(100),
		),
	)
}

func (v *ApplicantValidator) GetStatusList(a *request.ApplicantStatusList) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.UserHashKey,
			validation.Required,
		),
	)
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
