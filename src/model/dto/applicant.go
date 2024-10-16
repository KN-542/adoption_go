package dto

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"time"
)

// 検索
type SearchApplicant struct {
	request.SearchApplicant
	// ユーザー
	Users []string
}

// 予約表サブ
type ReserveTableSub struct {
	// 時間
	Time time.Time `json:"time"`
	// 予約可否
	IsReserve bool `json:"is_reserve"`
}

// ダウンロード時重複チェック
type CheckDuplDownloading struct {
	// チームID
	TeamID uint64
	// 企業ID
	CompanyID uint64
	// 媒体側IDリスト
	List []string
}

// 応募者原稿紐づけ
type ApplicantManuscriptAssociation struct {
	ddl.Applicant
	// 原稿ID
	ManuscriptID uint64 `json:"manuscript_id"`
	// 原稿ハッシュ
	ManuscriptHash string `json:"manuscript_hash"`
}
