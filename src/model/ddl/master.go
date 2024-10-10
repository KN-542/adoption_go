package ddl

/*
m_login_type
ログイン種別マスタ
*/
type LoginType struct {
	AbstractMasterModel
	// ログイン種別
	Type string `json:"type" gorm:"not null;unique;check:type <> '';type:varchar(15)"`
	// 遷移パス
	Path string `json:"path" gorm:"not null;unique;check:path <> '';type:varchar(15)"`
}

/*
m_site
媒体マスタ
*/
type Site struct {
	AbstractMasterModel
	// 媒体名
	SiteName string `json:"site_name" gorm:"unique;not null;check:site_name <> '';type:varchar(20)"`
	// ファイル名キーワード
	FileName string `json:"file_name" gorm:"not null;unique;check:file_name <> '';type:varchar(255)"`
	// 媒体側ID_index
	OuterIDIndex uint `json:"outer_id_index"`
	// 氏名_index
	NameIndex uint `json:"name_index"`
	// メールアドレス_index
	EmailIndex uint `json:"email_index"`
	// TEL_index
	TELIndex uint `json:"tel_index"`
	// 年齢_index
	AgeIndex uint `json:"age_index"`
	// 原稿_index
	ManuscriptIndex uint `json:"manuscript_index"`
	// 氏名_チェックタイプ
	NameCheckType uint `json:"name_check_type"`
	// カラム数
	NumOfColumn uint `json:"num_of_column"`
}

/*
m_role
ロールマスタ
*/
type Role struct {
	AbstractMasterModel
	// ロール名_日本語
	NameJa string `json:"name_ja" gorm:"unique;not null;check:name_ja <> '';type:varchar(30)"`
	// ロール名_英語
	NameEn string `json:"name_en" gorm:"unique;not null;check:name_en <> '';type:varchar(50)"`
	// ロール種別
	RoleType uint `json:"role_type"`
	// ログイン種別(外部キー)
	LoginType LoginType `gorm:"foreignKey:role_type;references:id"`
}

/*
m_sidebar
サイドバーマスタ
*/
type Sidebar struct {
	AbstractMasterModel
	// 機能名_日本語
	NameJa string `json:"name_ja" gorm:"not null;check:name_ja <> '';type:varchar(30)"`
	// 機能名_英語
	NameEn string `json:"name_en" gorm:"not null;check:name_en <> '';type:varchar(50)"`
	// 遷移パス
	Path string `json:"path" gorm:"unique;not null;check:path <> '';type:varchar(40)"`
	// 機能種別
	FuncType uint `json:"func_type"`
	// ログイン種別(外部キー)
	LoginType LoginType `gorm:"foreignKey:func_type;references:id"`
}

/*
m_sidebar_role_association
サイドバーロール紐づけマスタ
*/
type SidebarRoleAssociation struct {
	// サイドバーID
	SidebarID uint `json:"sidebar_id" gorm:"primaryKey"`
	// 操作可能ロールID
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	// サイドバー(外部キー)
	Sidebar Sidebar `gorm:"foreignKey:sidebar_id;references:id"`
	// ロールマスタ(外部キー)
	Role Role `gorm:"foreignKey:role_id;references:id"`
}

/*
m_schedule_freq_status
予定頻度マスタ
*/
type ScheduleFreqStatus struct {
	AbstractMasterModel
	// 頻度名
	FreqName string `json:"freq_name" gorm:"unique;not null;type:varchar(10)"`
	// 名前_日本語
	NameJa string `json:"name_ja" gorm:"unique;not null;check:name_ja <> '';type:varchar(10)"`
	// 名前_英語
	NameEn string `json:"name_en" gorm:"unique;not null;check:name_en <> '';type:varchar(10)"`
}

/*
m_apply_variable
適用変数種別マスタ
*/
type ApplyVariable struct {
	AbstractMasterModel
	// 種別名
	Name string `json:"name" gorm:"unique;not null;check:name <> '';type:varchar(20)"`
}

