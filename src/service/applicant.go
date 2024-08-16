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
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type IApplicantService interface {
	// 検索
	Search(req *request.SearchApplicant) (*response.SearchApplicant, *response.Error)
	// 取得
	Get(req *request.GetApplicant) (*response.GetApplicant, *response.Error)
	// サイト一覧取得
	GetSites() (*response.ApplicantSites, *response.Error)
	// 応募者ステータス一覧取得
	GetStatusList(req *request.ApplicantStatusList) (*response.ApplicantStatusList, *response.Error)
	// 応募者ダウンロード
	Download(req *request.ApplicantDownload) (*response.ApplicantDownload, *response.Error)
	// 予約表表示
	ReserveTable(req *request.ReserveTable) (*response.ReserveTable, *response.Error)
	// 書類アップロード(S3)
	S3Upload(req *request.FileUpload, fileHeader *multipart.FileHeader) *response.Error
	// 書類ダウンロード(S3)
	S3Download(req *request.FileDownload) ([]byte, *string, *response.Error)
	// 面接希望日登録
	InsertDesiredAt(req *request.InsertDesiredAt) *response.Error
	// 認証URL作成
	GetOauthURL(req *request.GetOauthURL) (*response.GetOauthURL, *response.Error)
	// GoogleMeetUrl発行
	GetGoogleMeetUrl(req *request.GetGoogleMeetUrl) (*response.GetGoogleMeetUrl, *response.Error)
	// 応募者ステータス変更
	UpdateStatus(req *request.UpdateStatus) *response.Error
	// 面接官割り振り
	AssignUser(req *request.AssignUser) *response.Error
}

type ApplicantService struct {
	r     repository.IApplicantRepository
	u     repository.IUserRepository
	m     repository.IMasterRepository
	a     repository.IAWSRepository
	g     repository.IGoogleRepository
	redis repository.IRedisRepository
	v     validator.IApplicantValidator
	d     repository.IDBRepository
	o     repository.IOuterIFRepository
}

func NewApplicantService(
	r repository.IApplicantRepository,
	u repository.IUserRepository,
	m repository.IMasterRepository,
	a repository.IAWSRepository,
	g repository.IGoogleRepository,
	redis repository.IRedisRepository,
	v validator.IApplicantValidator,
	d repository.IDBRepository,
	o repository.IOuterIFRepository,
) IApplicantService {
	return &ApplicantService{r, u, m, a, g, redis, v, d, o}
}

// 検索
func (s *ApplicantService) Search(req *request.SearchApplicant) (*response.SearchApplicant, *response.Error) {
	// バリデーション
	if err := s.v.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// Redisから取得
	ctx := context.Background()

	team, teamErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
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
	req.TeamID = teamID

	company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
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
	applicants, num, searchErr := s.r.Search(&dto.SearchApplicant{
		SearchApplicant: *req,
		Users:           req.Users,
	})
	if searchErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.SearchApplicant
	for _, applicant := range applicants {
		var filteredUsers []*ddl.User
		for _, user := range applicant.Users {
			u := &ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: user.HashKey,
				},
				Name: user.Name,
			}
			filteredUsers = append(filteredUsers, u)
		}
		applicant.Users = filteredUsers
		res = append(res, *applicant)
	}

	return &response.SearchApplicant{
		List: res,
		Num:  num,
	}, nil
}

// 取得
func (s *ApplicantService) Get(req *request.GetApplicant) (*response.GetApplicant, *response.Error) {
	// バリデーション
	if err := s.v.Get(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者情報取得
	applicant, err := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.GetApplicant{
		Applicant: entity.Applicant{
			Applicant: ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: applicant.HashKey,
				},
				SiteID: applicant.SiteID,
				Status: applicant.Status,
				Name:   applicant.Name,
				Email:  applicant.Email,
				Tel:    applicant.Tel,
				Age:    applicant.Age,
			},
			GoogleMeetURL: applicant.GoogleMeetURL,
		},
	}, nil
}

// サイト一覧取得
func (s *ApplicantService) GetSites() (*response.ApplicantSites, *response.Error) {
	sites, err := s.m.ListSite()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range sites {
		row.ID = 0
	}

	return &response.ApplicantSites{
		List: sites,
	}, nil
}

