package request

import "api/src/model/ddl"

// 検索
type ApplicantSearch struct {
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
type ApplicantDownloadSubRequest struct {
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
type ApplicantDownloadRequest struct {
	Abstract
	// サイトハッシュキー
	SiteHashKey string `json:"site_hash_key"`
	// 応募者
	Applicants []ApplicantDownloadSubRequest `json:"applicants"`
}
