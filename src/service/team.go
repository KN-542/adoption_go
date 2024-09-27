package service

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"net/http"
	"strconv"
)

type ITeamService interface {
	// 登録
	Create(req *request.CreateTeam) *response.Error
	// 更新
	Update(req *request.UpdateTeam) *response.Error
	// 基本情報更新
	UpdateBasic(req *request.UpdateBasicTeam) *response.Error
	// 削除
	Delete(req *request.DeleteTeam) *response.Error
	// 検索
	Search(req *request.SearchTeam) (*response.SearchTeam, *response.Error)
	// 取得
	Get(req *request.GetTeam) (*response.GetTeam, *response.Error)
	// 自チーム取得
	GetOwn(req *request.GetOwnTeam) (*response.GetOwnTeam, *response.Error)
	// 検索_同一企業
	SearchByCompany(req *request.SearchTeamByCompany) (*response.SearchTeamByCompany, *response.Error)
	// 毎ステータスイベント取得
	StatusEvents(req *request.StatusEventsByTeam) (*response.StatusEventsByTeam, *response.Error)
}

type TeamService struct {
	db         repository.IDBRepository
	redis      repository.IRedisRepository
	user       repository.IUserRepository
	team       repository.ITeamRepository
	schedule   repository.IScheduleRepository
	applicant  repository.IApplicantRepository
	role       repository.IRoleRepository
	manuscript repository.IManuscriptRepository
	master     repository.IMasterRepository
	v          validator.ITeamValidator
	outer      repository.IOuterIFRepository
}

func NewTeamService(
	db repository.IDBRepository,
	redis repository.IRedisRepository,
	user repository.IUserRepository,
	team repository.ITeamRepository,
	schedule repository.IScheduleRepository,
	applicant repository.IApplicantRepository,
	role repository.IRoleRepository,
	manuscript repository.IManuscriptRepository,
	master repository.IMasterRepository,
	v validator.ITeamValidator,
	outer repository.IOuterIFRepository,
) ITeamService {
	return &TeamService{db, redis, user, team, schedule, applicant, role, manuscript, master, v, outer}
}

// 検索
func (u *TeamService) Search(req *request.SearchTeam) (*response.SearchTeam, *response.Error) {
	// バリデーション
	if err := u.v.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 企業ID取得
	ctx := context.Background()
	company, companyErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyParseErr := strconv.ParseUint(*company, 10, 64)
	if companyParseErr != nil {
		log.Printf("%v", companyParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.CompanyID = companyID

	// 検索
	teams, err := u.team.Search(&req.Team)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.SearchTeam
	for _, team := range teams {
		team.ID = 0

		var filteredUsers []*ddl.User

		for _, row := range team.Users {
			filteredUsers = append(filteredUsers, &ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: row.HashKey,
				},
				Name: row.Name,
			})
		}
		team.Users = filteredUsers
		res = append(res, *team)
	}

	return &response.SearchTeam{
		List: res,
	}, nil
}

// 取得
func (u *TeamService) Get(req *request.GetTeam) (*response.GetTeam, *response.Error) {
	// バリデーション
	if err := u.v.Get(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 取得
	team, err := u.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	for _, row := range team.Users {
		row.ID = 0
	}

	return &response.GetTeam{
		Team: entity.Team{
			Team: ddl.Team{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: team.HashKey,
				},
				Name: team.Name,
			},
			Users: team.Users,
		},
	}, nil
}

