package static

const (
	ON  uint = 1
	OFF uint = 0
)

// m_role
const (
	// admin_ロール関連
	ROLE_ADMIN_ROLE_CREATE      uint = 1
	ROLE_ADMIN_ROLE_READ        uint = 2
	ROLE_ADMIN_ROLE_DETAIL_READ uint = 3
	ROLE_ADMIN_ROLE_EDIT        uint = 4
	ROLE_ADMIN_ROLE_DELETE      uint = 5
	ROLE_ADMIN_ROLE_ASSIGN      uint = 6
	// admin_企業関連
	ROLE_ADMIN_COMPANY_CREATE      uint = 101
	ROLE_ADMIN_COMPANY_READ        uint = 102
	ROLE_ADMIN_COMPANY_DETAIL_READ uint = 103
	ROLE_ADMIN_COMPANY_EDIT        uint = 104
	ROLE_ADMIN_COMPANY_DELETE      uint = 105
	// admin_ユーザー関連
	ROLE_ADMIN_USER_CREATE      uint = 201
	ROLE_ADMIN_USER_READ        uint = 202
	ROLE_ADMIN_USER_DETAIL_READ uint = 203
	ROLE_ADMIN_USER_EDIT        uint = 204
	ROLE_ADMIN_USER_DELETE      uint = 205
	// admin_操作ログ関連
	ROLE_ADMIN_LOG_READ        uint = 301
	ROLE_ADMIN_LOG_DETAIL_READ uint = 302
	// management_ロール関連
	ROLE_MANAGEMENT_ROLE_CREATE      uint = 1001
	ROLE_MANAGEMENT_ROLE_READ        uint = 1002
	ROLE_MANAGEMENT_ROLE_DETAIL_READ uint = 1003
	ROLE_MANAGEMENT_ROLE_EDIT        uint = 1004
	ROLE_MANAGEMENT_ROLE_DELETE      uint = 1005
	ROLE_MANAGEMENT_ROLE_ASSIGN      uint = 1006
	// management_ユーザー関連
	ROLE_MANAGEMENT_USER_CREATE      uint = 1101
	ROLE_MANAGEMENT_USER_READ        uint = 1102
	ROLE_MANAGEMENT_USER_DETAIL_READ uint = 1103
	ROLE_MANAGEMENT_USER_EDIT        uint = 1104
	ROLE_MANAGEMENT_USER_DELETE      uint = 1105
	// management_チーム関連
	ROLE_MANAGEMENT_TEAM_CREATE      uint = 1201
	ROLE_MANAGEMENT_TEAM_READ        uint = 1202
	ROLE_MANAGEMENT_TEAM_DETAIL_READ uint = 1203
	ROLE_MANAGEMENT_TEAM_EDIT        uint = 1204
	ROLE_MANAGEMENT_TEAM_DELETE      uint = 1205
	// management_予定関連
	ROLE_MANAGEMENT_SCHEDULE_CREATE      uint = 1301
	ROLE_MANAGEMENT_SCHEDULE_READ        uint = 1302
	ROLE_MANAGEMENT_SCHEDULE_DETAIL_READ uint = 1303
	ROLE_MANAGEMENT_SCHEDULE_EDIT        uint = 1304
	ROLE_MANAGEMENT_SCHEDULE_DELETE      uint = 1305
	// management_応募者関連
	ROLE_MANAGEMENT_APPLICANT_CREATE             uint = 1401
	ROLE_MANAGEMENT_APPLICANT_READ               uint = 1402
	ROLE_MANAGEMENT_APPLICANT_DETAIL_READ        uint = 1403
	ROLE_MANAGEMENT_APPLICANT_DOWNLOAD           uint = 1404
	ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL    uint = 1405
	ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER        uint = 1406
	ROLE_MANAGEMENT_APPLICANT_SETTING_MANUSCRIPT uint = 1407
	ROLE_MANAGEMENT_APPLICANT_SETTING_TYPE       uint = 1408
	ROLE_MANAGEMENT_APPLICANT_SETTING_STATUS     uint = 1409
	// management_原稿関連
	ROLE_MANAGEMENT_MANUSCRIPT_CREATE      uint = 1501
	ROLE_MANAGEMENT_MANUSCRIPT_READ        uint = 1502
	ROLE_MANAGEMENT_MANUSCRIPT_DETAIL_READ uint = 1503
	ROLE_MANAGEMENT_MANUSCRIPT_EDIT        uint = 1504
	ROLE_MANAGEMENT_MANUSCRIPT_DELETE      uint = 1505
	// management_メール関連
	ROLE_MANAGEMENT_MAIL_CREATE      uint = 1601
	ROLE_MANAGEMENT_MAIL_READ        uint = 1602
	ROLE_MANAGEMENT_MAIL_DETAIL_READ uint = 1603
	ROLE_MANAGEMENT_MAIL_EDIT        uint = 1604
	ROLE_MANAGEMENT_MAIL_DELETE      uint = 1605
	// management_変数関連
	ROLE_MANAGEMENT_VARIABLE_CREATE      uint = 1701
	ROLE_MANAGEMENT_VARIABLE_READ        uint = 1702
	ROLE_MANAGEMENT_VARIABLE_DETAIL_READ uint = 1703
	ROLE_MANAGEMENT_VARIABLE_EDIT        uint = 1704
	ROLE_MANAGEMENT_VARIABLE_DELETE      uint = 1705
	// management_分析関連
	ROLE_MANAGEMENT_ANALYSIS_READ uint = 1801
	// management_操作ログ関連
	ROLE_MANAGEMENT_LOG_READ        uint = 1901
	ROLE_MANAGEMENT_LOG_DETAIL_READ uint = 1902
	// management_設定関連
	ROLE_MANAGEMENT_SETTING_COMPANY uint = 2001
	ROLE_MANAGEMENT_SETTING_TEAM    uint = 2002
)

