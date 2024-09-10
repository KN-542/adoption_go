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
	"sort"
	"strconv"
)

type IUserService interface {
	// 登録
	Create(req *request.CreateUser) (*response.CreateUser, *response.Error)
	// 検索
	Search(req *request.SearchUser) (*response.SearchUser, *response.Error)
	// 検索_同一企業
	SearchByCompany(req *request.SearchUserByCompany) (*response.SearchUserByCompany, *response.Error)
	// 取得
	Get(req *request.GetUser) (*response.GetUser, *response.Error)
	// ステータスイベントマスタ一覧
	ListStatusEvent() (*response.ListStatusEvent, *response.Error)
	// アサイン関連マスタ取得
	AssignMaster() (*response.AssignMaster, *response.Error)
	// 面接官割り振り方法更新
	UpdateAssignMethod(req *request.UpdateAssignMethod) *response.Error
	// 書類提出ルールマスタ取得
	DocumentRuleMaster() (*response.DocumentRule, *response.Error)
	// 職種マスタ取得
	OccupationMaster() (*response.Occupation, *response.Error)
}

type UserService struct {
	user          repository.IUserRepository
	team          repository.ITeamRepository
	schedule      repository.IScheduleRepository
	role          repository.IRoleRepository
	applicant     repository.IApplicantRepository
	manuscript    repository.IManuscriptRepository
	master        repository.IMasterRepository
	validator     validator.IUserValidator
	validatorTeam validator.ITeamValidator
	db            repository.IDBRepository
	outer         repository.IOuterIFRepository
	redis         repository.IRedisRepository
}

func NewUserService(
	user repository.IUserRepository,
	team repository.ITeamRepository,
	schedule repository.IScheduleRepository,
	role repository.IRoleRepository,
	applicant repository.IApplicantRepository,
	manuscript repository.IManuscriptRepository,
	master repository.IMasterRepository,
	validator validator.IUserValidator,
	validatorTeam validator.ITeamValidator,
	db repository.IDBRepository,
	outer repository.IOuterIFRepository,
	redis repository.IRedisRepository,
) IUserService {
	return &UserService{user, team, schedule, role, applicant, manuscript, master, validator, validatorTeam, db, outer, redis}
}

