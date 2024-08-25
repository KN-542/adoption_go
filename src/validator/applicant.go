package validator

import (
	"api/src/model/request"
	"api/src/model/static"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IApplicantValidator interface {
	// 検索
	Search(a *request.SearchApplicant) error
	// 応募者ダウンロード
	Download(a *request.ApplicantDownload) error
	// 応募者ダウンロード_サブ構造体
	DownloadSub(a *request.ApplicantDownloadSub) error
	// 応募者ステータス一覧取得
	GetStatusList(a *request.ApplicantStatusList) error
	// 予約表表示
	ReserveTable(a *request.ReserveTable) error
	// 書類アップロード(S3)
	S3Upload(a *request.FileUpload) error
	// 書類ダウンロード(S3)
	S3Download(a *request.FileDownload) error
	// 取得
	Get(a *request.GetApplicant) error
	// 認証URL作成
	GetOauthURL(a *request.GetOauthURL) error
	// GetGoogleMeetUrl発行
	GetGoogleMeetUrl(a *request.GetGoogleMeetUrl) error
	// 面接希望日登録
	InsertDesiredAt(a *request.InsertDesiredAt) error
	// 応募者ステータス変更
	UpdateStatus(a *request.UpdateStatus) error
	// 応募者ステータス変更サブ
	UpdateStatusSub(a *request.UpdateStatusSub) error
	// 応募者ステータス変更サブ2
	UpdateStatusSub2(a *request.UpdateStatusSub2) error
	// 面接官割り振り
	AssignUser(a *request.AssignUser) error
	// 面接官割り振り可能判定
	CheckAssignableUser(a *request.CheckAssignableUser) error
}

type ApplicantValidator struct{}

func NewApplicantValidator() IApplicantValidator {
	return &ApplicantValidator{}
}

// 検索
func (v *ApplicantValidator) Search(a *request.SearchApplicant) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Sites,
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&a.ApplicantStatusList,
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&a.ResumeFlg,
			MinUintValidator{Min: 0},
			MaxUintValidator{Max: static.DOCUMENT_NOT_EXIST},
			IsUintValidator{},
		),
		validation.Field(
			&a.CurriculumVitaeFlg,
			MinUintValidator{Min: 0},
			MaxUintValidator{Max: static.DOCUMENT_NOT_EXIST},
			IsUintValidator{},
		),
		validation.Field(
			&a.Users,
			validation.Each(validation.Required),
			UniqueValidator{},
		),
		validation.Field(
			&a.InterviewerDateFrom,
			validation.By(func(value interface{}) error {
				return IsBeforeTime(a.InterviewerDateFrom, a.InterviewerDateTo)
			}),
		),
		validation.Field(
			&a.CreatedAtFrom,
			validation.By(func(value interface{}) error {
				return IsBeforeTime(a.CreatedAtFrom, a.CreatedAtTo)
			}),
		),
	)
}

// 応募者ダウンロード
func (v *ApplicantValidator) Download(a *request.ApplicantDownload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.SiteHashKey,
			validation.Required,
		),
	)
}

// 応募者ダウンロード_サブ構造体
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

// 応募者ステータス一覧取得
func (v *ApplicantValidator) GetStatusList(a *request.ApplicantStatusList) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.UserHashKey,
			validation.Required,
		),
	)
}

// 予約表表示
func (v *ApplicantValidator) ReserveTable(a *request.ReserveTable) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

// 書類アップロード(S3)
func (v *ApplicantValidator) S3Upload(a *request.FileUpload) error {
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

// 書類ダウンロード(S3)
func (v *ApplicantValidator) S3Download(a *request.FileDownload) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

// 取得
func (v *ApplicantValidator) Get(a *request.GetApplicant) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

// 認証URL作成
func (v *ApplicantValidator) GetOauthURL(a *request.GetOauthURL) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
	)
}

// GetGoogleMeetUrl発行
func (v *ApplicantValidator) GetGoogleMeetUrl(a *request.GetGoogleMeetUrl) error {
	return validation.ValidateStruct(
		a,
	)
}

// 面接希望日登録
func (v *ApplicantValidator) InsertDesiredAt(a *request.InsertDesiredAt) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Title,
			validation.Required,
			validation.Length(1, 75),
		),
	)
}

// 応募者ステータス変更
func (v *ApplicantValidator) UpdateStatus(a *request.UpdateStatus) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Status,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
		),
		validation.Field(
			&a.Association,
			validation.Required,
		),
		validation.Field(
			&a.Events,
		),
		validation.Field(
			&a.EventsOfInterview,
		),
	)
}

// 応募者ステータス変更サブ
func (v *ApplicantValidator) UpdateStatusSub(a *request.UpdateStatusSub) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.BeforeHash,
			validation.Required,
		),
	)
}

// 応募者ステータス変更サブ2
func (v *ApplicantValidator) UpdateStatusSub2(a *request.UpdateStatusSub2) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.EventHash,
			validation.Required,
		),
		validation.Field(
			&a.Status,
			validation.Min(1),
		),
	)
}

// 面接官割り振り
func (v *ApplicantValidator) AssignUser(a *request.AssignUser) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.HashKey,
			validation.Required,
		),
		validation.Field(
			&a.HashKeys,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
		),
	)
}

// 面接官割り振り可能判定
func (v *ApplicantValidator) CheckAssignableUser(a *request.CheckAssignableUser) error {
	return validation.ValidateStruct(
		a,
		validation.Field(
			&a.Start,
			validation.Required,
		),
		validation.Field(
			&a.HashKeys,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
		),
	)
}
