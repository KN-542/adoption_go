package enum

type Site int
type ApplicantStatus uint
type CalendarStatus uint

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
	// 毎週
	FREQ_WEEKLY CalendarStatus = 1
	// 毎月
	FREQ_MONTHLY CalendarStatus = 2
	// 毎年
	FREQ_YEARLY CalendarStatus = 3
)
