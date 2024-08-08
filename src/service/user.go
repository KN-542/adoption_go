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
	"time"
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
	// チーム登録
	CreateTeam(req *request.CreateTeam) *response.Error
	// チーム更新
	UpdateTeam(req *request.UpdateTeam) *response.Error
	// チーム基本情報更新
	UpdateBasicTeam(req *request.UpdateBasicTeam) *response.Error
	// チーム削除
	DeleteTeam(req *request.DeleteTeam) *response.Error
	// チーム検索
	SearchTeam(req *request.SearchTeam) (*response.SearchTeam, *response.Error)
	// チーム取得
	GetTeam(req *request.GetTeam) (*response.GetTeam, *response.Error)
	// 自チーム取得
	GetOwnTeam(req *request.GetOwnTeam) (*response.GetOwnTeam, *response.Error)
	// チーム検索_同一企業
	SearchTeamByCompany(req *request.SearchTeamByCompany) (*response.SearchTeamByCompany, *response.Error)
	// 予定登録種別一覧
	SearchScheduleType() (*response.SearchScheduleType, *response.Error)
	// 予定登録
	CreateSchedule(req *request.CreateSchedule) *response.Error
	// 予定更新
	UpdateSchedule(req *request.UpdateSchedule) *response.Error
	// 予定検索
	SearchSchedule(req *request.SearchSchedule) (*response.SearchSchedule, *response.Error)
	// 予定削除
	DeleteSchedule(req *request.DeleteSchedule) *response.Error
	// ステータスイベントマスタ一覧
	ListStatusEvent() (*response.ListStatusEvent, *response.Error)
	// チーム毎ステータスイベント取得
	StatusEventsByTeam(req *request.StatusEventsByTeam) (*response.StatusEventsByTeam, *response.Error)
	// アサイン関連マスタ取得
	AssignMaster() (*response.AssignMaster, *response.Error)
	// 面接官割り振り方法更新
	UpdateAssignMethod(req *request.UpdateAssignMethod) *response.Error
}

type UserService struct {
	user      repository.IUserRepository
	role      repository.IRoleRepository
	applicant repository.IApplicantRepository
	master    repository.IMasterRepository
	v         validator.IUserValidator
	d         repository.IDBRepository
	outer     repository.IOuterIFRepository
	redis     repository.IRedisRepository
}

func NewUserService(
	user repository.IUserRepository,
	role repository.IRoleRepository,
	applicant repository.IApplicantRepository,
	master repository.IMasterRepository,
	v validator.IUserValidator,
	d repository.IDBRepository,
	outer repository.IOuterIFRepository,
	redis repository.IRedisRepository,
) IUserService {
	return &UserService{user, role, applicant, master, v, d, outer, redis}
}