// 登録
func (u *UserService) Create(req *request.CreateUser) (*response.CreateUser, *response.Error) {
	// バリデーション
	if err := u.validator.Create(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ログイン種別、企業ID取得
	ctx := context.Background()
	login, loginTypeErr := u.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if loginTypeErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	loginType, loginTypeParseErr := strconv.ParseUint(*login, 10, 64)
	if loginTypeParseErr != nil {
		log.Printf("%v", loginTypeParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
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

	// ログイン種別が管理者の場合、チームに関するバリデーション
	if loginType == uint64(static.LOGIN_TYPE_MANAGEMENT) {
		// バリデーション
		if err := u.validator.CreateManagement(req); err != nil {
			log.Printf("%v", err)
			return nil, &response.Error{
				Status: http.StatusBadRequest,
			}
		}
	}

	// チームを取得
	teams, teamsErr := u.team.GetByHashKeys(req.Teams)
	if teamsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	var ids []uint64
	for _, team := range teams {
		ids = append(ids, team.ID)
	}

	// 各チームの面接官割り振り優先順位取得
	priorities, prioritiesErr := u.team.GetAssignPriorityTeams(ids)
	if prioritiesErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// メールアドレス重複チェック
	if err := u.user.EmailDuplCheck(&req.User); err != nil {
		return nil, &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_USER_EMAIL_DUPL,
		}
	}

	// ロール取得
	role, roleErr := u.role.Get(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.RoleHashKey,
		},
	})
	if roleErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 初回パスワード発行
	password, hashPassword, passwordErr := GenerateHash(8, 16)
	if passwordErr != nil {
		log.Printf("%v", passwordErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, hashErr := GenerateHash(1, 25)
	if hashErr != nil {
		log.Printf("%v", hashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	user, userCreateErr := u.user.Insert(tx, &ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   static.PRE_USER + "_" + *hashKey,
			CompanyID: companyID,
		},
		Name:         req.Name,
		Email:        req.Email,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
		RoleID:       role.ID,
		UserType:     static.LOGIN_TYPE_MANAGEMENT,
	})
	if userCreateErr != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム紐づけ一括登録
	var teamAssociations []*ddl.TeamAssociation
	for _, id := range ids {
		teamAssociations = append(teamAssociations, &ddl.TeamAssociation{
			TeamID: id,
			UserID: user.ID,
		})
	}
	if err := u.team.InsertsTeamAssociation(tx, teamAssociations); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	//  面接毎参加可能者登録(全員参加可能)
	for _, team := range teams {
		var possibleList []*ddl.TeamAssignPossible

		for i := 1; i <= int(team.NumOfInterview); i++ {
			possibleList = append(possibleList, &ddl.TeamAssignPossible{
				TeamID:         team.ID,
				UserID:         user.ID,
				NumOfInterview: uint(i),
			})
		}

		if err := u.team.InsertsAssignPossible(tx, possibleList); err != nil {
			if err := u.db.TxRollback(tx); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// 面接割り振り優先順位登録
	for _, id := range ids {
		var count uint
		for _, row := range priorities {
			if row.TeamID == id {
				count++
			}
		}

		if count > 0 {
			var list []*ddl.TeamAssignPriority
			list = append(list, &ddl.TeamAssignPriority{
				TeamID:   id,
				UserID:   user.ID,
				Priority: count + 1,
			})
			if err := u.team.InsertsAssignPriority(tx, list); err != nil {
				if err := u.db.TxRollback(tx); err != nil {
					return nil, &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
		}
	}

	if err := u.db.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	res := response.CreateUser{
		User: entity.User{
			User: ddl.User{
				Email:        user.Email,
				InitPassword: *password,
			},
		},
	}
	return &res, nil
}

// 検索
func (u *UserService) Search(req *request.SearchUser) (*response.SearchUser, *response.Error) {
	// バリデーション
	if err := u.validator.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// Redisから取得
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

	users, usersErr := u.user.Search(&dto.SearchUser{
		SearchUser: *req,
		CompanyID:  companyID,
	})
	if usersErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchUser{
		List: users,
	}, nil
}

// 検索_同一企業
func (u *UserService) SearchByCompany(req *request.SearchUserByCompany) (*response.SearchUserByCompany, *response.Error) {
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
	users, err := u.user.SearchByCompany(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchUserByCompany{
		List: users,
	}, nil
}

// 取得
func (u *UserService) Get(req *request.GetUser) (*response.GetUser, *response.Error) {
	// バリデーション
	if err := u.validator.Get(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	user, err := u.user.Get(&req.User)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.GetUser{
		User: entity.User{
			User: ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: user.HashKey,
				},
				Name:   user.Name,
				Email:  user.Email,
				RoleID: user.RoleID,
			},
		},
	}, nil
}

// ステータスイベントマスタ一覧
func (u *UserService) ListStatusEvent() (*response.ListStatusEvent, *response.Error) {
	res, err := u.master.ListSelectStatusEvent()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.ListStatusEvent{
		List: res,
	}, nil
}

// アサイン関連マスタ取得
func (u *UserService) AssignMaster() (*response.AssignMaster, *response.Error) {
	// アサインルール取得
	rules, rulesErr := u.master.ListAssignRule()
	if rulesErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 自動アサインルール取得
	autoRule, autoRulesErr := u.master.ListAutoAssignRule()
	if autoRulesErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.AssignMaster{
		Rule:     rules,
		AutoRule: autoRule,
	}, nil
}

// 面接官割り振り方法更新
func (u *UserService) UpdateAssignMethod(req *request.UpdateAssignMethod) *response.Error {
	// バリデーション
	if err := u.validatorTeam.UpdateAssignMethod(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}
	for _, row := range req.PossibleList {
		if err := u.validatorTeam.UpdateAssignMethod4(&row); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
			}
		}
	}

	// チームID取得
	ctx := context.Background()
	teamRedis, teamRedisErr := u.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamRedisErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*teamRedis, 10, 64)
	if teamIDErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	team, teamErr := u.team.GetByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ルールハッシュからルール取得
	rule, ruleErr := u.master.SelectAssignRule(&ddl.AssignRule{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: req.RuleHash,
		},
	})
	if ruleErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var autoRule entity.AutoAssignRule
	if rule.ID == static.ASSIGN_RULE_AUTO {
		// 相関バリデーション1
		if err := u.validatorTeam.UpdateAssignMethod2(req); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
			}
		}

		// 自動ルールハッシュから自動ルール取得
		temp, tempErr := u.master.SelectAutoAssignRule(&ddl.AutoAssignRule{
			AbstractMasterModel: ddl.AbstractMasterModel{
				HashKey: req.AutoRuleHash,
			},
		})
		if tempErr != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		autoRule.ID = temp.ID
	}

	var priority []uint64
	if autoRule.ID == static.AUTO_ASSIGN_RULE_ASC {
		// 相関バリデーション2
		if err := u.validatorTeam.UpdateAssignMethod3(req); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
			}
		}

		indexMap := make(map[string]int)
		for i, v := range req.Priority {
			indexMap[v] = i
		}

		sort.Slice(team.Users, func(i, j int) bool {
			return indexMap[team.Users[i].HashKey] < indexMap[team.Users[j].HashKey]
		})

		for _, row := range team.Users {
			priority = append(priority, row.ID)
		}
	}

	// 設定 ＆ 参加可能者ID取得
	var perList []*ddl.TeamPerInterview
	var possibleList []*ddl.TeamAssignPossible
	for _, possible := range req.PossibleList {
		perList = append(perList, &ddl.TeamPerInterview{
			TeamID:         team.ID,
			NumOfInterview: possible.NumOfInterview,
			UserMin:        possible.UserMin,
		})

		for _, user := range team.Users {
			for _, hashKey := range possible.HashKeys {
				if hashKey == user.HashKey {
					possibleList = append(possibleList, &ddl.TeamAssignPossible{
						TeamID:         teamID,
						NumOfInterview: possible.NumOfInterview,
						UserID:         user.ID,
					})
				}
			}
		}
	}

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ルール更新
	_, updateTeamErr := u.team.Update(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: team.HashKey,
		},
		RuleID: rule.ID,
	})
	if updateTeamErr != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 自動ルール削除
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

	// 優先順位削除
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

	// 面接毎設定削除
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

	// 面接参加可能者削除
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

	// 面接毎設定一括登録
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

	// 面接参加可能者一括登録
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

	// 自動ルール更新
	if rule.ID == static.ASSIGN_RULE_AUTO && autoRule.ID > 0 {
		// 登録
		if err := u.team.InsertAutoAssignRule(tx, &ddl.TeamAutoAssignRule{
			TeamID: team.ID,
			RuleID: autoRule.ID,
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

		// 優先順位更新
		if autoRule.ID == static.AUTO_ASSIGN_RULE_ASC && len(priority) == len(team.Users) {
			var list []*ddl.TeamAssignPriority
			for index, row := range priority {
				list = append(list, &ddl.TeamAssignPriority{
					TeamID:   team.ID,
					UserID:   row,
					Priority: uint(index + 1),
				})
			}

			// 一括登録
			if err := u.team.InsertsAssignPriority(tx, list); err != nil {
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

// 書類提出ルールマスタ取得
func (u *UserService) DocumentRuleMaster() (*response.DocumentRule, *response.Error) {
	res, err := u.master.ListDocumentRule()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	for _, r := range res {
		r.ID = 0
	}

	return &response.DocumentRule{
		List: res,
	}, nil
}

// 職種マスタ取得
func (u *UserService) OccupationMaster() (*response.Occupation, *response.Error) {
	res, err := u.master.ListOccupation()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	for _, r := range res {
		r.ID = 0
	}

	return &response.Occupation{
		List: res,
	}, nil
}
