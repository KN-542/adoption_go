package dto

import (
	"api/src/model/ddl"
	"api/src/model/request"
)

type SearchUser struct {
	request.SearchUser
	// 企業ID
	CompanyID uint64
	// TODO
}
type SearchTeamByCompany struct {
	ddl.Team
	request.Abstract
}
