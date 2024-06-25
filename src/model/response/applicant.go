package response

import (
	"api/src/model/dto"
	"api/src/model/entity"
	"time"
)

// 検索
type SearchApplicant struct {
	List []entity.SearchApplicant `json:"list"`
}

// サイト一覧取得
type ApplicantSites struct {
	List []entity.Site `json:"list"`
}

// 応募者ステータス一覧取得
type ApplicantStatusList struct {
	List []entity.ApplicantStatus `json:"list"`
}

// 応募者ダウンロード
type ApplicantDownload struct {
	UpdateNum int `json:"update_num"`
}

// 予約表
type ReserveTable struct {
	// 年月日
	Dates []time.Time `json:"date"`
	// 予約時間
	Options []dto.ReserveTableSub `json:"options"`
	// 予定
	Schedule time.Time `json:"schedule"`
	// 予定ハッシュキー
	ScheduleHashKey string `json:"schedule_hash_key"`
}

// 応募者取得
type GetApplicant struct {
	Applicant entity.Applicant `json:"applicant"`
}

// 認証URL作成
type GetOauthURL struct {
	Url string `json:"url"`
}

// GoogleMeetUrl発行
type GetGoogleMeetUrl struct {
	entity.Applicant
}
