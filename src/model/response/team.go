package response

import "api/src/model/entity"

// チーム検索
type SearchTeam struct {
	List []entity.SearchTeam `json:"list"`
}

// チーム取得
type GetTeam struct {
	entity.Team
}

// 自チーム取得
type GetOwnTeam struct {
	Team         entity.Team                    `json:"team"`
	Events       []entity.InterviewEventsByTeam `json:"events"`
	AutoRule     entity.TeamAutoAssignRule      `json:"auto_rule"`
	Priority     []entity.TeamAssignPriority    `json:"priority"`
	PerList      []entity.TeamPerInterview      `json:"per_list"`
	PossibleList []entity.TeamAssignPossible    `json:"possible_list"`
}

// チーム検索_同一企業
type SearchTeamByCompany struct {
	List []entity.SearchTeam `json:"list"`
}

// チーム毎ステータスイベント取得
type StatusEventsByTeam struct {
	List []entity.StatusEventsByTeam `json:"list"`
}
