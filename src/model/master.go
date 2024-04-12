package model

/*
	m_company
	企業マスタ
*/
type Company struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 企業名
	Name string `json:"name" gorm:"not null;type:varchar(30)"`
	// ロゴファイル名
	Logo string `json:"logo" gorm:"not null;type:varchar(30)"`
}

/*
	m_site
	媒体マスタ
*/
type Site struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 媒体名_日本語
	SiteNameJa string `json:"site_name_ja" gorm:"type:varchar(20)"`
	// 媒体名_英語
	SiteNameEn string `json:"site_name_en" gorm:"type:varchar(30)"`
}

/*
	m_role
	ロールマスタ
*/
type Role struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
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
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
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
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
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
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 種別名
	Name string `json:"name" gorm:"unique;not null;type:varchar(20)"`
}

/*
	m_view_role_operation
	画面操作ロールマスタ
*/
type ViewRoleOperation struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
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
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// イベント内容
	Event string `json:"event" gorm:"unique;not null;type:text"`
}

/*
	m_notice
	通知マスタ
*/
type Notice struct {
	// ID
	ID uint `json:"id" gorm:"primaryKey"`
	// 通知内容
	Notice string `json:"notice" gorm:"unique;not null;type:text"`
}

func (m Company) TableName() string {
	return "m_company"
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
func (m Notice) TableName() string {
	return "m_notice"
}

type Sites struct {
	List []Site `json:"list"`
}
type ApplicantStatusList struct {
	List []ApplicantStatus `json:"list"`
}
