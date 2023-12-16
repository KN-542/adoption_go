package enum

type CellIndex int
type IsAuthCodeExist int
type ItemOfSite int
type isDocumentExist uint

const (
	CELL_CREATED_AT     CellIndex = 0
	CELL_NAME           CellIndex = 1
	CELL_SEX            CellIndex = 2
	CELL_AGE            CellIndex = 3
	CELL_FROM           CellIndex = 4
	CELL_TELL           CellIndex = 5
	CELL_EMAIL          CellIndex = 6
	CELL_INTERVIEW_DATE CellIndex = 7
	CELL_PR             CellIndex = 11
)

const (
	AUTH_CODE_EXIST     IsAuthCodeExist = 1
	AUTH_CODE_NOT_EXIST IsAuthCodeExist = 0
)

const (
	/*
		リクルート
	*/
	RECRUIT_ID    ItemOfSite = 11
	RECRUIT_NAME  ItemOfSite = 12
	RECRUIT_EMAIL ItemOfSite = 17
	RECRUIT_TEL   ItemOfSite = 18
	RECRUIT_AGE   ItemOfSite = 14
)

const (
	DOCUMENT_EXIST     isDocumentExist = 1
	DOCUMENT_NOT_EXIST isDocumentExist = 2
)