// 応募者ステータス一覧取得
func (s *ApplicantService) GetStatusList(req *request.ApplicantStatusList) (*response.ApplicantStatusList, *response.Error) {
	// バリデーション
	if err := s.v.GetStatusList(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// チーム取得
	ctx := context.Background()
	team, teamErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
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

	// ステータス取得
	list, listErr := s.r.ListStatus(&ddl.SelectStatus{
		TeamID: teamID,
	})
	if listErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.ApplicantStatus
	for _, row := range list {
		res = append(res, entity.ApplicantStatus{
			SelectStatus: ddl.SelectStatus{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: row.HashKey,
				},
				StatusName: row.StatusName,
			},
		})
	}

	return &response.ApplicantStatusList{
		List: res,
	}, nil
}

// 応募者ダウンロード
func (s *ApplicantService) Download(req *request.ApplicantDownload) (*response.ApplicantDownload, *response.Error) {
	// バリデーション
	if err := s.v.Download(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}
	for _, row := range req.Applicants {
		if err := s.v.DownloadSub(&row); err != nil {
			log.Printf("%v", err)
			return nil, &response.Error{
				Status: http.StatusBadRequest,
			}
		}
	}

	// チーム、企業取得
	ctx := context.Background()
	team, teamErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
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
	company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
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

	// ステータス取得
	list, listErr := s.r.ListStatus(&ddl.SelectStatus{
		TeamID: teamID,
	})
	if listErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 重複チェック
	var request request.ApplicantDownload
	var outerIDs []string

	for _, row := range req.Applicants {
		outerIDs = append(outerIDs, row.OuterID)
	}

	duplApplicants, duplApplicantsErr := s.r.CheckDuplByOuterId(&dto.CheckDuplDownloading{
		TeamID:    teamID,
		CompanyID: companyID,
		List:      outerIDs,
	})
	if duplApplicantsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	for _, row := range req.Applicants {
		dupl := false
		for _, applicant := range duplApplicants {
			if row.OuterID == applicant.OuterID {
				dupl = true
				break
			}
		}
		if !dupl {
			request.Applicants = append(request.Applicants, row)
		}
	}

	// サイトID取得
	site, siteErr := s.m.SelectSite(&ddl.Site{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: req.SiteHashKey,
		},
	})
	if siteErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 一括登録
	var applicants []*ddl.Applicant
	if len(request.Applicants) > 0 {
		for _, row := range request.Applicants {
			// ハッシュキー生成
			_, hash, hashErr := GenerateHash(1, 25)
			if hashErr != nil {
				if err := s.d.TxRollback(tx); err != nil {
					return nil, &response.Error{
						Status: http.StatusInternalServerError,
					}
				}

				log.Printf("%v", hashErr)
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			// 構造体生成
			applicant := &ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey:   *hash,
					CompanyID: companyID,
				},
				OuterID:        row.OuterID,
				SiteID:         site.ID,
				Status:         list[0].ID,
				Name:           row.Name,
				Email:          row.Email,
				Tel:            row.Tel,
				Age:            uint(row.Age),
				TeamID:         teamID,
				NumOfInterview: 1,
			}
			applicants = append(applicants, applicant)
		}

		if err := s.r.Inserts(tx, applicants); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.ApplicantDownload{
		UpdateNum: len(request.Applicants),
	}, nil
}

