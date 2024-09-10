package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type ITeamRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *ddl.Team) (*entity.Team, error)
	// 検索
	Search(m *ddl.Team) ([]*entity.SearchTeam, error)
	// 取得
	Get(m *ddl.Team) (*entity.Team, error)
	// 取得_PK
	GetByPrimary(m *ddl.Team) (*entity.Team, error)
	// 更新
	Update(tx *gorm.DB, m *ddl.Team) (*entity.Team, error)
	// 削除
	Delete(tx *gorm.DB, m *ddl.Team) error
	// 検索_同一企業
	SearchByCompany(m *dto.SearchTeamByCompany) ([]entity.SearchTeam, error)
	// チーム紐づけ登録
	InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error
	// チーム紐づけ一括登録
	InsertsTeamAssociation(tx *gorm.DB, m []*ddl.TeamAssociation) error
	// チーム紐づけ一覧取得
	ListTeamAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error)
	// 所属チーム一覧取得
	ListBelongTeam(m *ddl.TeamAssociation) ([]entity.Team, error)
	// ユーザー紐づけ一覧取得
	ListUserAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error)
	// チーム紐づけ削除
	DeleteTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error
	// チーム毎ステータスイベント取得
	StatusEventsByTeam(m *ddl.Team) ([]entity.StatusEventsByTeam, error)
	// チーム面接毎イベント取得
	InterviewEventsByTeam(m *ddl.Team) ([]entity.InterviewEventsByTeam, error)
	// 選考状況登録
	InsertSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// 選考状況一括登録
	InsertsSelectStatus(tx *gorm.DB, m []*ddl.SelectStatus) (*entity.ApplicantStatusList, error)
	// 選考状況削除
	DeleteSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error
	// イベント一括登録
	InsertsEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error
	// イベント取得
	SelectEventAssociation(m *ddl.TeamEvent) ([]entity.TeamEvent, error)
	// イベント削除
	DeleteEventAssociation(tx *gorm.DB, m *ddl.TeamEvent) error
	// 面接毎イベント一括登録
	InsertsEventEachInterviewAssociation(tx *gorm.DB, m []*ddl.TeamEventEachInterview) error
	// ～次面接イベント取得
	GetEventEachInterviewAssociation(m *ddl.TeamEventEachInterview) (*entity.TeamEventEachInterview, error)
	// 面接毎イベント削除
	DeleteEventEachInterviewAssociation(tx *gorm.DB, m *ddl.TeamEventEachInterview) error
	// 面接自動割り当てルールイベント登録
	InsertAutoAssignRule(tx *gorm.DB, m *ddl.TeamAutoAssignRule) error
	// 面接自動割り当てルールイベント取得
	GetAutoAssignRule(m *ddl.TeamAutoAssignRule) (*entity.TeamAutoAssignRule, error)
	// 面接自動割り当てルールイベント取得_Find
	GetAutoAssignRuleFind(m *ddl.TeamAutoAssignRule) ([]entity.TeamAutoAssignRule, error)
	// 面接自動割り当てルールイベント削除
	DeleteAutoAssignRule(tx *gorm.DB, m *ddl.TeamAutoAssignRule) error
	// 面接割り振り優先順位一括登録
	InsertsAssignPriority(tx *gorm.DB, m []*ddl.TeamAssignPriority) error
	// 面接割り振り優先順位取得
	GetAssignPriority(m *ddl.TeamAssignPriority) ([]*entity.TeamAssignPriority, error)
	// 面接割り振り優先順位取得_結合なし
	GetAssignPriorityOnly(m *ddl.TeamAssignPriority) ([]*entity.TeamAssignPriorityOnly, error)
	// 面接割り振り優先順位取得_複数チーム
	GetAssignPriorityTeams(m []uint64) ([]*entity.TeamAssignPriority, error)
	// 面接割り振り優先順位削除
	DeleteAssignPriority(tx *gorm.DB, m *ddl.TeamAssignPriority) error
	// 面接毎参加可能者一括登録
	InsertsAssignPossible(tx *gorm.DB, m []*ddl.TeamAssignPossible) error
	// 面接毎参加可能者取得
	GetAssignPossible(m *ddl.TeamAssignPossible) ([]entity.TeamAssignPossible, error)
	// 面接毎参加可能者取得 by 面接回数
	GetAssignPossibleByNumOfInterview(m *ddl.TeamAssignPossible) ([]entity.TeamAssignPossible, error)
	// 面接毎参加可能者予定取得
	GetAssignPossibleSchedule(m *ddl.TeamAssignPossible) ([]entity.AssignPossibleSchedule, error)
	// 面接毎参加可能者削除
	DeleteAssignPossible(tx *gorm.DB, m *ddl.TeamAssignPossible) error
	// 面接毎設定一括登録
	InsertsPerInterview(tx *gorm.DB, m []*ddl.TeamPerInterview) error
	// 面接毎設定取得
	GetPerInterview(m *ddl.TeamPerInterview) ([]entity.TeamPerInterview, error)
	// 面接毎設定取得 by 面接回数
	GetPerInterviewByNumOfInterview(m *ddl.TeamPerInterview) (*entity.TeamPerInterview, error)
	// 面接毎設定削除
	DeletePerInterview(tx *gorm.DB, m *ddl.TeamPerInterview) error
	// チームID取得
	GetIDs(m []string) ([]uint64, error)
	// チーム取得_ハッシュキー配列
	GetByHashKeys(m []string) ([]entity.Team, error)
}

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) ITeamRepository {
	return &TeamRepository{db}
}

