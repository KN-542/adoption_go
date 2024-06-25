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
	"time"
)

type IUserService interface {
	// 登録
	Create(req *request.CreateUser) (*response.CreateUser, *response.Error)
	// 検索
	Search(req *request.SearchUser) (*response.SearchUser, *response.Error)
	// 取得
	Get(req *request.GetUser) (*response.GetUser, *response.Error)
	// チーム登録
	CreateTeam(req *request.CreateTeam) *response.Error
	// チーム検索
	SearchTeam(req *request.SearchTeam) (*response.SearchTeam, *response.Error)
	// 予定登録種別一覧
	SearchScheduleType() (*response.SearchScheduleType, *response.Error)
	// 予定登録
	CreateSchedule(req *request.CreateSchedule) (*response.CreateSchedule, *response.Error)
	// 予定更新
	UpdateSchedule(req *request.UpdateSchedule) *response.Error
	// 予定検索
	SearchSchedule(req *request.SearchSchedule) (*response.SearchSchedule, *response.Error)
	// 予定削除
	DeleteSchedule(req *request.DeleteSchedule) *response.Error
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

		// メールアドレス重複チェック
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

	// チーム取得
	ctx := context.Background()
	team, teamErr := u.redis.Get(ctx, req.HashKey, static.REDIS_USER_TEAM_ID)
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

	users, usersErr := u.user.Search(&dto.SearchUser{
		SearchUser: *req,
		TeamID:     teamID,
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
	req.CompanyID = companyID

	// チーム取得
	teams, err := u.user.SearchTeam(&req.Team)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 所属ユーザー取得
	var res []entity.SearchTeam
	for _, team := range teams {
		var t entity.SearchTeam

		t.Team = team.Team
		for _, row := range team.Users {
			t.Users = append(t.Users, &ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: row.HashKey,
				},
				Name: row.Name,
			})
		}
		res = append(res, t)
	}

	return &response.SearchTeam{
		List: res,
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
	req.HashKey = *hashKey
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

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
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
func (u *UserService) CreateSchedule(req *request.CreateSchedule) (*response.CreateSchedule, *response.Error) {
	// バリデーション
	if err := u.v.CreateSchedule(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー存在確認
	ids, idsErr := u.user.GetIDs(req.Users)
	if idsErr != nil {
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

	// 予定登録
	scheduleID, err := u.user.InsertSchedule(tx, &ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *hashKey,
		},
		InterviewFlg: req.InterviewFlg,
		FreqID:       req.FreqID,
		Start:        req.Start,
		End:          req.End,
		Title:        req.Title,
	})
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var userScheduleAssociations []*ddl.UserScheduleAssociation
	for _, id := range ids {
		userScheduleAssociations = append(userScheduleAssociations, &ddl.UserScheduleAssociation{
			UserScheduleID: *scheduleID,
			UserID:         id,
		})
	}
	// 予定紐づけ一括登録
	if err := u.user.InsertsScheduleAssociation(tx, userScheduleAssociations); err != nil {
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

	return &response.CreateSchedule{
		HashKey: *hashKey,
	}, nil
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
	scheduleID, err := u.user.UpdateSchedule(tx, &ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   req.HashKey,
			UpdatedAt: time.Now(),
		},
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
	// 問答無用で紐づけテーブルの該当予定IDのレコード削除
	if err := u.user.DeleteScheduleAssociation(tx, &ddl.UserScheduleAssociation{
		UserScheduleID: *scheduleID,
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

	var userScheduleAssociations []*ddl.UserScheduleAssociation
	for _, id := range ids {
		userScheduleAssociations = append(userScheduleAssociations, &ddl.UserScheduleAssociation{
			UserScheduleID: *scheduleID,
			UserID:         id,
		})
	}
	// 予定紐づけ一括登録
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
	// バリデーション
	if err := u.v.SearchSchedule(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

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

	schedulesBefore, sErr := u.user.SearchSchedule(&ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
		},
	})
	if sErr != nil {
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

	// 日付が過去の場合、更新or削除
	for _, schedule := range schedulesBefore {
		if schedule.Start.Before(time.Now()) {
			// なしの場合
			if schedule.FreqID == uint(static.FREQ_NONE) {
				if err := u.user.DeleteSchedule(tx, &ddl.UserSchedule{
					AbstractTransactionModel: ddl.AbstractTransactionModel{
						HashKey: schedule.HashKey,
					},
				}); err != nil {
					if err := u.d.TxRollback(tx); err != nil {
						return nil, &response.Error{
							Status: http.StatusInternalServerError,
						}
					}
					return nil, &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
			} else {
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

				_, updateErr := u.user.UpdateSchedule(tx, &ddl.UserSchedule{
					AbstractTransactionModel: ddl.AbstractTransactionModel{
						HashKey:   schedule.HashKey,
						UpdatedAt: time.Now(),
					},
					Start: s,
					End:   e,
				})
				if updateErr != nil {
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
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	schedulesAfter, err := u.user.SearchSchedule(&ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.UserSchedule
	for _, row := range schedulesAfter {
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

	tx, err := u.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := u.user.DeleteSchedule(tx, &req.UserSchedule); err != nil {
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