// m_sidebar
const (
	SIDEBAR_ADMIN_COMPANY         uint = 1
	SIDEBAR_ADMIN_USER            uint = 2
	SIDEBAR_ADMIN_ROLE            uint = 3
	SIDEBAR_ADMIN_LOG             uint = 4
	SIDEBAR_MANAGEMENT_APPLICANT  uint = 101
	SIDEBAR_MANAGEMENT_USER       uint = 102
	SIDEBAR_MANAGEMENT_TEAM       uint = 103
	SIDEBAR_MANAGEMENT_SCHEDULE   uint = 104
	SIDEBAR_MANAGEMENT_ROLE       uint = 105
	SIDEBAR_MANAGEMENT_MANUSCRIPT uint = 106
	SIDEBAR_MANAGEMENT_MAIL       uint = 107
	SIDEBAR_MANAGEMENT_VARIABLE   uint = 108
	SIDEBAR_MANAGEMENT_ANALYSIS   uint = 109
	SIDEBAR_MANAGEMENT_LOG        uint = 110
)

// m_login_type
const (
	LOGIN_TYPE_ADMIN      uint = 1
	LOGIN_TYPE_MANAGEMENT uint = 2
)

// m_select_status_event
const (
	STATUS_EVENT_DECIDE_SCHEDULE           uint = 1
	STATUS_EVENT_SUBMIT_DOCUMENTS          uint = 2
	STATUS_EVENT_SUBMIT_DOCUMENTS_NOT_PASS uint = 3
	STATUS_EVENT_SUBMIT_DOCUMENTS_PASS     uint = 4
	STATUS_EVENT_INTERVIEW_PASS            uint = 5
	STATUS_EVENT_INTERVIEW_FAIL            uint = 6
)

// m_assign_rule
const (
	ASSIGN_RULE_MANUAL uint = 1
	ASSIGN_RULE_AUTO   uint = 2
)
const (
	ASSIGN_RULE_CONFIG_UNREQUIRED uint = 0
	ASSIGN_RULE_CONFIG_REQUIRED   uint = 1
)

// m_auto_assign_rule
const (
	AUTO_ASSIGN_RULE_RANDOM        uint = 1
	AUTO_ASSIGN_RULE_ASC           uint = 2
	AUTO_ASSIGN_RULE_DESC_SCHEDULE uint = 3
)
const (
	AUTO_ASSIGN_RULE_CONFIG_UNREQUIRED uint = 0
	AUTO_ASSIGN_RULE_CONFIG_REQUIRED   uint = 1
)

