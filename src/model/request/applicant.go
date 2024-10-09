package request

import (
	"api/src/model/ddl"
	"time"
)

// 検索
type SearchApplicant struct {
	Abstract
	ddl.Applicant
	// ページ
	Page int `json:"page"`
	// ページサイズ
	PageSize int `json:"page_size"`
	// サイト一覧
	Sites []string `json:"sites"`
	// 応募者ステータス
	ApplicantStatusList []string `json:"applicant_status_list"`
	// 原稿一覧
	Manuscripts []string `json:"manuscripts"`
	// 種別一覧
	Types []string `json:"types"`
	// 履歴書フラグ
	ResumeFlg uint `json:"resume_flg"`
	// 職務経歴書フラグ
	CurriculumVitaeFlg uint `json:"curriculum_vitae_flg"`
	// 面接予定日_From
	InterviewerDateFrom time.Time `json:"interviewer_date_from"`
	// 面接予定日_To
	InterviewerDateTo time.Time `json:"interviewer_date_to"`
	// 登録日時_From
	CreatedAtFrom time.Time `json:"created_at_from"`
	// 登録日時_To
	CreatedAtTo time.Time `json:"created_at_to"`
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
	// リフレッシュトークン
	RefreshToken string `json:"refresh_token"`
	// 認証コード
	Code string `json:"code"`
}

// 面接希望日登録
type InsertDesiredAt struct {
	ddl.Schedule
	// 応募者ハッシュキー
	ApplicantHashKey string `json:"applicant_hash_key"`
	// 希望面接日時
	DesiredAt time.Time `json:"desired_at"`
	// タイトル
	Title string `json:"title"`
	// 履歴書拡張子
	ResumeExtension string `json:"resume_extension"`
	// 職務経歴書拡張子
	CurriculumVitaeExtension string `json:"curriculum_vitae_extension"`
}

// 応募者ステータス変更
type UpdateStatus struct {
	Abstract
	ddl.SelectStatus
	// ステータス
	Status []string `json:"status"`
	// ステータス紐づけ
	Association []UpdateStatusSub `json:"association"`
	// イベント
	Events []UpdateStatusSub2 `json:"events"`
	// 面接毎イベント
	EventsOfInterview []UpdateStatusSub3 `json:"events_of_interview"`
}

// 応募者ステータス変更サブ
type UpdateStatusSub struct {
	// 変更前ハッシュ
	BeforeHash string `json:"before_hash"`
	// 変更後インデックス
	AfterIndex int `json:"after_index"`
}

// 応募者ステータス変更サブ2
type UpdateStatusSub2 struct {
	// イベントマスタID
	EventID uint `json:"event_id"`
	// イベントマスタハッシュ
	EventHash string `json:"event_hash"`
	// ステータス
	Status int `json:"status"`
}

// 応募者ステータス変更サブ3
type UpdateStatusSub3 struct {
	// 面接回数
	Num uint `json:"num"`
	// 過程
	ProcessHash string `json:"process_hash"`
	// ステータス
	Status int `json:"status"`
}

// 面接官割り振り
type AssignUser struct {
	Abstract
	ddl.Applicant
	// ハッシュキーリスト
	HashKeys []string `json:"hash_keys"`
}

// 面接官割り振り可能判定
type CheckAssignableUser struct {
	Abstract
	// 開始時刻
	Start time.Time `json:"start"`
	// ハッシュキーリスト
	HashKeys []string `json:"hash_keys"`
	// 除外予定ハッシュリスト
	RemoveScheduleHashKeys []string `json:"remove_schedule_hash_keys"`
}

// 種別登録
type CreateApplicantType struct {
	Abstract
	// 種別名
	Name string `json:"name"`
	// 書類提出ルールハッシュ
	RuleHash string `json:"rule_hash"`
	// 職種ハッシュ
	OccupationHash string `json:"occupation_hash"`
}

// 応募者種別紐づけ登録
type CreateApplicantTypeAssociation struct {
	Abstract
	// ハッシュキー
	TypeHash string `json:"type_hash"`
	// 応募者
	Applicants []string `json:"applicants"`
}

// 種別一覧
type ListApplicantType struct {
	Abstract
}

// ステータス更新
type UpdateSelectStatus struct {
	Abstract
	// ハッシュキー
	StatusHash string `json:"status_hash"`
	// 応募者
	Applicants []string `json:"applicants"`
}
