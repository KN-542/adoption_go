package enum

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
	// management_ロール関連
	ROLE_MANAGEMENT_ROLE_CREATE      uint = 1001
	ROLE_MANAGEMENT_ROLE_READ        uint = 1002
	ROLE_MANAGEMENT_ROLE_DETAIL_READ uint = 1003
	ROLE_MANAGEMENT_ROLE_EDIT        uint = 1004
	ROLE_MANAGEMENT_ROLE_DELETE      uint = 1005
	ROLE_MANAGEMENT_ROLE_ASSIGN      uint = 1006
	// management_ユーザー関連
	ROLE_MANAGEMENT_USER_CREATE      uint = 2001
	ROLE_MANAGEMENT_USER_READ        uint = 2002
	ROLE_MANAGEMENT_USER_DETAIL_READ uint = 2003
	ROLE_MANAGEMENT_USER_EDIT        uint = 2004
	ROLE_MANAGEMENT_USER_DELETE      uint = 2005
	// management_チーム関連
	ROLE_MANAGEMENT_TEAM_CREATE      uint = 2101
	ROLE_MANAGEMENT_TEAM_READ        uint = 2102
	ROLE_MANAGEMENT_TEAM_DETAIL_READ uint = 2103
	ROLE_MANAGEMENT_TEAM_EDIT        uint = 2104
	ROLE_MANAGEMENT_TEAM_DELETE      uint = 2105
	// management_カレンダー関連
	ROLE_MANAGEMENT_CALENDAR_CREATE      uint = 2201
	ROLE_MANAGEMENT_CALENDAR_READ        uint = 2202
	ROLE_MANAGEMENT_CALENDAR_DETAIL_READ uint = 2203
	ROLE_MANAGEMENT_CALENDAR_EDIT        uint = 2204
	ROLE_MANAGEMENT_CALENDAR_DELETE      uint = 2205
	// management_応募者関連
	ROLE_MANAGEMENT_APPLICANT_CREATE          uint = 2301
	ROLE_MANAGEMENT_APPLICANT_READ            uint = 2302
	ROLE_MANAGEMENT_APPLICANT_DETAIL_READ     uint = 2303
	ROLE_MANAGEMENT_APPLICANT_DOWNLOAD        uint = 2304
	ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL uint = 2305
	ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER     uint = 2306
)

// m_sidebar
const (
	SIDEBAR_ADMIN_COMPANY        uint = 1
	SIDEBAR_ADMIN_USER           uint = 2
	SIDEBAR_ADMIN_ROLE           uint = 3
	SIDEBAR_ADMIN_LOG            uint = 4
	SIDEBAR_MANAGEMENT_APPLICANT uint = 101
	SIDEBAR_MANAGEMENT_USER      uint = 102
	SIDEBAR_MANAGEMENT_ROLE      uint = 103
	SIDEBAR_MANAGEMENT_MAIL      uint = 104
	SIDEBAR_MANAGEMENT_ANALYSIS  uint = 105
	SIDEBAR_MANAGEMENT_LOG       uint = 106
)

// m_login_type
const (
	LOGIN_TYPE_ADMIN      uint = 1
	LOGIN_TYPE_MANAGEMENT uint = 2
)

// m_hash_key_pre
const (
	PRE_COMPANY string = "company"
	PRE_ROLE    string = "role"
	PRE_USER    string = "user"
	PRE_TEAM    string = "team"
)

// m_site
const (
	// リクナビNEXT
	RECRUIT uint = 1
	// マイナビ
	MYNAVI uint = 2
	// DODA
	DODA uint = 3
	// その他
	OTHER uint = 999
)

// m_applicant_status
const (
	// 日程未回答
	SCHEDULE_UNANSWERED uint = 1
	// 書類未提出
	BOOK_CATEGORY_NOT_PRESENTED uint = 2
	// 1次面接
	INTERVIEW_1 uint = 3
	// 2次面接
	INTERVIEW_2 uint = 4
	// 3次面接
	INTERVIEW_3 uint = 5
	// 4次面接
	INTERVIEW_4 uint = 6
	// 5次面接
	INTERVIEW_5 uint = 7
	// 6次面接
	INTERVIEW_6 uint = 8
	// 7次面接
	INTERVIEW_7 uint = 9
	// 8次面接
	INTERVIEW_8 uint = 10
	// 9次面接
	INTERVIEW_9 uint = 11
	// 10次面接
	INTERVIEW_10 uint = 12
	// 1次面接後課題
	TASK_AFTER_INTERVIEW_1 uint = 13
	// 2次面接後課題
	TASK_AFTER_INTERVIEW_2 uint = 14
	// 3次面接後課題
	TASK_AFTER_INTERVIEW_3 uint = 15
	// 4次面接後課題
	TASK_AFTER_INTERVIEW_4 uint = 16
	// 5次面接後課題
	TASK_AFTER_INTERVIEW_5 uint = 17
	// 6次面接後課題
	TASK_AFTER_INTERVIEW_6 uint = 18
	// 7次面接後課題
	TASK_AFTER_INTERVIEW_7 uint = 19
	// 8次面接後課題
	TASK_AFTER_INTERVIEW_8 uint = 20
	// 9次面接後課題
	TASK_AFTER_INTERVIEW_9 uint = 21
	// 10次面接後課題
	TASK_AFTER_INTERVIEW_10 uint = 22
	// 1次面接落ち
	Failing_TO_PASS_INTERVIEW_1 uint = 23
	// 2次面接落ち
	Failing_TO_PASS_INTERVIEW_2 uint = 24
	// 3次面接落ち
	Failing_TO_PASS_INTERVIEW_3 uint = 25
	// 4次面接落ち
	Failing_TO_PASS_INTERVIEW_4 uint = 26
	// 5次面接落ち
	Failing_TO_PASS_INTERVIEW_5 uint = 27
	// 6次面接落ち
	Failing_TO_PASS_INTERVIEW_6 uint = 28
	// 7次面接落ち
	Failing_TO_PASS_INTERVIEW_7 uint = 29
	// 8次面接落ち
	Failing_TO_PASS_INTERVIEW_8 uint = 30
	// 9次面接落ち
	Failing_TO_PASS_INTERVIEW_9 uint = 31
	// 10次面接落ち
	Failing_TO_PASS_INTERVIEW_10 uint = 32
	// 内定
	OFFER uint = 33
	// 内定承諾
	OFFER_COMMITMENT uint = 34
	// 書類選考落ち
	Failing_TO_PASS_DOCUMENTS uint = 35
	// 選考辞退
	WITHDRAWAL uint = 36
	// 内定辞退
	OFFER_DISMISSAL uint = 37
	// 内定承諾後辞退
	OFFER_COMMITMENT_DISMISSAL uint = 38
)

// m_calendar_freq_status
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