// m_hash_key_pre
const (
	PRE_COMPANY        string = "company"
	PRE_ROLE           string = "role"
	PRE_USER           string = "user"
	PRE_TEAM           string = "team"
	PRE_SCHEDULE       string = "schedule"
	PRE_SELECT_STATUS  string = "select_status"
	PRE_APPLICANT      string = "applicant"
	PRE_APPLICANT_TYPE string = "applicant_type"
	PRE_MANUSCRIPT     string = "manuscript"
)

// m_site
const (
	// リクナビNEXT
	RECRUIT uint = 1
	// マイナビ
	MYNAVI uint = 2
	// DODA
	DODA uint = 3

	// リクナビNEXT
	FILE_NAME_RECRUIT        string = "oubosha"
	INDEX_RECRUIT_OUTER_ID   uint   = 11
	INDEX_RECRUIT_NAME       uint   = 12
	INDEX_RECRUIT_EMAIL      uint   = 17
	INDEX_RECRUIT_TEL        uint   = 18
	INDEX_RECRUIT_AGE        uint   = 14
	INDEX_RECRUIT_MANUSCRIPT uint   = 3
	COLUMNS_RECRUIT          uint   = 220
	// マイナビ
	FILE_NAME_MYNAVI        string = "mynavi"
	INDEX_MYNAVI_OUTER_ID   uint   = 0
	INDEX_MYNAVI_NAME       uint   = 1 // ※性: 1, 名: 2
	INDEX_MYNAVI_EMAIL      uint   = 9
	INDEX_MYNAVI_TEL        uint   = 11 // 空文字の場合は12の電話番号(自宅)をチェック
	INDEX_MYNAVI_AGE        uint   = 6
	INDEX_MYNAVI_MANUSCRIPT uint   = 18
	COLUMNS_MYNAVI          uint   = 381
	// DODA
	FILE_NAME_DODA        string = "Senko"
	INDEX_DODA_OUTER_ID   uint   = 186
	INDEX_DODA_NAME       uint   = 6 // ※性: 6, 名: 7
	INDEX_DODA_EMAIL      uint   = 13
	INDEX_DODA_TEL        uint   = 18 // 空文字の場合は19の電話番号(自宅)をチェック
	INDEX_DODA_AGE        uint   = 11
	INDEX_DODA_MANUSCRIPT uint   = 3
	COLUMNS_DODA          uint   = 186
)

// m_schedule_freq_status
const (
	// なし
	FREQ_NONE uint = 9
	// 毎日
	FREQ_DAILY uint = 1
	// 毎週
	FREQ_WEEKLY uint = 2
	// 毎月
	FREQ_MONTHLY uint = 3
	// 毎年
	FREQ_YEARLY uint = 4
)

func ConvertFreqStatus(value uint) string {
	switch {
	case value == FREQ_DAILY:
		return "day"
	case value == FREQ_WEEKLY:
		return "week"
	case value == FREQ_MONTHLY:
		return "month"
	case value == FREQ_YEARLY:
		return "year"
	}
	return ""
}

// m_interview_processing
const (
	INTERVIEW_PROCESSING_NOW  uint = 1
	INTERVIEW_PROCESSING_PASS uint = 2
	INTERVIEW_PROCESSING_FAIL uint = 3
)

// m_document_rule
const (
	DOCUMENT_RULE_REPUDIATE        uint = 1
	DOCUMENT_RULE_SUBMIT_CONFIRM   uint = 2
	DOCUMENT_RULE_REQUIRED_CONFIRM uint = 3
)

// m_occupation
const (
	OCCUPATION_ENGINEER             uint = 1
	OCCUPATION_PROJECT_MANAGER      uint = 2
	OCCUPATION_DESIGNER             uint = 3
	OCCUPATION_ACCOUNTANT           uint = 4
	OCCUPATION_HR_MANAGER           uint = 5
	OCCUPATION_CONSULTANT           uint = 6
	OCCUPATION_SALES_REPRESENTATIVE uint = 7
	OCCUPATION_MARKETING_SPECIALIST uint = 8
	OCCUPATION_CUSTOMER_SUPPORT     uint = 9
	OCCUPATION_CEO                  uint = 10
)