// 登録
func (u *TeamRepository) Insert(tx *gorm.DB, m *ddl.Team) (*entity.Team, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.Team{
		Team: *m,
	}, nil
}

// 検索
func (u *TeamRepository) Search(m *ddl.Team) ([]*entity.SearchTeam, error) {
	var l []*entity.SearchTeam

	query := u.db.Table("t_team").
		Select(`
			t_team.id,
			t_team.hash_key,
			t_team.name
		`).
		Where("t_team.company_id = ?", m.CompanyID)

	if err := query.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name")
	}).Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return l, nil
}

// 取得
func (u *TeamRepository) Get(m *ddl.Team) (*entity.Team, error) {
	var res entity.Team
	if err := u.db.Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name, email")
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 取得_PK
func (u *TeamRepository) GetByPrimary(m *ddl.Team) (*entity.Team, error) {
	var res entity.Team
	if err := u.db.Table("t_team").
		Select(`
			t_team.*,
			m_assign_rule.hash_key as rule_hash
		`).
		Joins("LEFT JOIN m_assign_rule ON t_team.rule_id = m_assign_rule.id").
		Where(
			&ddl.Team{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					ID: m.ID,
				},
			},
		).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Table("t_user").Select("id, hash_key, name, email")
	}).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 更新
func (u *TeamRepository) Update(tx *gorm.DB, m *ddl.Team) (*entity.Team, error) {
	team := ddl.Team{
		Name:           m.Name,
		NumOfInterview: m.NumOfInterview,
		RuleID:         m.RuleID,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.Team{}).Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(team).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var updatedTeam ddl.Team
	if err := tx.Where(
		&ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&updatedTeam).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.Team{
		Team: updatedTeam,
	}, nil
}

