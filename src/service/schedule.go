package service

import (
	"api/src/model/ddl"
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

type IScheduleService interface {
	// 予定登録種別一覧
	SearchScheduleType() (*response.SearchScheduleType, *response.Error)
	// 予定登録
	Create(req *request.CreateSchedule) *response.Error
	// 予定更新
	Update(req *request.UpdateSchedule) *response.Error
	// 予定検索
	Search(req *request.SearchSchedule) (*response.SearchSchedule, *response.Error)
	// 予定削除
	Delete(req *request.DeleteSchedule) *response.Error
}

type ScheduleService struct {
	db         repository.IDBRepository
	redis      repository.IRedisRepository
	user       repository.IUserRepository
	team       repository.ITeamRepository
	schedule   repository.IScheduleRepository
	applicant  repository.IApplicantRepository
	role       repository.IRoleRepository
	manuscript repository.IManuscriptRepository
	master     repository.IMasterRepository
	v          validator.IScheduleValidator
	outer      repository.IOuterIFRepository
}

func NewScheduleService(
	db repository.IDBRepository,
	redis repository.IRedisRepository,
	user repository.IUserRepository,
	team repository.ITeamRepository,
	schedule repository.IScheduleRepository,
	applicant repository.IApplicantRepository,
	role repository.IRoleRepository,
	manuscript repository.IManuscriptRepository,
	master repository.IMasterRepository,
	v validator.IScheduleValidator,
	outer repository.IOuterIFRepository,
) IScheduleService {
	return &ScheduleService{db, redis, user, team, schedule, applicant, role, manuscript, master, v, outer}
}

// 予定登録種別一覧
func (u *ScheduleService) SearchScheduleType() (*response.SearchScheduleType, *response.Error) {
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
func (u *ScheduleService) Create(req *request.CreateSchedule) *response.Error {
	// バリデーション
	if err := u.v.Create(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
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

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定登録
	scheduleID, err := u.schedule.Insert(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   static.PRE_SCHEDULE + "_" + *hashKey,
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
		if err := u.db.TxRollback(tx); err != nil {
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
	if err := u.schedule.InsertsScheduleAssociation(tx, userScheduleAssociations); err != nil {
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

// 予定更新
func (u *ScheduleService) Update(req *request.UpdateSchedule) *response.Error {
	// バリデーション
	if err := u.v.Update(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 予定取得
	schedule, scheduleErr := u.schedule.Get(&ddl.Schedule{
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

	tx, txErr := u.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定更新
	if err := u.schedule.Update(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   schedule.HashKey,
			UpdatedAt: time.Now(),
		},
		FreqID: req.FreqID,
		Start:  req.Start,
		End:    req.End,
		Title:  req.Title,
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
	// 問答無用で紐づけテーブルの該当予定IDのレコード削除
	if err := u.schedule.DeleteScheduleAssociation(tx, &ddl.ScheduleAssociation{
		ScheduleID: schedule.ID,
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

	// 予定紐づけ一括登録
	var userScheduleAssociations []*ddl.ScheduleAssociation
	for _, id := range ids {
		userScheduleAssociations = append(userScheduleAssociations, &ddl.ScheduleAssociation{
			ScheduleID: schedule.ID,
			UserID:     id,
		})
	}
	if err := u.schedule.InsertsScheduleAssociation(tx, userScheduleAssociations); err != nil {
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

// 予定検索 (バッチでも実行したい)
func (u *ScheduleService) Search(req *request.SearchSchedule) (*response.SearchSchedule, *response.Error) {
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

	schedulesBefore, sErr := u.schedule.Search(&ddl.Schedule{
		TeamID: teamID,
	})
	if sErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var deleteList []uint64
	var editList []*ddl.Schedule

	tx, txErr := u.db.TxStart()
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

	if err := u.schedule.Deletes(tx, deleteList); err != nil {
		if err := u.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if len(editList) > 0 {
		if err := u.schedule.Inserts(tx, editList); err != nil {
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

	if err := u.db.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	schedulesAfter, err := u.schedule.Search(&ddl.Schedule{
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
func (u *ScheduleService) Delete(req *request.DeleteSchedule) *response.Error {
	// バリデーション
	if err := u.v.Delete(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 予定取得
	schedule, scheduleErr := u.schedule.Get(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if scheduleErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.db.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 紐づけ削除
	if err := u.schedule.DeleteScheduleAssociation(tx, &ddl.ScheduleAssociation{
		ScheduleID: schedule.ID,
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
	if err := u.schedule.Delete(tx, &ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: schedule.HashKey,
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
