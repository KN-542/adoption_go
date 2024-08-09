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
}

// 応募者ユーザー紐づけ
type ApplicantUserAssociation struct {
	ddl.ApplicantUserAssociation
}

// 検索
type SearchApplicant struct {
	ddl.Applicant
	// ステータス
	StatusNameJa string `json:"status_name"`
	// サイト名
	SiteName string `json:"site_name"`
	// 予定用ハッシュキー
	ScheduleHashKey string `json:"schedule_hash_key"`
	// 開始時刻
	Start time.Time `json:"start"`
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
