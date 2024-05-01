package enum

type Protection uint
type HashKeyPre string
type Role uint
type LoginType uint
type Site uint
type ApplicantStatus uint
type CalendarStatus uint

const (
	ON  Protection = 1
	OFF Protection = 0
)

// m_role
const (
	// admin_ロール関連
	ROLE_ADMIN_ROLE_CREATE      Role = 1
	ROLE_ADMIN_ROLE_READ        Role = 2
	ROLE_ADMIN_ROLE_DETAIL_READ Role = 3
	ROLE_ADMIN_ROLE_EDIT        Role = 4
	ROLE_ADMIN_ROLE_DELETE      Role = 5
	ROLE_ADMIN_ROLE_ASSIGN      Role = 6
	// admin_企業関連
	ROLE_ADMIN_COMPANY_CREATE      Role = 101
	ROLE_ADMIN_COMPANY_READ        Role = 102
	ROLE_ADMIN_COMPANY_DETAIL_READ Role = 103
	ROLE_ADMIN_COMPANY_EDIT        Role = 104
	ROLE_ADMIN_COMPANY_DELETE      Role = 105
	// management_ロール関連
	ROLE_MANAGEMENT_ROLE_CREATE      Role = 1001
	ROLE_MANAGEMENT_ROLE_READ        Role = 1002
	ROLE_MANAGEMENT_ROLE_DETAIL_READ Role = 1003
	ROLE_MANAGEMENT_ROLE_EDIT        Role = 1004
	ROLE_MANAGEMENT_ROLE_DELETE      Role = 1005
	ROLE_MANAGEMENT_ROLE_ASSIGN      Role = 1006
	// management_ユーザー関連
	ROLE_MANAGEMENT_USER_CREATE      Role = 2001
	ROLE_MANAGEMENT_USER_READ        Role = 2002
	ROLE_MANAGEMENT_USER_DETAIL_READ Role = 2003
	ROLE_MANAGEMENT_USER_EDIT        Role = 2004
	ROLE_MANAGEMENT_USER_DELETE      Role = 2005
	// management_チーム関連
	ROLE_MANAGEMENT_TEAM_CREATE      Role = 2101
	ROLE_MANAGEMENT_TEAM_READ        Role = 2102
	ROLE_MANAGEMENT_TEAM_DETAIL_READ Role = 2103
	ROLE_MANAGEMENT_TEAM_EDIT        Role = 2104
	ROLE_MANAGEMENT_TEAM_DELETE      Role = 2105
	// management_カレンダー関連
	ROLE_MANAGEMENT_CALENDAR_CREATE      Role = 2201
	ROLE_MANAGEMENT_CALENDAR_READ        Role = 2202
	ROLE_MANAGEMENT_CALENDAR_DETAIL_READ Role = 2203
	ROLE_MANAGEMENT_CALENDAR_EDIT        Role = 2204
	ROLE_MANAGEMENT_CALENDAR_DELETE      Role = 2205
	// management_応募者関連
	ROLE_MANAGEMENT_APPLICANT_CREATE          Role = 2301
	ROLE_MANAGEMENT_APPLICANT_READ            Role = 2302
	ROLE_MANAGEMENT_APPLICANT_DETAIL_READ     Role = 2303
	ROLE_MANAGEMENT_APPLICANT_DOWNLOAD        Role = 2304
	ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL Role = 2305
	ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER     Role = 2306
)

// m_login_type
const (
	LOGIN_TYPE_ADMIN      LoginType = 1
	LOGIN_TYPE_MANAGEMENT LoginType = 2
)

// m_hash_key_pre
const (
	PRE_COMPANY HashKeyPre = "company"
	PRE_ROLE    HashKeyPre = "role"
	PRE_USER    HashKeyPre = "user"
)

// m_site
const (
	// リクナビNEXT
	RECRUIT Site = 1
	// マイナビ
	MYNAVI Site = 2
	// DODA
	DODA Site = 3
	// その他
	OTHER Site = 999
)

