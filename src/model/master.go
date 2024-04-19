package model

/*
m_login_type
ログイン種別マスタ
*/
type LoginType struct {
	AbstractMasterModel
	// ログイン種別
	Type string `json:"type" gorm:"not null;unique;type:varchar(15)"`
	// 遷移パス
	Path string `json:"path" gorm:"not null;unique;type:varchar(15)"`
}

/*
m_site
媒体マスタ
*/
type Site struct {
	AbstractMasterModel
	// 媒体名
	SiteName string `json:"site_name" gorm:"unique;not null;type:varchar(20)"`
}

/*
m_role
ロールマスタ
*/
type Role struct {
	AbstractMasterModel
	// ロール名_日本語
	NameJa string `json:"name_ja" gorm:"unique;not null;type:varchar(20)"`
	// ロール名_英語
	NameEn string `json:"name_en" gorm:"unique;not null;type:varchar(30)"`
}

/*
m_applicant_status
選考状況マスタ
*/
type ApplicantStatus struct {
	AbstractMasterModel
	// ステータス名_日本語
	StatusNameJa string `json:"status_name_ja" gorm:"unique;not null;type:varchar(20)"`
	// ステータス名_英語
	StatusNameEn string `json:"status_name_en" gorm:"unique;not null;type:varchar(30)"`
}

/*
m_calendar_freq_status
予定頻度マスタ
*/
type CalendarFreqStatus struct {
	AbstractMasterModel
	// 頻度
	Freq string `json:"freq" gorm:"unique;not null;type:varchar(10)"`
	// 名前_日本語
	NameJa string `json:"name_ja" gorm:"unique;not null;type:varchar(10)"`
	// 名前_英語
	NameEn string `json:"name_en" gorm:"unique;not null;type:varchar(10)"`
}

/*
m_apply_variable
適用変数種別マスタ
*/
type ApplyVariable struct {
	AbstractMasterModel
	// 種別名
	Name string `json:"name" gorm:"unique;not null;type:varchar(20)"`
}

/*
m_view_role_operation
画面操作ロールマスタ
*/
type ViewRoleOperation struct {
	AbstractMasterModel
	// 操作名
	Name string `json:"name" gorm:"unique;not null;type:varchar(30)"`
	// 操作ロール
	Roles string `json:"roles" gorm:"type:text"`
}

/*
m_operation_log_event
操作ログイベントマスタ
*/
type OperationLogEvent struct {
	AbstractMasterModel
	// イベント内容
	Event string `json:"event" gorm:"unique;not null;type:text"`
}

/*
m_notice
通知マスタ
*/
type NoticeType struct {
	AbstractMasterModel
	// 通知内容
	NoticeType string `json:"notice" gorm:"unique;not null;type:text"`
}

/*
m_analysis_term
分析項目マスタ
*/
type AnalysisTerm struct {
	AbstractMasterModel
	// 項目_日本語
	TermJa string `json:"term_ja" gorm:"unique;not null;type:varchar(20)"`
	// 項目_英語
	TermEn string `json:"term_en" gorm:"unique;not null;type:varchar(30)"`
}

/*
m_hash_key_pre
ハッシュキープレビューマスタ
*/
type HashKeyPre struct {
	AbstractMasterModel
	// プレビュー
	Pre string `json:"pre" gorm:"unique;not null;type:varchar(10)"`
}

/*
m_s3_name_pre
S3ファイル名プレビューマスタ
*/
type S3NamePre struct {
	AbstractMasterModel
	// プレビュー
	Pre string `json:"pre" gorm:"unique;not null;type:varchar(20)"`
}

func (m LoginType) TableName() string {
	return "m_login_type"
}
func (m Site) TableName() string {
	return "m_site"
}
func (m Role) TableName() string {
	return "m_role"
}
func (m ApplicantStatus) TableName() string {
	return "m_applicant_status"
}
func (m CalendarFreqStatus) TableName() string {
	return "m_calendar_freq_status"
}
func (m ApplyVariable) TableName() string {
	return "m_apply_variable"
}
func (m ViewRoleOperation) TableName() string {
	return "m_view_role_operation"
}
func (m OperationLogEvent) TableName() string {
	return "m_operation_log_event"
}
func (m NoticeType) TableName() string {
	return "m_notice"
}
func (m AnalysisTerm) TableName() string {
	return "m_analysis_term"
}
func (m HashKeyPre) TableName() string {
	return "m_hash_key_pre"
}
func (m S3NamePre) TableName() string {
	return "m_s3_name_pre"
}

type Sites struct {
	List []Site `json:"list"`
}
type ApplicantStatusList struct {
	List []ApplicantStatus `json:"list"`
}
