package dto

import (
	"api/src/model/ddl"
	"api/src/model/request"
)

type SearchTeamByCompany struct {
	ddl.Team
	request.Abstract
}
