package static

// 書類状況
const (
	DOCUMENT_EXIST     uint = 1
	DOCUMENT_NOT_EXIST uint = 2
)

// 面接官予定重複フラグ
const (
	DUPLICATION_OUT     uint = 1
	DUPLICATION_WARNING uint = 2
	DUPLICATION_SAFE    uint = 3
)

// 応募者種別最大数
const APPLICANT_TYPE_SIZE int = 25

// 面接官表示
const (
	INTERVIEWER_DISPLAY     uint = 0
	INTERVIEWER_NOT_DISPLAY uint = 1
)

// 書類通過
const (
	DOCUMENT_PROCESS uint = 0
	DOCUMENT_PASS    uint = 1
	DOCUMENT_FAIL    uint = 2
)