// 自チーム取得
func (u *TeamService) GetOwn(req *request.GetOwnTeam) (*response.GetOwnTeam, *response.Error) {
	// ID取得
	ctx := context.Background()
	team, teamErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*team, 10, 64)
	if teamIDErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 取得
	res, resErr := u.team.GetByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if resErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	for _, row := range res.Users {
		row.ID = 0
	}

	// 面接毎イベント取得
	events, eventsErr := u.team.InterviewEventsByTeam(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if eventsErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 面接自動割り当てルール取得
	var autoRule entity.TeamAutoAssignRule
	if res.RuleID == static.ASSIGN_RULE_AUTO {
		temp, tempErr := u.team.GetAutoAssignRule(&ddl.TeamAutoAssignRule{
			TeamID: teamID,
		})
		if tempErr != nil {
			return nil, &response.Error{
				Status: http.StatusBadRequest,
			}
		}
		autoRule = *temp
	}

	// 面接割り振り優先順位取得
	var priority []entity.TeamAssignPriority
	if autoRule.RuleID == static.AUTO_ASSIGN_RULE_ASC {
		tempList, tempErr := u.team.GetAssignPriority(&ddl.TeamAssignPriority{
			TeamID: teamID,
		})
		if tempErr != nil {
			return nil, &response.Error{
				Status: http.StatusBadRequest,
			}
		}
		for _, row := range tempList {
			priority = append(priority, *row)
		}
	}

	// 面接毎参加可能者取得
	perList, perListErr := u.team.GetPerInterview(&ddl.TeamPerInterview{
		TeamID: teamID,
	})
	if perListErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 面接毎参加可能者取得
	possibleList, possibleErr := u.team.GetAssignPossible(&ddl.TeamAssignPossible{
		TeamID: teamID,
	})
	if possibleErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	return &response.GetOwnTeam{
		Team: entity.Team{
			Team: ddl.Team{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: res.HashKey,
				},
				Name:           res.Name,
				NumOfInterview: res.NumOfInterview,
			},
			RuleHash: res.RuleHash,
			Users:    res.Users,
		},
		Events: events,
		AutoRule: entity.TeamAutoAssignRule{
			HashKey: autoRule.HashKey,
		},
		Priority:     priority,
		PerList:      perList,
		PossibleList: possibleList,
	}, nil
}

