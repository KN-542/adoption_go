package entity

import (
	"api/src/model/ddl"
	"time"
)

// 応募者
type Applicant struct {
	ddl.Applicant
	// 予定ID
	ScheduleID uint64 `json:"schedule_id"`
	// 面接開始日
	Start time.Time `json:"start"`
	// 履歴書拡張子
	ResumeExtension string `json:"resume_extension"`
	// 職務経歴書拡張子
	CurriculumVitaeExtension string `json:"curriculum_vitae_extension"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url"`
	// 原稿ハッシュ
	ManuscriptHash string `json:"manuscript_hash"`
}

// 応募者種別
type ApplicantType struct {
	ddl.ApplicantType
	// 書類提出ルール_日本語
	RuleJa string `json:"rule_ja"`
	// 書類提出ルール_英語
	RuleEn string `json:"rule_en"`
	// 職種名_日本語
	NameJa string `json:"name_ja"`
	// 職種名_英語
	NameEn string `json:"name_en"`
	// 書類確認必要性
	IsDocumentConfirm bool `json:"is_document_confirm"`
}

// 応募者ユーザー紐づけ
type ApplicantUserAssociation struct {
	ddl.ApplicantUserAssociation
}

// 検索
type SearchApplicant struct {
	ddl.Applicant
	// ステータス
	StatusName string `json:"status_name"`
	// サイト名
	SiteName string `json:"site_name"`
	// 過程ハッシュ
	ProcessHash string `json:"process_hash"`
	// 予定用ハッシュキー
	ScheduleHashKey string `json:"schedule_hash_key"`
	// 開始時刻
	Start time.Time `json:"start"`
	// 履歴書拡張子
	ResumeExtension string `json:"resume_extension"`
	// 職務経歴書拡張子
	CurriculumVitaeExtension string `json:"curriculum_vitae_extension"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url"`
	// 原稿内容
	Content string `json:"content"`
	// 種別
	Type string `json:"type"`
	// 担当面接官
	Users []*ddl.User `json:"users" gorm:"many2many:t_applicant_user_association;foreignKey:id;joinForeignKey:applicant_id;References:id;joinReferences:user_id"`
}

// 応募者ステータス
type ApplicantStatus struct {
	ddl.SelectStatus
}

// 応募者ステータス
type ApplicantStatusList struct {
	List []*ddl.SelectStatus
}

// Google Meet URL
type ApplicantURLAssociation struct {
	ddl.ApplicantURLAssociation
}
