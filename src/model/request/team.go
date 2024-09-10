package request

import "api/src/model/ddl"

// チーム登録
type CreateTeam struct {
	Abstract
	ddl.Team
	// ユーザーリスト
	Users []string `json:"users"`
}

// チーム更新
type UpdateTeam struct {
	Abstract
	ddl.Team
	// ユーザーリスト
	Users []string `json:"users"`
}

// チーム基本情報更新
type UpdateBasicTeam struct {
	Abstract
	ddl.Team
}

// チーム削除
type DeleteTeam struct {
	Abstract
	ddl.Team
}

// チーム取得
type GetTeam struct {
	Abstract
	ddl.Team
}

// 自チーム取得
type GetOwnTeam struct {
	Abstract
}

// チーム存在確認
type ConfirmTeamApplicant struct {
	ddl.Team
}

// チーム検索
type SearchTeam struct {
	Abstract
	ddl.Team
}

// チーム検索_同一企業
type SearchTeamByCompany struct {
	ddl.User
}

// チーム毎ステータスイベント取得
type StatusEventsByTeam struct {
	Abstract
}

// 面接官割り振り方法更新
type UpdateAssignMethod struct {
	Abstract
	// 最低面接人数
	UserMin uint `json:"user_min"`
	// ルールハッシュ
	RuleHash string `json:"rule_hash"`
	// 自動ルールハッシュ
	AutoRuleHash string `json:"auto_rule_hash"`
	// 優先順位
	Priority []string `json:"priority"`
	// 各面接参加可能者
	PossibleList []UpdateAssignMethodSub `json:"possible_list"`
}
type UpdateAssignMethodSub struct {
	ddl.TeamPerInterview
	// ハッシュキーリスト
	HashKeys []string `json:"hash_keys"`
}