// 登録
func (u *TeamService) Create(req *request.CreateTeam) *response.Error {
	// バリデーション
	if err := u.v.Create(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 企業ID取得
	ctx := context.Background()
	company, companyErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyParseErr := strconv.ParseUint(*company, 10, 64)
	if companyParseErr != nil {
		log.Printf("%v", companyParseErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.CompanyID = companyID

	// ユーザー存在確認
	ids, idsErr := u.user.GetIDs(req.Users)
	if idsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, hashErr := GenerateHash(1, 25)
	if hashErr != nil {
		log.Printf("%v", hashErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	req.HashKey = string(static.PRE_TEAM) + "_" + *hashKey
	req.NumOfInterview = 3
	req.RuleID = static.ASSIGN_RULE_MANUAL
	team, err := u.team.Insert(tx, &req.Team)
	if err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var teamAssociations []*ddl.TeamAssociation
	for _, id := range ids {
		teamAssociations = append(teamAssociations, &ddl.TeamAssociation{
			TeamID: team.ID,
			UserID: id,
		})
	}
	// 紐づけ一括登録
	if err := u.team.InsertsTeamAssociation(tx, teamAssociations); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接毎設定登録 & 面接毎参加可能者登録(全員参加可能)
	var perList []*ddl.TeamPerInterview
	var possibleList []*ddl.TeamAssignPossible
	for i := 1; i <= int(team.NumOfInterview); i++ {
		perList = append(perList, &ddl.TeamPerInterview{
			TeamID:         team.ID,
			NumOfInterview: uint(i),
			UserMin:        1,
		})
		for _, id := range ids {
			possibleList = append(possibleList, &ddl.TeamAssignPossible{
				TeamID:         team.ID,
				UserID:         id,
				NumOfInterview: uint(i),
			})
		}
	}
	if err := u.team.InsertsPerInterview(tx, perList); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if err := u.team.InsertsAssignPossible(tx, possibleList); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	_, selectHash, selectHashErr := GenerateHash(1, 25)
	if selectHashErr != nil {
		log.Printf("%v", selectHashErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	_, selectHash2, selectHash2Err := GenerateHash(1, 25)
	if selectHash2Err != nil {
		log.Printf("%v", selectHash2Err)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 選考状況登録
	if err := u.team.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash,
			CompanyID: companyID,
		},
		TeamID:     team.ID,
		StatusName: "日程未回答",
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if err := u.team.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash2,
			CompanyID: companyID,
		},
		TeamID:     team.ID,
		StatusName: "日程回答済み",
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 更新
func (u *TeamService) Update(req *request.UpdateTeam) *response.Error {
	// バリデーション
	if err := u.v.Update(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 企業ID取得
	ctx := context.Background()
	company, companyErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyParseErr := strconv.ParseUint(*company, 10, 64)
	if companyParseErr != nil {
		log.Printf("%v", companyParseErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.CompanyID = companyID

	// 取得
	team, teamErr := u.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー存在確認
	ids, idsErr := u.user.GetIDs(req.Users)
	if idsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 優先順位取得
	priority, priorityErr := u.team.GetAssignPriority(&ddl.TeamAssignPriority{
		TeamID: team.ID,
	})
	if priorityErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 更新
	_, err := u.team.Update(tx, &req.Team)
	if err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if len(ids) > 0 {
		var teamAssociations []*ddl.TeamAssociation
		var priorityAssociations []*ddl.TeamAssignPriority

		for index, id := range ids {
			teamAssociations = append(teamAssociations, &ddl.TeamAssociation{
				TeamID: team.ID,
				UserID: id,
			})

			if len(priority) > 0 {
				priorityAssociations = append(priorityAssociations, &ddl.TeamAssignPriority{
					TeamID:   team.ID,
					UserID:   id,
					Priority: uint(len(priority) + index + 1),
				})
			}
		}

		// 紐づけ一括登録
		if err := u.team.InsertsTeamAssociation(tx, teamAssociations); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// 優先順位登録
		if len(priorityAssociations) > 0 {
			if err := u.team.InsertsAssignPriority(tx, priorityAssociations); err != nil {
				if err := u.db.TxRollback(tx); err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
		}
	}

	if err := u.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 基本情報更新
func (u *TeamService) UpdateBasic(req *request.UpdateBasicTeam) *response.Error {
	// バリデーション
	if err := u.v.UpdateBasic(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ID取得
	ctx := context.Background()
	t, teamErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*t, 10, 64)
	if teamIDErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 取得
	team, teamErr := u.team.GetByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}
	req.HashKey = team.HashKey

	// ユーザー取得
	users, usersErr := u.team.ListUserAssociation(&ddl.TeamAssociation{
		TeamID: teamID,
	})
	if usersErr != nil {
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 更新
	updateTeam, updateErr := u.team.Update(tx, &req.Team)
	if updateErr != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if team.NumOfInterview < updateTeam.NumOfInterview {
		// 面接毎設定登録 & 面接毎参加可能者登録(全員参加可能)
		var perList []*ddl.TeamPerInterview
		var possibleList []*ddl.TeamAssignPossible
		for i := team.NumOfInterview + 1; i <= updateTeam.NumOfInterview; i++ {
			perList = append(perList, &ddl.TeamPerInterview{
				TeamID:         team.ID,
				NumOfInterview: uint(i),
				UserMin:        1,
			})
			for _, user := range users {
				possibleList = append(possibleList, &ddl.TeamAssignPossible{
					TeamID:         team.ID,
					UserID:         user.UserID,
					NumOfInterview: uint(i),
				})
			}
		}

		if err := u.team.InsertsPerInterview(tx, perList); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		if err := u.team.InsertsAssignPossible(tx, possibleList); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	} else if team.NumOfInterview > updateTeam.NumOfInterview {
		// 面接毎設定 & 面接毎参加可能者削除
		if err := u.team.DeletePerInterviewByNum(tx, &ddl.TeamPerInterview{
			TeamID:         team.ID,
			NumOfInterview: updateTeam.NumOfInterview,
		}); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		if err := u.team.DeleteAssignPossibleByNum(tx, &ddl.TeamAssignPossible{
			TeamID:         team.ID,
			NumOfInterview: updateTeam.NumOfInterview,
		}); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 削除
func (u *TeamService) Delete(req *request.DeleteTeam) *response.Error {
	// バリデーション
	if err := u.v.Delete(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 取得
	team, teamErr := u.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// redisのID取得、一致なら削除不可
	ctx := context.Background()
	t, teamErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*t, 10, 64)
	if teamIDErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if team.ID == teamID {
		return &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_TEAM_USER_CANNOT_DELETE_TEAM,
		}
	}

	// 削除可能判定
	// t_applicant
	apps, appsErr := u.applicant.GetByTeamID(&ddl.Applicant{
		TeamID: team.ID,
	})
	if appsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(apps) > 0 {
		return &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_TEAM_USER_CANNOT_DELETE_APPLICANT,
		}
	}
	// t_schedule
	schedules, schedulesErr := u.schedule.GetByTeamID(&ddl.Schedule{
		TeamID: team.ID,
	})
	if schedulesErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(schedules) > 0 {
		return &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_TEAM_USER_CANNOT_DELETE_SCHEDULE,
		}
	}
	// t_manuscript_team_association
	manuscripts, manuscriptsErr := u.manuscript.GetAssociationByTeamID(
		&ddl.ManuscriptTeamAssociation{
			TeamID: team.ID,
		},
	)
	if manuscriptsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(manuscripts) > 0 {
		return &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_TEAM_USER_CANNOT_DELETE_MANUSCRIPT,
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 関連する紐づけ削除
	// t_team_association
	if err := u.team.DeleteTeamAssociation(tx, &ddl.TeamAssociation{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	//  t_team_per_interview
	if err := u.team.DeletePerInterview(tx, &ddl.TeamPerInterview{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	//  t_team_assign_possible
	if err := u.team.DeleteAssignPossible(tx, &ddl.TeamAssignPossible{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	//  t_team_assign_priority
	if err := u.team.DeleteAssignPriority(tx, &ddl.TeamAssignPriority{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// t_select_status
	if err := u.team.DeleteSelectStatus(tx, &ddl.SelectStatus{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// t_select_status
	if err := u.team.DeleteSelectStatus(tx, &ddl.SelectStatus{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// t_team_auto_assign_rule_association
	if err := u.team.DeleteAutoAssignRule(tx, &ddl.TeamAutoAssignRule{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// t_team_event_each_interview
	if err := u.team.DeleteEventEachInterviewAssociation(tx, &ddl.TeamEventEachInterview{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// t_team_event
	if err := u.team.DeleteEventAssociation(tx, &ddl.TeamEvent{
		TeamID: team.ID,
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := u.team.Delete(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	}); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 検索_同一企業
func (u *TeamService) SearchByCompany(req *request.SearchTeamByCompany) (*response.SearchTeamByCompany, *response.Error) {
	// 企業ID取得
	ctx := context.Background()
	company, companyErr := u.redis.Get(ctx, req.HashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyParseErr := strconv.ParseUint(*company, 10, 64)
	if companyParseErr != nil {
		log.Printf("%v", companyParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 検索
	team, err := u.team.SearchByCompany(&dto.SearchTeamByCompany{
		Team: ddl.Team{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: companyID,
			},
		},
		Abstract: request.Abstract{
			UserHashKey: req.HashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchTeamByCompany{
		List: team,
	}, nil
}

// 毎ステータスイベント取得
func (u *TeamService) StatusEvents(req *request.StatusEventsByTeam) (*response.StatusEventsByTeam, *response.Error) {
	// 取得
	ctx := context.Background()
	teamRedis, teamRedisErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamRedisErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*teamRedis, 10, 64)
	if teamIDErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 取得
	res, err := u.team.StatusEventsByTeam(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.StatusEventsByTeam{
		List: res,
	}, nil
}
