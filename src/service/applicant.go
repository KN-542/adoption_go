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
	"fmt"
	"log"
	"math/rand"
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
	// 面接官割り振り可能判定
	CheckAssignableUser(req *request.CheckAssignableUser, domainFlg bool) (*response.CheckAssignableUser, *response.Error)
	// 種別登録
	CreateApplicantType(req *request.CreateApplicantType) *response.Error
	// 種別一覧
	ListApplicantType(req *request.ListApplicantType) (*response.ListApplicantType, *response.Error)
	// 応募者種別紐づけ登録
	CreateApplicantTypeAssociation(req *request.CreateApplicantTypeAssociation) *response.Error
	// ステータス更新
	UpdateSelectStatus(req *request.UpdateSelectStatus) *response.Error
}

type ApplicantService struct {
	r     repository.IApplicantRepository
	u     repository.IUserRepository
	t     repository.ITeamRepository
	s     repository.IScheduleRepository
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
	t repository.ITeamRepository,
	s repository.IScheduleRepository,
	m repository.IMasterRepository,
	a repository.IAWSRepository,
	g repository.IGoogleRepository,
	redis repository.IRedisRepository,
	v validator.IApplicantValidator,
	d repository.IDBRepository,
	o repository.IOuterIFRepository,
) IApplicantService {
	return &ApplicantService{r, u, t, s, m, a, g, redis, v, d, o}
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
		applicant.ID = 0
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

	// コミットID生成
	commitID, _, hashErr := GenerateHash(16, 28)
	if hashErr != nil {
		log.Printf("%v", hashErr)
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
					HashKey:   static.PRE_APPLICANT + "_" + *hash,
					CompanyID: companyID,
				},
				OuterID:        row.OuterID,
				SiteID:         site.ID,
				Status:         list[0].ID,
				Name:           row.Name,
				Email:          row.Email,
				Tel:            row.Tel,
				Age:            uint(row.Age),
				CommitID:       *commitID,
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

	// 面接設定取得
	setting, settingError := s.t.GetPerInterviewByNumOfInterview(&ddl.TeamPerInterview{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if settingError != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接参加可能者取得
	models, modelsErr := s.t.GetAssignPossibleSchedule(&ddl.TeamAssignPossible{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if modelsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接予定取得
	var res entity.Schedule
	if applicant.ScheduleID > 0 {
		applicantSchedule, applicantScheduleErr := s.s.GetByPrimary(&ddl.Schedule{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: applicant.ScheduleID,
			},
		})
		if applicantScheduleErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		res.ID = applicantSchedule.ID
		res.Start = applicantSchedule.Start
		res.End = applicantSchedule.End
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
					if res.ID == schedule.ID {
						continue
					}

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
					IsReserve: len(ableUsers) >= int(setting.UserMin),
				})
			}
		}
	}

	return &response.ReserveTable{
		Dates:             times,
		Options:           reserveTime,
		Schedule:          res.Start,
		ScheduleHashKey:   res.HashKey,
		IsResume:          setting.UserMin == 1 || applicant.ResumeExtension == "",
		IsCurriculumVitae: setting.UserMin == 1 || applicant.CurriculumVitaeExtension == "",
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

	// Google Meet URL取得
	associations, associationsErr := s.r.GetApplicantURLAssociation(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if associationsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(associations) == 1 {
		return &response.GetOauthURL{
			GoogleMeetURL: associations[0].URL,
		}, nil
	} else {
		// リフレッシュトークン紐づけ取得
		var refreshToken string
		tokenAssociations, tokenAssociationsErr := s.u.GetUserRefreshTokenAssociationByHashKey(&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.UserHashKey,
			},
		})
		if tokenAssociationsErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		if len(tokenAssociations) != 1 {
			refreshToken = ""
		} else {
			refreshToken = tokenAssociations[0].RefreshToken
		}

		ctx := context.Background()
		if err := s.redis.Set(
			ctx,
			req.UserHashKey,
			static.REDIS_APPLICANT_HASH_KEY,
			&req.HashKey,
			24*time.Hour,
		); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		if refreshToken == "" {
			res, resErr := s.g.GetOauthURL()
			if resErr != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.GetOauthURL{
				AuthURL: *res,
			}, nil
		} else {
			service, serviceErr := s.GetGoogleMeetUrl(&request.GetGoogleMeetUrl{
				Abstract:     req.Abstract,
				RefreshToken: refreshToken,
			})
			if serviceErr != nil {
				return nil, &response.Error{
					Status: serviceErr.Status,
				}
			}

			return &response.GetOauthURL{
				GoogleMeetURL: service.Url,
			}, nil
		}
	}
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

	// チーム取得
	team, teamErr := s.t.GetByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: applicant.TeamID,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接設定取得
	setting, settingError := s.t.GetPerInterviewByNumOfInterview(&ddl.TeamPerInterview{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if settingError != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官取得
	interviewers, interviewersError := s.r.GetUserAssociation(&ddl.ApplicantUserAssociation{
		ApplicantID: applicant.ID,
	})
	if interviewersError != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// イベント取得
	events, eventsErr := s.t.SelectEventAssociation(&ddl.TeamEvent{
		TeamID: applicant.TeamID,
	})
	if eventsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ～次面接イベント取得
	tEvent, tEventErr := s.t.GetEventEachInterviewAssociation(&ddl.TeamEventEachInterview{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if tEventErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官自動割り振りルール取得
	autoRules, autoRulesErr := s.t.GetAutoAssignRuleFind(&ddl.TeamAutoAssignRule{
		TeamID: applicant.TeamID,
	})
	if autoRulesErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(autoRules) > 1 {
		log.Printf("duplicate auto rule.")
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 優先順位取得
	priorities, prioritiesErr := s.t.GetAssignPriorityOnly(&ddl.TeamAssignPriority{
		TeamID: team.ID,
	})
	if prioritiesErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー取得_予定数順
	usersDesc, usersDescErr := s.u.GetUsersSortedByScheduleCount(&ddl.Schedule{
		TeamID: team.ID,
	})
	if usersDescErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官割り振り可能判定
	var userHashKeys []string
	for _, user := range team.Users {
		userHashKeys = append(userHashKeys, user.HashKey)
	}
	service, serviceErr := s.CheckAssignableUser(&request.CheckAssignableUser{
		Start:    req.DesiredAt,
		HashKeys: userHashKeys,
	}, true)
	if serviceErr != nil {
		return &response.Error{
			Status: serviceErr.Status,
		}
	}

	// 面接参加可能者取得
	models, modelsErr := s.t.GetAssignPossibleSchedule(&ddl.TeamAssignPossible{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if modelsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var ableUsers []*response.CheckAssignableUserSub
	for _, s := range service.List {
		flg := false
		for _, m := range models {
			if m.UserID == s.User.ID {
				flg = true
				break
			}
		}

		if flg {
			userCopy := s // コピーを作成
			ableUsers = append(ableUsers, &userCopy)
		}
	}

	for _, s := range ableUsers {
		fmt.Print(*s)
		fmt.Print("\n")
		fmt.Print(s.User)
		fmt.Print("\n")
		fmt.Print(s.User.ID)
		fmt.Print("\n")
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var scheduleID uint64
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
		id, scheduleErr := s.s.Insert(tx, &ddl.Schedule{
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
		scheduleID = *id

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
		if err := s.s.UpdateByPrimary(tx, &ddl.Schedule{
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
	if tEvent.NumOfInterview > 1 {
		if err := s.r.Update(tx, &ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: applicant.HashKey,
			},
			Status: tEvent.StatusID,
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
		eventStatus := uint(0)
		if req.ResumeExtension != "" && req.CurriculumVitaeExtension != "" {
			eventStatus = uint(static.STATUS_EVENT_SUBMIT_DOCUMENTS)
		} else {
			eventStatus = uint(static.STATUS_EVENT_DECIDE_SCHEDULE)
		}
		for _, event := range events {
			if event.EventID == eventStatus {
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
	}

	// 面接官割り振り
	var users []entity.User
	for _, s := range ableUsers {
		fmt.Print(*s)
		fmt.Print("\n")
		fmt.Print(s.DuplFlg)
		fmt.Print("\n")
		if s.DuplFlg == static.DUPLICATION_SAFE {
			users = append(users, s.User)
		}
	}
	if len(users) < int(setting.UserMin) {
		log.Printf("Not assignable.")
		return &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_APPLICANT_CANNOT_ASSIGN_USER,
		}
	}

	if len(interviewers) == 0 && applicant.ScheduleID == 0 && len(autoRules) == 1 && team.RuleID == static.ASSIGN_RULE_AUTO {
		autoRule := autoRules[0]
		var applicantUsers []*ddl.ApplicantUserAssociation

		// ランダム
		if autoRule.RuleID == static.AUTO_ASSIGN_RULE_RANDOM {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))

			rng.Shuffle(len(users), func(i, j int) {
				users[i], users[j] = users[j], users[i]
			})

			for _, user := range users[:int(setting.UserMin)] {
				applicantUsers = append(applicantUsers, &ddl.ApplicantUserAssociation{
					ApplicantID: applicant.ID,
					UserID:      user.ID,
				})
			}
		}

		// 定められた優先順位に従う
		if autoRule.RuleID == static.AUTO_ASSIGN_RULE_ASC && len(priorities) > 0 {
			var filterPriorities []entity.TeamAssignPriorityOnly
			for _, row := range priorities {
				flg := false
				for _, user := range users {
					if user.ID == row.UserID {
						flg = true
						break
					}
				}

				if flg {
					filterPriorities = append(filterPriorities, *row)
				}
			}

			count := 0
			for _, row := range filterPriorities {
				count++
				applicantUsers = append(applicantUsers, &ddl.ApplicantUserAssociation{
					ApplicantID: applicant.ID,
					UserID:      row.UserID,
				})
				if count == int(setting.UserMin) {
					break
				}
			}
		}

		// 予定の少ない順
		if autoRule.RuleID == static.AUTO_ASSIGN_RULE_DESC_SCHEDULE {
			var filterList []entity.User
			for _, row := range usersDesc {
				flg := false
				for _, user := range users {
					if user.ID == row.ID {
						flg = true
						break
					}
				}

				if flg {
					filterList = append(filterList, row)
				}
			}

			count := 0
			for _, row := range filterList {
				count++
				applicantUsers = append(applicantUsers, &ddl.ApplicantUserAssociation{
					ApplicantID: applicant.ID,
					UserID:      row.ID,
				})
				if count == int(setting.UserMin) {
					break
				}
			}
		}

		if len(applicantUsers) > 0 {
			// 応募者ユーザー紐づけ登録
			if err := s.r.InsertsUserAssociation(tx, applicantUsers); err != nil {
				if err := s.d.TxRollback(tx); err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			// 予定ユーザー紐づけ登録
			var list []*ddl.ScheduleAssociation
			for _, row := range applicantUsers {
				list = append(list, &ddl.ScheduleAssociation{
					ScheduleID: scheduleID,
					UserID:     row.UserID,
				})
			}
			if err := s.s.InsertsScheduleAssociation(tx, list); err != nil {
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
	schedule, scheduleErr := s.s.GetByPrimary(&ddl.Schedule{
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

	// トークン取得
	token, tokenErr := s.g.GetAccessToken(&req.RefreshToken, &req.Code)
	if tokenErr != nil {
		log.Printf("%v", tokenErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Google Meet Url 発行
	googleMeetUrl, googleMeetUrlErr := s.g.GetGoogleMeetUrl(
		token,
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

	if req.RefreshToken == "" && token.RefreshToken != "" {
		// リフレッシュトークン格納
		if err := s.u.InsertUserRefreshTokenAssociation(tx, &ddl.UserRefreshTokenAssociation{
			UserID:       user.ID,
			RefreshToken: token.RefreshToken,
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
	ids, idsErr := s.t.InsertsSelectStatus(tx, status)
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
	if err := s.t.DeleteEventAssociation(tx, &ddl.TeamEvent{
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
		if err := s.t.InsertsEventAssociation(tx, eventsDDL); err != nil {
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
	if err := s.t.DeleteEventEachInterviewAssociation(tx, &ddl.TeamEventEachInterview{
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
		if err := s.t.InsertsEventEachInterviewAssociation(tx, eventsDDL); err != nil {
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
	if err := s.r.DeleteStatusByPrimaryAndTeam(tx, &ddl.SelectStatus{
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

	if applicant.ScheduleID == 0 {
		log.Printf("Schedule does not exist.")
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_APPLICANT_SCHEDULE_DOES_NOT_EXIST,
		}
	}

	// 予定取得
	schedule, scheduleErr := s.s.GetByPrimary(&ddl.Schedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: applicant.ScheduleID,
		},
	})
	if scheduleErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム取得
	team, teamErr := s.t.GetByPrimary(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: applicant.TeamID,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接設定取得
	setting, settingError := s.t.GetPerInterviewByNumOfInterview(&ddl.TeamPerInterview{
		TeamID:         applicant.TeamID,
		NumOfInterview: applicant.NumOfInterview,
	})
	if settingError != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接可能ユーザー取得
	ableTempUsers, ableTempUsersErr := s.t.GetAssignPossibleByNumOfInterview(&ddl.TeamAssignPossible{
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

	// 面接可能判定1
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

	// 面接可能判定2
	service, serviceErr := s.CheckAssignableUser(&request.CheckAssignableUser{
		Start:    schedule.Start,
		HashKeys: req.HashKeys,
	}, true)
	if serviceErr != nil {
		return &response.Error{
			Status: serviceErr.Status,
		}
	}

	var users2 []*ddl.ApplicantUserAssociation
	for _, s := range service.List {
		if s.DuplFlg != static.DUPLICATION_SAFE {
			continue
		}
		for _, user := range users {
			if s.User.ID == user.UserID {
				users2 = append(users2, user)
			}
		}
	}

	if len(users2) < int(setting.UserMin) {
		log.Printf("Shortage user min.")
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_APPLICANT_SHORTAGE_USER_MIN,
		}
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接官割り振り削除
	if err := s.r.DeleteUserAssociation(tx, &ddl.ApplicantUserAssociation{
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
	if err := s.r.InsertsUserAssociation(tx, users2); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 予定紐づけ削除
	if err := s.s.DeleteScheduleAssociation(tx, &ddl.ScheduleAssociation{
		ScheduleID: applicant.ScheduleID,
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
	var schedules []*ddl.ScheduleAssociation
	for _, row := range users2 {
		schedules = append(schedules, &ddl.ScheduleAssociation{
			ScheduleID: applicant.ScheduleID,
			UserID:     row.UserID,
		})
	}
	if err := s.s.InsertsScheduleAssociation(tx, schedules); err != nil {
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

// 面接官割り振り可能判定
func (s *ApplicantService) CheckAssignableUser(req *request.CheckAssignableUser, domainFlg bool) (*response.CheckAssignableUser, *response.Error) {
	// バリデーション
	if err := s.v.CheckAssignableUser(req); err != nil {
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

	// ユーザー取得
	users, usersErr := s.u.GetByHashKeys(req.HashKeys)
	if usersErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var list []response.CheckAssignableUserSub
	for _, user := range users {
		// ユーザー単位予定取得
		models, modelsErr := s.s.GetScheduleByUser(&dto.GetScheduleByUser{
			ScheduleAssociation: ddl.ScheduleAssociation{
				UserID: user.ID,
			},
			RemoveScheduleHashKeys: req.RemoveScheduleHashKeys,
		})
		if modelsErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// 整形
		var schedules []entity.Schedule
		var schedulesJST []entity.Schedule

		start := req.Start.In(jst)
		end := start.Add(1 * time.Hour)

		for _, model := range models {
			model.Start = model.Start.In(jst)
			model.End = model.End.In(jst)
			schedulesJST = append(schedulesJST, entity.Schedule{
				Schedule: model.Schedule,
			})
		}

		for _, scheduleJST := range schedulesJST {
			// なしの場合 (変更なし)
			if scheduleJST.FreqID == uint(static.FREQ_NONE) {
				if (end.After(scheduleJST.Start) || end.Equal(scheduleJST.Start)) &&
					(start.Before(scheduleJST.End) || start.Equal(scheduleJST.End)) {
					schedules = append(schedules, scheduleJST)
				}
			}

			// 毎日の場合 (変更なし)
			if scheduleJST.FreqID == uint(static.FREQ_DAILY) {
				scheduleStart := time.Date(start.Year(), start.Month(), start.Day(), scheduleJST.Start.Hour(), scheduleJST.Start.Minute(), scheduleJST.Start.Second(), 0, jst)
				scheduleEnd := time.Date(start.Year(), start.Month(), start.Day(), scheduleJST.End.Hour(), scheduleJST.End.Minute(), scheduleJST.End.Second(), 0, jst)

				if scheduleEnd.Before(scheduleStart) {
					scheduleEnd = scheduleEnd.Add(24 * time.Hour)
				}

				if (end.After(scheduleStart) || end.Equal(scheduleStart)) &&
					(start.Before(scheduleEnd) || start.Equal(scheduleEnd)) {
					schedules = append(schedules, scheduleJST)
				}
			}

			// 毎週の場合 (変更なし)
			if scheduleJST.FreqID == uint(static.FREQ_WEEKLY) {
				if start.Weekday() == scheduleJST.Start.Weekday() {
					scheduleStart := time.Date(start.Year(), start.Month(), start.Day(), scheduleJST.Start.Hour(), scheduleJST.Start.Minute(), scheduleJST.Start.Second(), 0, jst)
					scheduleEnd := time.Date(start.Year(), start.Month(), start.Day(), scheduleJST.End.Hour(), scheduleJST.End.Minute(), scheduleJST.End.Second(), 0, jst)

					if scheduleEnd.Before(scheduleStart) {
						scheduleEnd = scheduleEnd.Add(24 * time.Hour)
					}

					if (end.After(scheduleStart) || end.Equal(scheduleStart)) &&
						(start.Before(scheduleEnd) || start.Equal(scheduleEnd)) {
						schedules = append(schedules, scheduleJST)
					}
				}
			}

			// 毎月の場合 (存在しない日付の場合は何もしない)
			if scheduleJST.FreqID == uint(static.FREQ_MONTHLY) {
				scheduleDay := scheduleJST.Start.Day()

				if start.Day() == scheduleDay {
					scheduleStart := time.Date(start.Year(), start.Month(), scheduleDay, scheduleJST.Start.Hour(), scheduleJST.Start.Minute(), scheduleJST.Start.Second(), 0, jst)
					scheduleEnd := time.Date(start.Year(), start.Month(), scheduleDay, scheduleJST.End.Hour(), scheduleJST.End.Minute(), scheduleJST.End.Second(), 0, jst)

					// 存在しない日付の場合は無視
					if scheduleStart.Month() == start.Month() {
						if scheduleEnd.Before(scheduleStart) {
							scheduleEnd = scheduleEnd.Add(24 * time.Hour)
						}

						if (end.After(scheduleStart) || end.Equal(scheduleStart)) &&
							(start.Before(scheduleEnd) || start.Equal(scheduleEnd)) {
							schedules = append(schedules, scheduleJST)
						}
					}
				}
			}

			// 毎年の場合 (修正: 2/29のみ特別処理、その他の存在しない日付は考慮しない)
			if scheduleJST.FreqID == uint(static.FREQ_YEARLY) {
				scheduleMonth := scheduleJST.Start.Month()
				scheduleDay := scheduleJST.Start.Day()

				// 2月29日の特別処理
				if scheduleMonth == time.February && scheduleDay == 29 {
					if !isLeapYear(start.Year()) {
						scheduleDay = 1
						scheduleMonth = time.March
					}
				}

				if start.Month() == scheduleMonth && start.Day() == scheduleDay {
					scheduleStart := time.Date(start.Year(), scheduleMonth, scheduleDay, scheduleJST.Start.Hour(), scheduleJST.Start.Minute(), scheduleJST.Start.Second(), 0, jst)
					scheduleEnd := time.Date(start.Year(), scheduleMonth, scheduleDay, scheduleJST.End.Hour(), scheduleJST.End.Minute(), scheduleJST.End.Second(), 0, jst)

					if scheduleEnd.Before(scheduleStart) {
						scheduleEnd = scheduleEnd.Add(24 * time.Hour)
					}

					if (end.After(scheduleStart) || end.Equal(scheduleStart)) &&
						(start.Before(scheduleEnd) || start.Equal(scheduleEnd)) {
						schedules = append(schedules, scheduleJST)
					}
				}
			}
		}

		flg := false
		id := uint64(0)
		if domainFlg {
			id = user.ID
		}
		resUser := entity.User{
			User: ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					ID:      id,
					HashKey: user.HashKey,
				},
				Name:  user.Name,
				Email: user.Email,
			},
		}

		// 重複リスト作成
		for _, s := range schedules {
			if s.InterviewFlg == uint(static.USER_INTERVIEW) {
				list = append(list, response.CheckAssignableUserSub{
					User:    resUser,
					DuplFlg: static.DUPLICATION_OUT,
				})
			}
			flg = true
			break
		}

		if !flg && len(schedules) > 0 {
			list = append(list, response.CheckAssignableUserSub{
				User:    resUser,
				DuplFlg: static.DUPLICATION_WARNING,
			})
		}

		if !flg && len(schedules) == 0 {
			list = append(list, response.CheckAssignableUserSub{
				User:    resUser,
				DuplFlg: static.DUPLICATION_SAFE,
			})
		}
	}

	return &response.CheckAssignableUser{
		List: list,
	}, nil
}

// 種別登録
func (s *ApplicantService) CreateApplicantType(req *request.CreateApplicantType) *response.Error {
	// バリデーション
	if err := s.v.CreateApplicantType(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// チーム、企業取得
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

	// 書類提出ルール取得
	documentRule, documentRuleErr := s.m.SelectDocumentRuleByHash(&ddl.DocumentRule{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: req.RuleHash,
		},
	})
	if documentRuleErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 種別取得
	occupation, occupationErr := s.m.SelectOccupationByHash(&ddl.Occupation{
		AbstractMasterModel: ddl.AbstractMasterModel{
			HashKey: req.OccupationHash,
		},
	})
	if occupationErr != nil {
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

	// ハッシュキー生成
	_, hash, hashErr := GenerateHash(1, 25)
	if hashErr != nil {
		log.Printf("%v", hashErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	if err := s.r.InsertType(tx, &ddl.ApplicantType{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   static.PRE_APPLICANT_TYPE + "_" + *hash,
			CompanyID: companyID,
		},
		TeamID:       teamID,
		RuleID:       documentRule.ID,
		OccupationID: occupation.ID,
		Name:         req.Name,
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

	if err := s.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 種別一覧
func (s *ApplicantService) ListApplicantType(req *request.ListApplicantType) (*response.ListApplicantType, *response.Error) {
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

	// 種別一覧
	res, err := s.r.ListType(&ddl.ApplicantType{
		TeamID: teamID,
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.ListApplicantType{
		List: res,
	}, nil
}

// 応募者種別紐づけ登録
func (s *ApplicantService) CreateApplicantTypeAssociation(req *request.CreateApplicantTypeAssociation) *response.Error {
	// バリデーション
	if err := s.v.CreateApplicantTypeAssociation(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 種別取得
	appType, typeErr := s.r.GetType(&ddl.ApplicantType{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.TypeHash,
		},
	})
	if typeErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者ID取得
	ids, idsErr := s.r.GetIDs(req.Applicants)
	if idsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var associations []*ddl.ApplicantTypeAssociation
	for _, id := range ids {
		associations = append(associations, &ddl.ApplicantTypeAssociation{
			TypeID:      appType.ID,
			ApplicantID: id,
		})
	}

	tx, txErr := s.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := s.r.DeleteTypeAssociation(tx, ids); err != nil {
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
	if err := s.r.InsertsTypeAssociation(tx, associations); err != nil {
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

// ステータス更新
func (s *ApplicantService) UpdateSelectStatus(req *request.UpdateSelectStatus) *response.Error {
	// バリデーション
	if err := s.v.UpdateSelectStatus(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ステータス取得
	status, statusErr := s.r.GetStatus(&ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.StatusHash,
		},
	})
	if statusErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者ID取得
	ids, idsErr := s.r.GetIDs(req.Applicants)
	if idsErr != nil {
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

	// 更新
	if err := s.r.UpdatesByPrimary(tx, &ddl.Applicant{
		Status: status.ID,
	}, ids); err != nil {
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
