package response

import (
	"api/src/model/dto"
	"api/src/model/entity"
	"time"
)

// 検索
type SearchApplicant struct {
	List []entity.SearchApplicant `json:"list"`
	// 総数
	Num int64 `json:"num"`
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
	// 履歴書表示
	IsResume bool `json:"is_resume"`
	// 職務経歴書表示
	IsCurriculumVitae bool `json:"is_curriculum_vitae"`
}

// 応募者取得
type GetApplicant struct {
	Applicant entity.Applicant `json:"applicant"`
}

// 認証URL作成
type GetOauthURL struct {
	// 認証URL
	AuthURL string `json:"auth_url"`
	// Google Meet URL
	GoogleMeetURL string `json:"google_meet_url"`
}

// GoogleMeetUrl発行
type GetGoogleMeetUrl struct {
	Url string `json:"url"`
}

// 面接官割り振り可能判定
type CheckAssignableUser struct {
	List []CheckAssignableUserSub `json:"list"`
}
type CheckAssignableUserSub struct {
	// ユーザー
	User entity.User `json:"user"`
	// 面接官予定重複フラグ
	DuplFlg uint `json:"dupl_flg"`
}

// 応募者種別一覧
type ListApplicantType struct {
	List []entity.ApplicantType `json:"list"`
}