// 登録
func (u *UserService) Create(req *request.CreateUser) (*response.CreateUser, *response.Error) {
	// バリデーション
	if err := u.v.Create(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

	// チームのIDを取得
	teams, teamsErr := u.user.GetTeamIDs(req.Teams)
	if teamsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ログイン種別が管理者の場合、チームに関するバリデーション
	if loginType == uint64(static.LOGIN_TYPE_MANAGEMENT) {
		// バリデーション
		if err := u.v.CreateManagement(req); err != nil {
			log.Printf("%v", err)
			return nil, &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
			}
		}

		// メールアドレス重複チェック_管理者
		if err := u.user.EmailDuplCheckManagement(&req.User, teams); err != nil {
			return nil, &response.Error{
				Status: http.StatusConflict,
				Code:   static.CODE_USER_EMAIL_DUPL,
			}
		}
	} else {
		// メールアドレス重複チェック
		if err := u.user.EmailDuplCheck(&req.User); err != nil {
			return nil, &response.Error{
				Status: http.StatusConflict,
				Code:   static.CODE_USER_EMAIL_DUPL,
			}
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

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	user, userCreateErr := u.user.Insert(tx, &ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   *hashKey,
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
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var teamAssociations []*ddl.TeamAssociation
	for _, id := range teams {
		teamAssociations = append(teamAssociations, &ddl.TeamAssociation{
			TeamID: id,
			UserID: user.ID,
		})
	}
	// チーム紐づけ一括登録
	if err := u.user.InsertsTeamAssociation(tx, teamAssociations); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
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
	if err := u.v.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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
	if err := u.v.Get(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

// チーム検索
func (u *UserService) SearchTeam(req *request.SearchTeam) (*response.SearchTeam, *response.Error) {
	// バリデーション
	if err := u.v.SearchTeam(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

	// チーム検索
	teams, err := u.user.SearchTeam(&req.Team)
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

// チーム取得
func (u *UserService) GetTeam(req *request.GetTeam) (*response.GetTeam, *response.Error) {
	// バリデーション
	if err := u.v.GetTeam(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 取得
	team, err := u.user.GetTeam(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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
func (u *UserService) GetOwnTeam(req *request.GetOwnTeam) (*response.GetOwnTeam, *response.Error) {
	// チームID取得
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
	res, resErr := u.user.GetTeamByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if resErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	for _, row := range res.Users {
		row.ID = 0
	}

	// チーム面接毎イベント取得
	events, eventsErr := u.user.InterviewEventsByTeam(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if eventsErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 面接自動割り当てルール取得
	var autoRule entity.TeamAutoAssignRule
	if res.RuleID == static.ASSIGN_RULE_AUTO {
		temp, tempErr := u.user.GetAutoAssignRule(&ddl.TeamAutoAssignRule{
			TeamID: teamID,
		})
		if tempErr != nil {
			return nil, &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
			}
		}
		autoRule = *temp
	}

	// 面接割り振り優先順位取得
	var priority []entity.TeamAssignPriority
	if autoRule.RuleID == static.AUTO_ASSIGN_RULE_ASC {
		tempList, tempErr := u.user.GetAssignPriority(&ddl.TeamAssignPriority{
			TeamID: teamID,
		})
		if tempErr != nil {
			return nil, &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
			}
		}
		for _, row := range tempList {
			priority = append(priority, *row)
		}
	}

	// 面接毎参加可能者取得
	possibleList, possibleErr := u.user.GetAssignPossible(&ddl.TeamAssignPossible{
		TeamID: teamID,
	})
	if possibleErr != nil {
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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
				UserMin:        res.UserMin,
			},
			RuleHash: res.RuleHash,
			Users:    res.Users,
		},
		Events: events,
		AutoRule: entity.TeamAutoAssignRule{
			HashKey: autoRule.HashKey,
		},
		Priority:     priority,
		PossibleList: possibleList,
	}, nil
}

// チーム登録
func (u *UserService) CreateTeam(req *request.CreateTeam) *response.Error {
	// バリデーション
	if err := u.v.CreateTeam(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム登録
	req.HashKey = string(static.PRE_TEAM) + "_" + *hashKey
	req.NumOfInterview = 3
	req.UserMin = 1
	req.RuleID = static.ASSIGN_RULE_MANUAL
	team, err := u.user.InsertTeam(tx, &req.Team)
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
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
	// チーム紐づけ一括登録
	if err := u.user.InsertsTeamAssociation(tx, teamAssociations); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接毎参加可能者登録(全員参加可能)
	var possibleList []*ddl.TeamAssignPossible
	for i := 1; i <= int(team.NumOfInterview); i++ {
		for _, id := range ids {
			possibleList = append(possibleList, &ddl.TeamAssignPossible{
				TeamID:         team.ID,
				UserID:         id,
				NumOfInterview: uint(i),
			})
		}
	}
	if err := u.user.InsertsAssignPossible(tx, possibleList); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
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
	if err := u.user.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash,
			CompanyID: companyID,
		},
		TeamID:     team.ID,
		StatusName: "日程未回答",
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if err := u.user.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash2,
			CompanyID: companyID,
		},
		TeamID:     team.ID,
		StatusName: "日程回答済み",
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// チーム更新
func (u *UserService) UpdateTeam(req *request.UpdateTeam) *response.Error {
	// バリデーション
	if err := u.v.UpdateTeam(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム更新
	team, err := u.user.UpdateTeam(tx, &req.Team)
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
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
		for _, id := range ids {
			teamAssociations = append(teamAssociations, &ddl.TeamAssociation{
				TeamID: team.ID,
				UserID: id,
			})
		}

		// チーム紐づけ一括登録
		if err := u.user.InsertsTeamAssociation(tx, teamAssociations); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// チーム基本情報更新
func (u *UserService) UpdateBasicTeam(req *request.UpdateBasicTeam) *response.Error {
	// バリデーション
	if err := u.v.UpdateBasicTeam(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// チームID取得
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

	// チーム取得
	team, err := u.user.GetTeamByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: teamID,
		},
	})
	if err != nil {
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}
	req.HashKey = team.HashKey

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム更新
	_, updateErr := u.user.UpdateTeam(tx, &req.Team)
	if updateErr != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// チーム削除
func (u *UserService) DeleteTeam(req *request.DeleteTeam) *response.Error {
	// バリデーション
	if err := u.v.DeleteTeam(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// チーム取得
	team, teamErr := u.user.GetTeam(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除可能判定
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
			Code:   static.CODE_TEAM_USER_CANNOT_DELETE,
		}
	}

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 紐づけ削除
	if err := u.user.DeleteTeamAssociation(tx, &ddl.TeamAssociation{
		TeamID: team.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 選考状況削除
	if err := u.applicant.DeleteStatus(tx, &ddl.SelectStatus{
		TeamID: team.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム削除
	if err := u.user.DeleteTeam(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// チーム検索_同一企業
func (u *UserService) SearchTeamByCompany(req *request.SearchTeamByCompany) (*response.SearchTeamByCompany, *response.Error) {
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
	team, err := u.user.SearchTeamByCompany(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
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

// 予定登録種別一覧
func (u *UserService) SearchScheduleType() (*response.SearchScheduleType, *response.Error) {
	res, err := u.master.SelectScheduleFreqStatus()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchScheduleType{
		List: res,
	}, nil
}

// 予定登録
func (u *UserService) CreateSchedule(req *request.CreateSchedule) *response.Error {
	// バリデーション
	if err := u.v.CreateSchedule(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// Redisから取得
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

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定登録
	scheduleID, err := u.user.InsertSchedule(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   *hashKey,
			CompanyID: companyID,
		},
		InterviewFlg: req.InterviewFlg,
		FreqID:       req.FreqID,
		Start:        req.Start,
		End:          req.End,
		Title:        req.Title,
		TeamID:       teamID,
	})
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定紐づけ一括登録
	var userScheduleAssociations []*ddl.ScheduleAssociation
	for _, id := range ids {
		userScheduleAssociations = append(userScheduleAssociations, &ddl.ScheduleAssociation{
			ScheduleID: *scheduleID,
			UserID:     id,
		})
	}
	if err := u.user.InsertsScheduleAssociation(tx, userScheduleAssociations); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 予定更新
func (u *UserService) UpdateSchedule(req *request.UpdateSchedule) *response.Error {
	// バリデーション
	if err := u.v.UpdateSchedule(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 予定取得
	schedule, scheduleErr := u.user.GetSchedule(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if scheduleErr != nil {
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

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定更新
	if err := u.user.UpdateSchedule(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   schedule.HashKey,
			UpdatedAt: time.Now(),
		},
		FreqID: req.FreqID,
		Start:  req.Start,
		End:    req.End,
		Title:  req.Title,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// 問答無用で紐づけテーブルの該当予定IDのレコード削除
	if err := u.user.DeleteScheduleAssociation(tx, &ddl.ScheduleAssociation{
		ScheduleID: schedule.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定紐づけ一括登録
	var userScheduleAssociations []*ddl.ScheduleAssociation
	for _, id := range ids {
		userScheduleAssociations = append(userScheduleAssociations, &ddl.ScheduleAssociation{
			ScheduleID: schedule.ID,
			UserID:     id,
		})
	}
	if err := u.user.InsertsScheduleAssociation(tx, userScheduleAssociations); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 予定検索 (バッチでも実行したい)
func (u *UserService) SearchSchedule(req *request.SearchSchedule) (*response.SearchSchedule, *response.Error) {
	// チームID取得
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

	schedulesBefore, sErr := u.user.SearchSchedule(&ddl.Schedule{
		TeamID: teamID,
	})
	if sErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var deleteList []uint64
	var editList []*ddl.Schedule

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 日付が過去の場合、更新or削除
	if len(schedulesBefore) > 0 {
		for _, schedule := range schedulesBefore {
			if schedule.Start.Before(time.Now()) {
				deleteList = append(deleteList, schedule.ID)

				// なし以外の場合
				if schedule.FreqID != uint(static.FREQ_NONE) {
					s := schedule.Start
					e := schedule.End
					if schedule.FreqID == uint(static.FREQ_DAILY) {
						s = s.AddDate(0, 0, 1)
						e = e.AddDate(0, 0, 1)
					}
					if schedule.FreqID == uint(static.FREQ_WEEKLY) {
						s = s.AddDate(0, 0, 7)
						e = e.AddDate(0, 0, 7)
					}
					if schedule.FreqID == uint(static.FREQ_MONTHLY) {
						s = s.AddDate(0, 1, 0)
						e = e.AddDate(0, 1, 0)
					}
					if schedule.FreqID == uint(static.FREQ_YEARLY) {
						s = s.AddDate(1, 0, 0)
						e = e.AddDate(1, 0, 0)
					}

					editList = append(editList, &ddl.Schedule{
						AbstractTransactionModel: ddl.AbstractTransactionModel{
							HashKey:   schedule.HashKey,
							CompanyID: schedule.CompanyID,
							CreatedAt: schedule.CreatedAt,
							UpdatedAt: time.Now(),
						},
						Start:        s,
						End:          e,
						Title:        schedule.Title,
						FreqID:       schedule.FreqID,
						InterviewFlg: schedule.InterviewFlg,
						TeamID:       schedule.TeamID,
					})
				}
			}
		}
	}

	if err := u.user.DeletesSchedule(tx, deleteList); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if len(editList) > 0 {
		if err := u.user.InsertsSchedule(tx, editList); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	schedulesAfter, err := u.user.SearchSchedule(&ddl.Schedule{
		TeamID: teamID,
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.Schedule
	for _, row := range schedulesAfter {
		row.ID = 0
		for _, row2 := range row.Users {
			row2.ID = 0
		}

		res = append(res, *row)
	}

	return &response.SearchSchedule{
		List: res,
	}, nil
}

// 予定削除
func (u *UserService) DeleteSchedule(req *request.DeleteSchedule) *response.Error {
	// バリデーション
	if err := u.v.DeleteSchedule(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 予定取得
	schedule, scheduleErr := u.user.GetSchedule(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if scheduleErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 紐づけ削除
	if err := u.user.DeleteScheduleAssociation(tx, &ddl.ScheduleAssociation{
		ScheduleID: schedule.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := u.user.DeleteSchedule(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: schedule.HashKey,
		},
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
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

// チーム毎ステータスイベント取得
func (u *UserService) StatusEventsByTeam(req *request.StatusEventsByTeam) (*response.StatusEventsByTeam, *response.Error) {
	// チーム取得
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
	res, err := u.user.StatusEventsByTeam(&ddl.Team{
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
	if err := u.v.UpdateAssignMethod(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}
	for _, row := range req.PossibleList {
		if err := u.v.UpdateAssignMethod4(&row); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
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
	team, teamErr := u.user.GetTeamByPrimary(&ddl.Team{
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
		if err := u.v.UpdateAssignMethod2(req); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
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
		if err := u.v.UpdateAssignMethod3(req); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
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

	// 参加可能者ID取得
	var possibleList []*ddl.TeamAssignPossible
	for _, possible := range req.PossibleList {
		newPossible := &ddl.TeamAssignPossible{
			TeamID:         teamID,
			NumOfInterview: possible.NumOfInterview,
		}

		for _, user := range team.Users {
			if possible.HashKey == user.HashKey {
				newPossible.UserID = user.ID
				break
			}
		}

		possibleList = append(possibleList, newPossible)
	}

	tx, txErr := u.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ルール更新
	_, updateTeamErr := u.user.UpdateTeam(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: team.HashKey,
		},
		RuleID:  rule.ID,
		UserMin: req.UserMin,
	})
	if updateTeamErr != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 自動ルール削除
	if err := u.user.DeleteAutoAssignRule(tx, &ddl.TeamAutoAssignRule{
		TeamID: team.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 優先順位削除
	if err := u.user.DeleteAssignPriority(tx, &ddl.TeamAssignPriority{
		TeamID: team.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接参加可能者削除
	if err := u.user.DeleteAssignPossible(tx, &ddl.TeamAssignPossible{
		TeamID: team.ID,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接参加可能者一括登録
	if err := u.user.InsertsAssignPossible(tx, possibleList); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
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
		if err := u.user.InsertAutoAssignRule(tx, &ddl.TeamAutoAssignRule{
			TeamID: team.ID,
			RuleID: autoRule.ID,
		}); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
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
			if err := u.user.InsertsAssignPriority(tx, list); err != nil {
				if err := u.d.TxRollback(tx); err != nil {
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

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