// 削除
func (u *TeamRepository) Delete(tx *gorm.DB, m *ddl.Team) error {
	if err := tx.Where(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: m.HashKey,
		},
	}).Delete(&ddl.Team{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索_同一企業
func (u *TeamRepository) SearchByCompany(m *dto.SearchTeamByCompany) ([]entity.SearchTeam, error) {
	var l []entity.SearchTeam

	if err := u.db.Model(&entity.SearchTeam{}).
		Select("t_team.hash_key, t_team.name").
		Joins("inner join t_team_association ON t_team_association.team_id = t_team.id").
		Joins("inner join t_user ON t_user.id = t_team_association.user_id").
		Where("t_team.company_id = ?", m.CompanyID).
		Where("t_user.hash_key = ?", m.UserHashKey).
		Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// チーム紐づけ登録
func (u *TeamRepository) InsertTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム紐づけ一括登録
func (u *TeamRepository) InsertsTeamAssociation(tx *gorm.DB, m []*ddl.TeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 所属チーム一覧取得
func (u *TeamRepository) ListBelongTeam(m *ddl.TeamAssociation) ([]entity.Team, error) {
	var res []entity.Team
	if err := u.db.
		Joins("LEFT JOIN t_team_association ON t_team_association.team_id = t_team.id").
		Where("t_team_association.user_id = ?", m.UserID).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("team not exist")
	}

	return res, nil
}

// チーム紐づけ一覧取得
func (u *TeamRepository) ListTeamAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error) {
	var res []entity.TeamAssociation
	if err := u.db.Where(
		&ddl.TeamAssociation{
			UserID: m.UserID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// ユーザー紐づけ一覧取得
func (u *TeamRepository) ListUserAssociation(m *ddl.TeamAssociation) ([]entity.TeamAssociation, error) {
	var res []entity.TeamAssociation
	if err := u.db.Where(
		&ddl.TeamAssociation{
			TeamID: m.TeamID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// チーム紐づけ削除
func (u *TeamRepository) DeleteTeamAssociation(tx *gorm.DB, m *ddl.TeamAssociation) error {
	if err := tx.Where(m).Delete(&ddl.TeamAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チーム毎ステータスイベント取得
func (u *TeamRepository) StatusEventsByTeam(m *ddl.Team) ([]entity.StatusEventsByTeam, error) {
	var res []entity.StatusEventsByTeam

	query := u.db.Table("t_team_event").
		Select(`
			m_select_status_event.hash_key as event_hash_key,
			m_select_status_event.desc_ja as desc_ja,
			m_select_status_event.desc_en as desc_en,
			t_select_status.hash_key as select_status_hash_key,
			t_select_status.status_name as status_name
		`).
		Joins("LEFT JOIN t_select_status ON t_team_event.status_id = t_select_status.id").
		Joins("LEFT JOIN m_select_status_event ON t_team_event.event_id = m_select_status_event.id").
		Where("t_team_event.team_id = ?", m.ID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// チーム面接毎イベント取得
func (u *TeamRepository) InterviewEventsByTeam(m *ddl.Team) ([]entity.InterviewEventsByTeam, error) {
	var res []entity.InterviewEventsByTeam

	query := u.db.Table("t_team_event_each_interview").
		Select(`
			t_team_event_each_interview.num_of_interview as num_of_interview,
			t_select_status.hash_key as select_status_hash_key,
			t_select_status.status_name as status_name
		`).
		Joins("LEFT JOIN t_select_status ON t_team_event_each_interview.status_id = t_select_status.id").
		Where("t_team_event_each_interview.team_id = ?", m.ID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 選考状況登録
func (u *TeamRepository) InsertSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 選考状況一括登録
func (u *TeamRepository) InsertsSelectStatus(tx *gorm.DB, m []*ddl.SelectStatus) (*entity.ApplicantStatusList, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &entity.ApplicantStatusList{
		List: m,
	}, nil
}

// 選考状況削除
func (u *TeamRepository) DeleteSelectStatus(tx *gorm.DB, m *ddl.SelectStatus) error {
	if err := tx.Model(&ddl.SelectStatus{}).Where(&ddl.SelectStatus{
		TeamID: m.TeamID,
	}).Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント登録
func (u *TeamRepository) InsertEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント一括登録
func (u *TeamRepository) InsertsEventAssociation(tx *gorm.DB, m []*ddl.TeamEvent) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// イベント取得
func (u *TeamRepository) SelectEventAssociation(m *ddl.TeamEvent) ([]entity.TeamEvent, error) {
	var res []entity.TeamEvent
	if err := u.db.Where(
		&ddl.TeamEvent{
			TeamID: m.TeamID,
		},
	).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// イベント削除
func (u *TeamRepository) DeleteEventAssociation(tx *gorm.DB, m *ddl.TeamEvent) error {
	if err := tx.Model(&ddl.TeamEvent{}).Where(&ddl.TeamEvent{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamEvent{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎イベント一括登録
func (u *TeamRepository) InsertsEventEachInterviewAssociation(tx *gorm.DB, m []*ddl.TeamEventEachInterview) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// ～次面接イベント取得
func (u *TeamRepository) GetEventEachInterviewAssociation(m *ddl.TeamEventEachInterview) (*entity.TeamEventEachInterview, error) {
	var res entity.TeamEventEachInterview
	if err := u.db.Table("t_team_event_each_interview").Where(m).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 面接毎イベント削除
func (u *TeamRepository) DeleteEventEachInterviewAssociation(tx *gorm.DB, m *ddl.TeamEventEachInterview) error {
	if err := tx.Model(&ddl.TeamEventEachInterview{}).Where(&ddl.TeamEventEachInterview{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamEventEachInterview{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接自動割り当てルールイベント登録
func (u *TeamRepository) InsertAutoAssignRule(tx *gorm.DB, m *ddl.TeamAutoAssignRule) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接自動割り当てルールイベント取得
func (u *TeamRepository) GetAutoAssignRule(m *ddl.TeamAutoAssignRule) (*entity.TeamAutoAssignRule, error) {
	var res *entity.TeamAutoAssignRule

	query := u.db.Table("t_team_auto_assign_rule_association").
		Select(`
			t_team_auto_assign_rule_association.*,
			m_auto_assign_rule.hash_key as hash_key
		`).
		Joins(`
			LEFT JOIN 
				m_auto_assign_rule
			ON
				m_auto_assign_rule.id = t_team_auto_assign_rule_association.rule_id
		`).
		Where(&ddl.TeamAutoAssignRule{
			TeamID: m.TeamID,
		})

	if err := query.First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接自動割り当てルールイベント取得_Find
func (u *TeamRepository) GetAutoAssignRuleFind(m *ddl.TeamAutoAssignRule) ([]entity.TeamAutoAssignRule, error) {
	var res []entity.TeamAutoAssignRule

	query := u.db.Table("t_team_auto_assign_rule_association").
		Select(`
			t_team_auto_assign_rule_association.*,
			m_auto_assign_rule.hash_key as hash_key
		`).
		Joins(`
			LEFT JOIN 
				m_auto_assign_rule
			ON
				m_auto_assign_rule.id = t_team_auto_assign_rule_association.rule_id
		`).
		Where(&ddl.TeamAutoAssignRule{
			TeamID: m.TeamID,
		})

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接自動割り当てルールイベント削除
func (u *TeamRepository) DeleteAutoAssignRule(tx *gorm.DB, m *ddl.TeamAutoAssignRule) error {
	if err := tx.Where(&ddl.TeamAutoAssignRule{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamAutoAssignRule{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接割り振り優先順位一括登録
func (u *TeamRepository) InsertsAssignPriority(tx *gorm.DB, m []*ddl.TeamAssignPriority) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接割り振り優先順位取得
func (u *TeamRepository) GetAssignPriority(m *ddl.TeamAssignPriority) ([]*entity.TeamAssignPriority, error) {
	var res []*entity.TeamAssignPriority

	query := u.db.Table("t_team_assign_priority").
		Select(`
			t_team_assign_priority.priority,
			t_user.hash_key as hash_key,
			t_user.name as name
		`).
		Joins(`LEFT JOIN t_user ON t_team_assign_priority.user_id = t_user.id`).
		Where(&ddl.TeamAssignPriority{
			TeamID: m.TeamID,
		}).
		Order("t_team_assign_priority.priority ASC")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接割り振り優先順位取得_結合なし
func (u *TeamRepository) GetAssignPriorityOnly(m *ddl.TeamAssignPriority) ([]*entity.TeamAssignPriorityOnly, error) {
	var res []*entity.TeamAssignPriorityOnly

	query := u.db.Table("t_team_assign_priority").
		Where(&ddl.TeamAssignPriority{
			TeamID: m.TeamID,
		}).
		Order("t_team_assign_priority.priority ASC")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接割り振り優先順位取得_複数チーム
func (u *TeamRepository) GetAssignPriorityTeams(m []uint64) ([]*entity.TeamAssignPriority, error) {
	var res []*entity.TeamAssignPriority

	query := u.db.Table("t_team_assign_priority").
		Select(`
			t_team_assign_priority.priority,
			t_user.hash_key as hash_key,
			t_user.name as name
		`).
		Joins(`LEFT JOIN t_user ON t_team_assign_priority.user_id = t_user.id`).
		Where("t_team_assign_priority.team_id IN ?", m).
		Order("t_team_assign_priority.priority ASC")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接割り振り優先順位削除
func (u *TeamRepository) DeleteAssignPriority(tx *gorm.DB, m *ddl.TeamAssignPriority) error {
	if err := tx.Where(&ddl.TeamAssignPriority{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamAssignPriority{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎参加可能者一括登録
func (u *TeamRepository) InsertsAssignPossible(tx *gorm.DB, m []*ddl.TeamAssignPossible) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎参加可能者取得
func (u *TeamRepository) GetAssignPossible(m *ddl.TeamAssignPossible) ([]entity.TeamAssignPossible, error) {
	var res []entity.TeamAssignPossible

	query := u.db.Table("t_team_assign_possible").
		Select(`
			t_team_assign_possible.num_of_interview,
			t_user.hash_key as hash_key,
			t_user.name as name,
			t_user.email as email
		`).
		Joins(`LEFT JOIN t_user ON t_team_assign_possible.user_id = t_user.id`).
		Where(&ddl.TeamAssignPossible{
			TeamID: m.TeamID,
		})

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接毎参加可能者取得 by 面接回数
func (u *TeamRepository) GetAssignPossibleByNumOfInterview(m *ddl.TeamAssignPossible) ([]entity.TeamAssignPossible, error) {
	var res []entity.TeamAssignPossible

	query := u.db.Table("t_team_assign_possible").
		Select(`t_user.hash_key as hash_key`).
		Joins(`LEFT JOIN t_user ON t_team_assign_possible.user_id = t_user.id`).
		Where(&ddl.TeamAssignPossible{
			TeamID:         m.TeamID,
			NumOfInterview: m.NumOfInterview,
		})

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接毎参加可能者予定取得
func (u *TeamRepository) GetAssignPossibleSchedule(m *ddl.TeamAssignPossible) ([]entity.AssignPossibleSchedule, error) {
	var res []entity.AssignPossibleSchedule

	query := u.db.Table("t_team_assign_possible").
		Select(`
			t_user.id as user_id,
			t_user.hash_key as user_hash_key
		`).
		Joins(`LEFT JOIN t_user ON t_team_assign_possible.user_id = t_user.id`).
		Where(m)

	if err := query.Preload("Schedules", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Users")
	}).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接毎参加可能者削除
func (u *TeamRepository) DeleteAssignPossible(tx *gorm.DB, m *ddl.TeamAssignPossible) error {
	if err := tx.Where(&ddl.TeamAssignPossible{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamAssignPossible{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎設定一括登録
func (u *TeamRepository) InsertsPerInterview(tx *gorm.DB, m []*ddl.TeamPerInterview) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 面接毎設定取得
func (u *TeamRepository) GetPerInterview(m *ddl.TeamPerInterview) ([]entity.TeamPerInterview, error) {
	var res []entity.TeamPerInterview

	query := u.db.Select("num_of_interview, user_min").Table("t_team_per_interview").Where(&ddl.TeamPerInterview{
		TeamID: m.TeamID,
	})

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 面接毎設定取得 by 面接回数
func (u *TeamRepository) GetPerInterviewByNumOfInterview(m *ddl.TeamPerInterview) (*entity.TeamPerInterview, error) {
	var res entity.TeamPerInterview

	query := u.db.Table("t_team_per_interview").Where(&ddl.TeamPerInterview{
		TeamID:         m.TeamID,
		NumOfInterview: m.NumOfInterview,
	})

	if err := query.First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &res, nil
}

// 面接毎設定削除
func (u *TeamRepository) DeletePerInterview(tx *gorm.DB, m *ddl.TeamPerInterview) error {
	if err := tx.Where(&ddl.TeamPerInterview{
		TeamID: m.TeamID,
	}).Delete(&ddl.TeamPerInterview{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// チームID取得
func (u *TeamRepository) GetIDs(m []string) ([]uint64, error) {
	var res []entity.Team
	if err := u.db.Model(&ddl.Team{}).
		Select("id").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var IDs []uint64
	for _, row := range res {
		IDs = append(IDs, row.ID)
	}

	return IDs, nil
}

// チーム取得_ハッシュキー配列
func (u *TeamRepository) GetByHashKeys(m []string) ([]entity.Team, error) {
	var res []entity.Team
	if err := u.db.Model(&ddl.Team{}).
		Select("id, hash_key, name, num_of_interview").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}