// 予約表表示 ※TODO 要方針決定
func (s *ApplicantService) ReserveTable(req *request.ReserveTable) (*response.ReserveTable, *response.Error) {
	const WEEKS = 7
	const RESERVE_DURATION = 2 * WEEKS
	var schedules []entity.Schedule

	// バリデーション
	if err := s.v.ReserveTable(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// TZをAsia/Tokyoに
	jst, jstErr := time.LoadLocation("Asia/Tokyo")
	if jstErr != nil {
		log.Printf("%v", jstErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 日本の休日取得
	holidays, hErr := s.o.HolidaysJp(time.Now().Year())
	if hErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者取得
	applicant, applicantErr := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if applicantErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接参加可能者取得
	models, modelsErr := s.u.GetAssignPossibleSchedule(&ddl.TeamAssignPossible{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if modelsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 整形
	for _, model := range models {
		// TZを日本に
		var schedulesJST []entity.Schedule
		for _, row := range model.Schedules {
			row.Start = row.Start.In(jst)
			row.End = row.End.In(jst)
			schedulesJST = append(schedulesJST, *row)
		}

		// スケジュールの頻度が「毎日」と「毎週」の場合、コピー
		start := time.Now().AddDate(0, 0, WEEKS).In(jst)
		s := time.Date(
			start.Year(),
			start.Month(),
			start.Day(),
			0,
			0,
			0,
			0,
			jst,
		)
		for _, row := range schedulesJST {
			if row.FreqID == uint(static.FREQ_NONE) || row.FreqID == uint(static.FREQ_MONTHLY) || row.FreqID == uint(static.FREQ_YEARLY) {
				schedules = append(schedules, row)
				continue
			}

			for i := 0; i < RESERVE_DURATION; i++ {
				s_0 := time.Date(
					s.AddDate(0, 0, i).Year(),
					s.AddDate(0, 0, i).Month(),
					s.AddDate(0, 0, i).Day(),
					row.Start.Hour(),
					row.Start.Minute(),
					row.Start.Second(),
					row.Start.Nanosecond(),
					jst,
				)
				e_0 := time.Date(
					s.AddDate(0, 0, i).Year(),
					s.AddDate(0, 0, i).Month(),
					s.AddDate(0, 0, i).Day(),
					row.End.Hour(),
					row.End.Minute(),
					row.End.Second(),
					row.End.Nanosecond(),
					jst,
				)
				if row.FreqID == uint(static.FREQ_DAILY) || (row.FreqID == uint(static.FREQ_WEEKLY) && s_0.Weekday() == row.Start.Weekday()) {
					row.Start = s_0
					row.End = e_0
					schedules = append(schedules, row)
				}
			}
		}
	}

	// 各日チェック
	start := time.Now().AddDate(0, 0, WEEKS).In(jst)
	var times []time.Time
	var reserveTime []dto.ReserveTableSub
	for i := 0; i < RESERVE_DURATION; i++ {
		isReserve := true

		s_0 := time.Date(
			start.AddDate(0, 0, i).Year(),
			start.AddDate(0, 0, i).Month(),
			start.AddDate(0, 0, i).Day(),
			0,
			0,
			0,
			0,
			jst,
		)

		times = append(times, s_0)

		for _, holiday := range holidays {
			y1, m1, d1 := s_0.Date()
			y2, m2, d2 := holiday.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				isReserve = false
				break
			}
		}

		for d := s_0.Add(9 * time.Hour); d.Day() == s_0.Day() && d.Hour() <= 20; d = d.Add(30 * time.Minute) {
			if !isReserve {
				reserveTime = append(reserveTime, dto.ReserveTableSub{
					Time:      d,
					Users:     []string{},
					IsReserve: false,
				})
			} else {
				var tempUsers []string
				for _, schedule := range schedules {
					// 時刻が対象範囲の場合
					if d.After(schedule.Start.Add(-1*time.Minute)) && d.Before(schedule.End.Add(1*time.Minute)) {
						for _, user := range schedule.Users {
							tempUsers = append(tempUsers, user.HashKey)
						}
					}
				}

				var users []string
				for _, model := range models {
					users = append(users, model.UserHashKey)
				}

				seen := make(map[string]bool)
				unableUsers := []string{}

				for _, str := range tempUsers {
					if !seen[str] {
						seen[str] = true
						unableUsers = append(unableUsers, str)
					}
				}

				var ableUsers []string
				for _, a := range users {
					found := false
					for _, b := range unableUsers {
						if a == b {
							found = true
							break
						}
					}

					if !found {
						ableUsers = append(ableUsers, a)
					}
				}

				reserveTime = append(reserveTime, dto.ReserveTableSub{
					Time:      d,
					Users:     ableUsers,
					IsReserve: len(ableUsers) > int(applicant.NumOfInterview),
				})
			}
		}
	}

	// 面接予定取得
	var res entity.Schedule
	if applicant.ScheduleID > 0 {
		applicantSchedule, applicantScheduleErr := s.u.GetScheduleByPrimary(&ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: applicant.ScheduleID,
			},
		})
		if applicantScheduleErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		res.Start = applicantSchedule.Start
		res.End = applicantSchedule.End
	}

	return &response.ReserveTable{
		Dates:             times,
		Options:           reserveTime,
		Schedule:          res.Start,
		ScheduleHashKey:   res.HashKey,
		IsResume:          applicant.NumOfInterview == 1 || applicant.ResumeExtension == "",
		IsCurriculumVitae: applicant.NumOfInterview == 1 || applicant.CurriculumVitaeExtension == "",
	}, nil
}

// 書類アップロード(S3)
func (s *ApplicantService) S3Upload(req *request.FileUpload, fileHeader *multipart.FileHeader) *response.Error {
	// バリデーション
	if err := s.v.S3Upload(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者取得
	applicant, applicantErr := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if applicantErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// S3 Upload
	objName := req.NamePre + "_" + applicant.Name + "_" + applicant.Email + "." + req.Extension
	if err := s.a.S3Upload(objName, fileHeader); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 書類ダウンロード(S3)
func (s *ApplicantService) S3Download(req *request.FileDownload) ([]byte, *string, *response.Error) {
	// バリデーション
	if err := s.v.S3Download(req); err != nil {
		log.Printf("%v", err)
		return nil, nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ファイル名取得
	applicant, err := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		log.Printf("%v", err)
		return nil, nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// S3からダウンロード
	if applicant.ResumeExtension != "" && req.NamePre == "resume" {
		fileName := req.NamePre + "_" + applicant.Name + "_" + applicant.Email + "." + applicant.ResumeExtension
		file, err := s.a.S3Download(fileName)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &fileName, nil
	}
	if applicant.CurriculumVitaeExtension != "" && req.NamePre == "curriculum_vitae" {
		fileName := req.NamePre + "_" + applicant.Name + "_" + applicant.Email + "." + applicant.CurriculumVitaeExtension
		file, err := s.a.S3Download(fileName)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &fileName, nil
	}

	return nil, nil, nil
}

// 認証URL作成
func (s *ApplicantService) GetOauthURL(req *request.GetOauthURL) (*response.GetOauthURL, *response.Error) {
	// バリデーション
	if err := s.v.GetOauthURL(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	ctx := context.Background()
	if err := s.redis.Set(
		ctx,
		req.UserHashKey,
		static.REDIS_APPLICANT_HASH_KEY,
		&req.Applicant.HashKey,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	res, err := s.g.GetOauthURL()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return &response.GetOauthURL{
		Url: *res,
	}, nil
}

// 面接希望日登録
func (s *ApplicantService) InsertDesiredAt(req *request.InsertDesiredAt) *response.Error {
	// バリデーション
	if err := s.v.InsertDesiredAt(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者取得
	applicant, applicantErr := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.ApplicantHashKey,
		},
	})
	if applicantErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// イベント取得
	events, eventsErr := s.u.SelectEventAssociation(&ddl.TeamEvent{
		TeamID: applicant.TeamID,
	})
	if eventsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if applicant.ScheduleID == 0 {
		// ハッシュキー生成
		_, hash, hashErr := GenerateHash(1, 25)
		if hashErr != nil {
			log.Printf("%v", hashErr)
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// 予定登録
		id, scheduleErr := s.u.InsertSchedule(tx, &ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey:   static.PRE_SCHEDULE + "_" + *hash,
				CompanyID: applicant.CompanyID,
			},
			InterviewFlg: uint(static.USER_INTERVIEW),
			FreqID:       static.FREQ_NONE,
			Start:        req.DesiredAt,
			End:          req.DesiredAt.Add(time.Hour),
			Title:        req.Title,
			TeamID:       applicant.TeamID,
		})
		if scheduleErr != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// 応募者紐づけ登録
		if err := s.r.InsertApplicantScheduleAssociation(tx, &ddl.ApplicantScheduleAssociation{
			ApplicantID: applicant.ID,
			ScheduleID:  *id,
		}); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	} else {
		// 予定更新
		if err := s.u.UpdateScheduleByPrimary(tx, &ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: applicant.ScheduleID,
			},
			Start: req.DesiredAt,
			End:   req.DesiredAt.Add(time.Hour),
			Title: req.Title,
		}); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if applicant.NumOfInterview == 1 {
		// 履歴書登録
		if req.ResumeExtension != "" {
			// 削除
			if err := s.r.DeleteApplicantResumeAssociation(tx, &ddl.ApplicantResumeAssociation{
				ApplicantID: applicant.ID,
			}); err != nil {
				if err := s.d.TxRollback(tx); err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			// 登録
			if err := s.r.InsertApplicantResumeAssociation(tx, &ddl.ApplicantResumeAssociation{
				ApplicantID: applicant.ID,
				Extension:   req.ResumeExtension,
			}); err != nil {
				if err := s.d.TxRollback(tx); err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
		}

		// 職務経歴書登録
		if req.CurriculumVitaeExtension != "" {
			// 削除
			if err := s.r.DeleteApplicantCurriculumVitaeAssociation(tx, &ddl.ApplicantCurriculumVitaeAssociation{
				ApplicantID: applicant.ID,
			}); err != nil {
				if err := s.d.TxRollback(tx); err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			// 登録
			if err := s.r.InsertApplicantCurriculumVitaeAssociation(tx, &ddl.ApplicantCurriculumVitaeAssociation{
				ApplicantID: applicant.ID,
				Extension:   req.CurriculumVitaeExtension,
			}); err != nil {
				if err := s.d.TxRollback(tx); err != nil {
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

	// 応募者ステータス決定＆更新
	if req.ResumeExtension != "" && req.CurriculumVitaeExtension != "" {
		for _, event := range events {
			if event.EventID == static.STATUS_EVENT_SUBMIT_DOCUMENTS {
				if err := s.r.Update(tx, &ddl.Applicant{
					AbstractTransactionModel: ddl.AbstractTransactionModel{
						HashKey: applicant.HashKey,
					},
					Status: event.StatusID,
				}); err != nil {
					if err := s.d.TxRollback(tx); err != nil {
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
	} else {
		for _, event := range events {
			if event.EventID == static.STATUS_EVENT_DECIDE_SCHEDULE {
				if err := s.r.Update(tx, &ddl.Applicant{
					Status: event.StatusID,
				}); err != nil {
					if err := s.d.TxRollback(tx); err != nil {
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
	}

	if err := s.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// GoogleMeetUrl発行
func (s *ApplicantService) GetGoogleMeetUrl(req *request.GetGoogleMeetUrl) (*response.GetGoogleMeetUrl, *response.Error) {
	// バリデーション
	if err := s.v.GetGoogleMeetUrl(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者
	ctx := context.Background()
	applicantHashKey, applicantHashKeyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_APPLICANT_HASH_KEY)
	if applicantHashKeyErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者取得
	applicant, applicantErr := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *applicantHashKey,
		},
	})
	if applicantErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定取得
	schedule, scheduleErr := s.u.GetScheduleByPrimary(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: applicant.ScheduleID,
		},
	})
	if scheduleErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー取得
	user, userErr := s.u.Get(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.UserHashKey,
		},
	})
	if userErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// アクセストークン取得
	accessToken, accessTokenErr := s.g.GetAccessToken(&user.RefreshToken, &req.Code)
	if accessTokenErr != nil {
		log.Printf("%v", accessTokenErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Google Meet Url 発行
	googleMeetUrl, googleMeetUrlErr := s.g.GetGoogleMeetUrl(
		accessToken,
		user.Name,
		schedule.Start,
		schedule.End,
	)
	if googleMeetUrlErr != nil {
		log.Printf("%v", googleMeetUrlErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Google Meet Url 格納
	if err := s.r.InsertApplicantURLAssociation(tx, &ddl.ApplicantURLAssociation{
		ApplicantID: applicant.ID,
		URL:         *googleMeetUrl,
	}); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.GetGoogleMeetUrl{
		Url: *googleMeetUrl,
	}, nil
}

// 応募者ステータス変更
func (s *ApplicantService) UpdateStatus(req *request.UpdateStatus) *response.Error {
	// バリデーション
	if err := s.v.UpdateStatus(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}
	for _, row := range req.Association {
		if err := s.v.UpdateStatusSub(&row); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
			}
		}
	}
	for _, row := range req.Events {
		if err := s.v.UpdateStatusSub2(&row); err != nil {
			log.Printf("%v", err)
			return &response.Error{
				Status: http.StatusBadRequest,
			}
		}
	}

	// Redisから取得
	ctx := context.Background()
	team, teamErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*team, 10, 64)
	if teamIDErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.TeamID = teamID

	company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
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

	// 旧スタータス取得
	oldStatus, oldStatusErr := s.r.ListStatus(&ddl.SelectStatus{
		TeamID: req.TeamID,
	})
	if oldStatusErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// イベントハッシュキーからイベント取得
	var hashKeys []string
	for _, row := range req.Events {
		hashKeys = append(hashKeys, row.EventHash)
	}
	events, eventsErr := s.m.SelectSelectStatusEventByHashKeys(hashKeys)
	if eventsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for index := range req.Events {
		req.Events[index].EventID = events[index].ID
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 新ステータス登録
	var status []*ddl.SelectStatus
	for _, row := range req.Status {
		_, hashKey, hashErr := GenerateHash(1, 25)
		if hashErr != nil {
			log.Printf("%v", hashErr)
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		status = append(status, &ddl.SelectStatus{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *hashKey,
				CompanyID: companyID,
			},
			TeamID:     teamID,
			StatusName: row,
		})
	}
	ids, idsErr := s.u.InsertsSelectStatus(tx, status)
	if idsErr != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 各応募者のステータス更新
	for _, row := range req.Association {
		if err := s.r.UpdateSelectStatus(tx, &ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: row.BeforeHash,
			},
			Status: ids.List[row.AfterIndex].ID,
		}); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// イベントを一度全削除
	if err := s.u.DeleteEventAssociation(tx, &ddl.TeamEvent{
		TeamID: teamID,
	}); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// イベントに登録
	if len(req.Events) > 0 {
		var eventsDDL []*ddl.TeamEvent
		for _, row := range req.Events {
			eventsDDL = append(eventsDDL, &ddl.TeamEvent{
				TeamID:   teamID,
				EventID:  row.EventID,
				StatusID: ids.List[row.Status].ID,
			})
		}
		if err := s.u.InsertsEventAssociation(tx, eventsDDL); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// 面接毎イベントを一度全削除
	if err := s.u.DeleteEventEachInterviewAssociation(tx, &ddl.TeamEventEachInterview{
		TeamID: teamID,
	}); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接毎イベントに登録
	if len(req.EventsOfInterview) > 0 {
		var eventsDDL []*ddl.TeamEventEachInterview
		for _, row := range req.EventsOfInterview {
			eventsDDL = append(eventsDDL, &ddl.TeamEventEachInterview{
				TeamID:         teamID,
				NumOfInterview: row.Num,
				StatusID:       ids.List[row.Status].ID,
			})
		}
		if err := s.u.InsertsEventEachInterviewAssociation(tx, eventsDDL); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// 旧ステータス削除
	var oldStatusIds []uint64
	for _, row := range oldStatus {
		oldStatusIds = append(oldStatusIds, row.ID)
	}
	if err := s.r.DeleteStatusByPrimary(tx, &ddl.SelectStatus{
		TeamID: teamID,
	}, oldStatusIds); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 面接官割り振り
func (s *ApplicantService) AssignUser(req *request.AssignUser) *response.Error {
	// バリデーション
	if err := s.v.AssignUser(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者取得
	applicant, applicantErr := s.r.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if applicantErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム取得
	team, teamErr := s.u.GetTeamByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: applicant.TeamID,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接可能ユーザー取得
	ableTempUsers, ableTempUsersErr := s.u.GetAssignPossibleByNumOfInterview(&ddl.TeamAssignPossible{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if ableTempUsersErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var ableUsers []*ddl.User
	for _, user := range ableTempUsers {
		for _, user2 := range team.Users {
			if user.HashKey == user2.HashKey {
				ableUsers = append(ableUsers, user2)
				break
			}
		}
	}

	// 面接可能判定
	var users []*ddl.ApplicantUserAssociation
	for _, user := range req.HashKeys {
		for _, user2 := range ableUsers {
			if user == user2.HashKey {
				users = append(users, &ddl.ApplicantUserAssociation{
					ApplicantID: applicant.ID,
					UserID:      user2.ID,
				})
			}
		}
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官割り振り削除
	if err := s.u.DeleteUserAssociation(tx, &ddl.ApplicantUserAssociation{
		ApplicantID: applicant.ID,
	}); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官割り振り登録
	if err := s.u.InsertsUserAssociation(tx, users); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
