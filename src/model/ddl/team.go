package ddl

/*
t_team
チーム
*/
type Team struct {
	AbstractTransactionModel
	// チーム名
	Name string `json:"name" gorm:"not null;check:name <> '';type:varchar(30);index"`
	// 最大面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"check:num_of_interview >= 1 AND num_of_interview <= 30"`
	// ルールID
	RuleID uint `json:"rule_id"`
	// ルール(外部キー)
	Rule AssignRule `gorm:"foreignKey:rule_id;references:id"`
}

/*
t_team_event
チームイベント
*/
type TeamEvent struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// イベントID
	EventID uint `json:"event_id" gorm:"primaryKey"`
	// ステータスID
	StatusID uint64 `json:"status_id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// イベント(外部キー)
	Event SelectStatusEvent `gorm:"foreignKey:event_id;references:id"`
	// ステータス(外部キー)
	Status SelectStatus `gorm:"foreignKey:status_id;references:id"`
}

/*
t_team_event_each_interview
チーム面接毎イベント
*/
type TeamEventEachInterview struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// 面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"primaryKey"`
	// ステータスID
	StatusID uint64 `json:"status_id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ステータス(外部キー)
	Status SelectStatus `gorm:"foreignKey:status_id;references:id"`
}

/*
t_team_association
チーム紐づけ
*/
type TeamAssociation struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_team_auto_assign_rule_association
チーム面接自動割り当てルール紐づけ
*/
type TeamAutoAssignRule struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// ルールID
	RuleID uint `json:"rule_id"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ルール(外部キー)
	Rule AutoAssignRule `gorm:"foreignKey:rule_id;references:id"`
}

/*
t_team_assign_priority
面接割り振り優先順位
*/
type TeamAssignPriority struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// 優先順位
	Priority uint `json:"priority"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_team_per_interview
面接毎設定
*/
type TeamPerInterview struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// 面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"primaryKey"`
	// 最低人数
	UserMin uint `json:"user_min" gorm:"check:user_min >= 1 AND user_min <= 6"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
}

/*
t_team_assign_possible
面接毎参加可能者
*/
type TeamAssignPossible struct {
	// チームID
	TeamID uint64 `json:"team_id" gorm:"primaryKey"`
	// 面接回数
	NumOfInterview uint `json:"num_of_interview" gorm:"primaryKey"`
	// ユーザーID
	UserID uint64 `json:"user_id" gorm:"primaryKey"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_select_status
選考状況
*/
type SelectStatus struct {
	AbstractTransactionModel
	// チームID
	TeamID uint64 `json:"team_id"`
	// ステータス名
	StatusName string `json:"status_name" gorm:"not null;type:varchar(50)"`
	// チーム(外部キー)
	Team Team `gorm:"foreignKey:team_id;references:id"`
}

func (t Team) TableName() string {
	return "t_team"
}
func (t TeamEvent) TableName() string {
	return "t_team_event"
}
func (t TeamEventEachInterview) TableName() string {
	return "t_team_event_each_interview"
}
func (t TeamAssociation) TableName() string {
	return "t_team_association"
}
func (t TeamAutoAssignRule) TableName() string {
	return "t_team_auto_assign_rule_association"
}
func (t TeamAssignPriority) TableName() string {
	return "t_team_assign_priority"
}
func (t TeamPerInterview) TableName() string {
	return "t_team_per_interview"
}
func (t TeamAssignPossible) TableName() string {
	return "t_team_assign_possible"
}
func (t SelectStatus) TableName() string {
	return "t_select_status"
}
