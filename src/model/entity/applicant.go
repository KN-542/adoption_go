package entity

import (
	"api/src/model/ddl"
	"time"
)

// 応募者
type Applicant struct {
	ddl.Applicant
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
	// カレンダー用ハッシュキー
	CalendarHashKey string `json:"calendar_hash_key"`
	// 開始時刻
	Start time.Time `json:"start"`
	// 担当面接官(hash_key)
	Users string `json:"users"`
	// 担当面接官(氏名)
	UserNames string `json:"user_names"`
}

// 応募者ステータス
type ApplicantStatus struct {
	ddl.SelectStatus
}
