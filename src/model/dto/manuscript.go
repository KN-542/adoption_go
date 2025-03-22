package dto

import "api/src/model/request"

// 検索
type SearchManuscript struct {
	request.SearchManuscript
	// 企業ID
	CompanyID uint64 `json:"company_id"`
}

// 検索_チーム＆サイト
type SearchManuscriptByTeamAndSite struct {
	// チームID
	TeamID uint64 `json:"team_id"`
	// サイトID
	SiteID uint `json:"site_id"`
}
