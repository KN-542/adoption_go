package request

import (
	"api/src/model/ddl"
	"time"
)

// 検索
type SearchApplicant struct {
	Abstract
	ddl.Applicant
	// サイト一覧
	Sites []string `json:"sites"`
	// 応募者ステータス
	ApplicantStatusList []string `json:"applicant_status_list"`
	// 履歴書フラグ
	ResumeFlg uint `json:"resume_flg"`
	// 職務経歴書フラグ
	CurriculumVitaeFlg uint `json:"curriculum_vitae_flg"`
	// 面接官
	Users []string `json:"users"`
	// ソート(key)
	SortKey string `json:"sort_key"`
	// ソート(向き)
	SortAsc bool `json:"sort_asc"`
}

// 応募者ステータス一覧取得
type ApplicantStatusList struct {
	Abstract
}

// 応募者ダウンロード sub request
type ApplicantDownloadSub struct {
	// 媒体側ID
	OuterID string `json:"outer_id"`
	// 氏名
	Name string `json:"name"`
	// メールアドレス
	Email string `json:"email"`
	// TEL
	Tel string `json:"tel"`
	// 年齢
	Age int64 `json:"age"`
}

// 応募者ダウンロード
type ApplicantDownload struct {
	Abstract
	// サイトハッシュキー
	SiteHashKey string `json:"site_hash_key"`
	// 応募者
	Applicants []ApplicantDownloadSub `json:"applicants"`
}

// 予約表表示
type ReserveTable struct {
	ddl.Applicant
}

// ファイルアップロード
type FileUpload struct {
	ddl.Applicant
	// ファイル拡張子
	Extension string `json:"extension"`
	// ファイル名(Pre)
	NamePre string `json:"name_pre"`
	// ファイル名
	Name string `json:"name"`
}

// ファイルダウンロード
type FileDownload struct {
	Abstract
	ddl.Applicant
	// ファイル名(Pre)
	NamePre string `json:"name_pre"`
}

// 取得
type GetApplicant struct {
	Abstract
	ddl.Applicant
}

// 認証URL作成
type GetOauthURL struct {
	Abstract
	ddl.Applicant
}

// GoogleMeetUrl発行
type GetGoogleMeetUrl struct {
	Abstract
	ddl.Applicant
	Code string `json:"code"`
}

// 面接希望日登録
type InsertDesiredAt struct {
	// ハッシュキー
	HashKey string `json:"hash_key"`
	// 希望面接日時
	DesiredAt time.Time `json:"desired_at"`
	// タイトル
	Title string `json:"title"`
	// カレンダーハッシュキー
	CalendarHashKey string `json:"calendar_hash_key"`
}