/*
m_operation_log_event
操作ログイベントマスタ
*/
type OperationLogEvent struct {
	AbstractMasterModel
	// イベント内容
	Event string `json:"event" gorm:"unique;not null;check:event <> '';type:text"`
}

/*
m_notice
通知マスタ
*/
type NoticeType struct {
	AbstractMasterModel
	// 通知内容
	Notice string `json:"notice" gorm:"unique;not null;check:notice <> '';type:text"`
}

/*
m_analysis_term
分析項目マスタ
*/
type AnalysisTerm struct {
	AbstractMasterModel
	// 項目_日本語
	TermJa string `json:"term_ja" gorm:"unique;not null;check:term_ja <> '';type:varchar(20)"`
	// 項目_英語
	TermEn string `json:"term_en" gorm:"unique;not null;check:term_en <> '';type:varchar(30)"`
}

/*
m_hash_key_pre
ハッシュキープレビューマスタ
*/
type HashKeyPre struct {
	AbstractMasterModel
	// プレビュー
	Pre string `json:"pre" gorm:"unique;not null;check:pre <> '';type:varchar(10)"`
}

/*
m_s3_name_pre
S3ファイル名プレビューマスタ
*/
type S3NamePre struct {
	AbstractMasterModel
	// プレビュー
	Pre string `json:"pre" gorm:"unique;not null;check:pre <> '';type:varchar(20)"`
}

/*
m_select_status_event
応募者ステータスイベントマスタ
*/
type SelectStatusEvent struct {
	AbstractMasterModel
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
}

/*
m_assign_rule
面接アサインルールマスタ
*/
type AssignRule struct {
	AbstractMasterModel
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
	// 追加設定必要性
	AdditionalConfiguration uint `json:"additional_configuration"`
}

/*
m_auto_assign_rule
面接自動割り当てルールマスタ
*/
type AutoAssignRule struct {
	AbstractMasterModel
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
	// 追加設定必要性
	AdditionalConfiguration uint `json:"additional_configuration"`
}

/*
m_document_rule
書類提出ルールマスタ
*/
type DocumentRule struct {
	AbstractMasterModel
	// ルール_日本語
	RuleJa string `json:"rule_ja" gorm:"text"`
	// ルール_英語
	RuleEn string `json:"rule_en" gorm:"text"`
}

/*
m_occupation
職種マスタ
*/
type Occupation struct {
	AbstractMasterModel
	// 職種名_日本語
	NameJa string `json:"name_ja" gorm:"text"`
	// 職種名_英語
	NameEn string `json:"name_en" gorm:"text"`
}

/*
m_interview_processing
面接過程マスタ
*/
type Processing struct {
	AbstractMasterModel
	// 過程
	Processing string `json:"processing" gorm:"text"`
	// 説明_日本語
	DescJa string `json:"desc_ja" gorm:"text"`
	// 説明_英語
	DescEn string `json:"desc_en" gorm:"text"`
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
func (m Sidebar) TableName() string {
	return "m_sidebar"
}
func (m SidebarRoleAssociation) TableName() string {
	return "m_sidebar_role_association"
}
func (m ScheduleFreqStatus) TableName() string {
	return "m_schedule_freq_status"
}
func (m ApplyVariable) TableName() string {
	return "m_apply_variable"
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
func (m SelectStatusEvent) TableName() string {
	return "m_select_status_event"
}
func (m AssignRule) TableName() string {
	return "m_assign_rule"
}
func (m AutoAssignRule) TableName() string {
	return "m_auto_assign_rule"
}
func (m DocumentRule) TableName() string {
	return "m_document_rule"
}
func (m Occupation) TableName() string {
	return "m_occupation"
}
func (m Processing) TableName() string {
	return "m_interview_processing"
}

type Sites struct {
	List []Site `json:"list"`
}