// m_applicant_status
const (
	// 日程未回答
	SCHEDULE_UNANSWERED ApplicantStatus = 1
	// 書類未提出
	BOOK_CATEGORY_NOT_PRESENTED ApplicantStatus = 2
	// 1次面接
	INTERVIEW_1 ApplicantStatus = 3
	// 2次面接
	INTERVIEW_2 ApplicantStatus = 4
	// 3次面接
	INTERVIEW_3 ApplicantStatus = 5
	// 4次面接
	INTERVIEW_4 ApplicantStatus = 6
	// 5次面接
	INTERVIEW_5 ApplicantStatus = 7
	// 6次面接
	INTERVIEW_6 ApplicantStatus = 8
	// 7次面接
	INTERVIEW_7 ApplicantStatus = 9
	// 8次面接
	INTERVIEW_8 ApplicantStatus = 10
	// 9次面接
	INTERVIEW_9 ApplicantStatus = 11
	// 10次面接
	INTERVIEW_10 ApplicantStatus = 12
	// 1次面接後課題
	TASK_AFTER_INTERVIEW_1 ApplicantStatus = 13
	// 2次面接後課題
	TASK_AFTER_INTERVIEW_2 ApplicantStatus = 14
	// 3次面接後課題
	TASK_AFTER_INTERVIEW_3 ApplicantStatus = 15
	// 4次面接後課題
	TASK_AFTER_INTERVIEW_4 ApplicantStatus = 16
	// 5次面接後課題
	TASK_AFTER_INTERVIEW_5 ApplicantStatus = 17
	// 6次面接後課題
	TASK_AFTER_INTERVIEW_6 ApplicantStatus = 18
	// 7次面接後課題
	TASK_AFTER_INTERVIEW_7 ApplicantStatus = 19
	// 8次面接後課題
	TASK_AFTER_INTERVIEW_8 ApplicantStatus = 20
	// 9次面接後課題
	TASK_AFTER_INTERVIEW_9 ApplicantStatus = 21
	// 10次面接後課題
	TASK_AFTER_INTERVIEW_10 ApplicantStatus = 22
	// 1次面接落ち
	Failing_TO_PASS_INTERVIEW_1 ApplicantStatus = 23
	// 2次面接落ち
	Failing_TO_PASS_INTERVIEW_2 ApplicantStatus = 24
	// 3次面接落ち
	Failing_TO_PASS_INTERVIEW_3 ApplicantStatus = 25
	// 4次面接落ち
	Failing_TO_PASS_INTERVIEW_4 ApplicantStatus = 26
	// 5次面接落ち
	Failing_TO_PASS_INTERVIEW_5 ApplicantStatus = 27
	// 6次面接落ち
	Failing_TO_PASS_INTERVIEW_6 ApplicantStatus = 28
	// 7次面接落ち
	Failing_TO_PASS_INTERVIEW_7 ApplicantStatus = 29
	// 8次面接落ち
	Failing_TO_PASS_INTERVIEW_8 ApplicantStatus = 30
	// 9次面接落ち
	Failing_TO_PASS_INTERVIEW_9 ApplicantStatus = 31
	// 10次面接落ち
	Failing_TO_PASS_INTERVIEW_10 ApplicantStatus = 32
	// 内定
	OFFER ApplicantStatus = 33
	// 内定承諾
	OFFER_COMMITMENT ApplicantStatus = 34
	// 書類選考落ち
	Failing_TO_PASS_DOCUMENTS ApplicantStatus = 35
	// 選考辞退
	WITHDRAWAL ApplicantStatus = 36
	// 内定辞退
	OFFER_DISMISSAL ApplicantStatus = 37
	// 内定承諾後辞退
	OFFER_COMMITMENT_DISMISSAL ApplicantStatus = 38
)

// m_calendar_freq_status
const (
	// なし
	FREQ_NONE CalendarStatus = 9
	// 毎日
	FREQ_DAILY CalendarStatus = 1
	// 毎週
	FREQ_WEEKLY CalendarStatus = 2
	// 毎月
	FREQ_MONTHLY CalendarStatus = 3
	// 毎年
	FREQ_YEARLY CalendarStatus = 4
)

func ConvertFreqStatus(value CalendarStatus) string {
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
