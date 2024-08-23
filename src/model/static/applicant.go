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
