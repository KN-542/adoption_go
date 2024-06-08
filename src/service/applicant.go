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
	"strings"
	"time"
)

type IApplicantService interface {
	// 検索*
	Search(req *request.ApplicantSearch) (*response.ApplicantSearch, *response.Error)
	// サイト一覧取得
	GetSites() (*response.ApplicantSites, *response.Error)
	// 応募者ステータス一覧取得
	GetStatusList(req *request.ApplicantStatusList) (*response.ApplicantStatusList, *response.Error)
	// 応募者ダウンロード
	Download(req *request.ApplicantDownload) (*response.ApplicantDownload, *response.Error)
	// Google 認証URL作成*
	GetOauthURL(req *ddl.ApplicantAndUser) (*ddl.GetOauthURLResponse, *response.Error)
	// 応募者取得(1件)*
	Get(req *ddl.Applicant) (*ddl.Applicant, *response.Error)
	// 書類アップロード(S3)*
	S3Upload(req *ddl.FileUpload, fileHeader *multipart.FileHeader) *response.Error
	// 書類ダウンロード(S3)*
	S3Download(req *ddl.FileDownload) ([]byte, *string, *response.Error)
	// 面接希望日登録*
	InsertDesiredAt(req *ddl.ApplicantDesired) *response.Error
	// Google Meet Url 発行*
	GetGoogleMeetUrl(req *ddl.ApplicantAndUser) (*ddl.Applicant, *response.Error)
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
) IApplicantService {
	return &ApplicantService{r, u, m, a, g, redis, v, d}
}

// 検索
func (s *ApplicantService) Search(req *request.ApplicantSearch) (*response.ApplicantSearch, *response.Error) {
	// バリデーション
	if err := s.v.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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
	req.TeamID = teamID

	// 面接官取得
	var users []uint64
	for _, hash := range req.Users {
		user, err := s.u.Get(&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: hash,
			},
		})
		if err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		users = append(users, user.ID)
	}

	// 検索
	applicants, searchErr := s.r.Search(&dto.ApplicantSearch{
		ApplicantSearch: *req,
		UserIDs:         users,
	})
	if searchErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー紐づけ取得
	for _, row := range applicants {
		users, usersErr := s.u.GetUserAssociation(&ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: row.ID,
			},
		})
		if usersErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		var userHashKeys []string
		var userNames []string
		for _, row2 := range users {
			userHashKeys = append(userHashKeys, row2.HashKey)
			userNames = append(userNames, row2.Name)
		}
		row.Users = strings.Join(userHashKeys, ",")
		row.UserNames = strings.Join(userNames, ",")
		row.ID = 0
	}

	var res []entity.ApplicantSearch
	for _, row := range applicants {
		res = append(res, *row)
	}

	return &response.ApplicantSearch{
		List: res,
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
			Code:   static.CODE_BAD_REQUEST,
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
			Code:   static.CODE_BAD_REQUEST,
		}
	}
	for _, row := range req.Applicants {
		if err := s.v.DownloadSub(&row); err != nil {
			log.Printf("%v", err)
			return nil, &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_BAD_REQUEST,
			}
		}
	}

	// 重複チェック
	var request request.ApplicantDownload
	for _, row := range req.Applicants {
		count, err := s.r.CheckDuplByOuterId(&ddl.Applicant{
			OuterID: row.OuterID,
		})
		if err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		if *count == 0 {
			request.Applicants = append(request.Applicants, row)
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

	// サイトID取得
	site, siteErr := s.m.SelectSiteByHashKey(&ddl.Site{
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

		// 登録
		if err := s.r.Insert(tx, &ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey:   *hash,
				CompanyID: companyID,
			},
			OuterID: row.OuterID,
			SiteID:  site.ID,
			Status:  1,
			Name:    row.Name,
			Email:   row.Email,
			Tel:     row.Tel,
			Age:     uint(row.Age),
			TeamID:  teamID,
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

	return &response.ApplicantDownload{
		UpdateNum: len(request.Applicants),
	}, nil
}

// 認証URL作成
func (s *ApplicantService) GetOauthURL(req *ddl.ApplicantAndUser) (*ddl.GetOauthURLResponse, *response.Error) {
	// バリデーション
	if err := s.v.HashKeyValidate(&req.Applicant); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}
	if err := s.v.HashKeyValidate(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.UserHashKey,
		},
	}); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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
	return res, nil
}

// 応募者取得(1件)
func (s *ApplicantService) Get(req *ddl.Applicant) (*ddl.Applicant, *response.Error) {
	// バリデーション
	if err := s.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 応募者情報取得
	applicant, err := s.r.GetByHashKey(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return applicant, nil
}

// 書類アップロード(S3)
func (s *ApplicantService) S3Upload(req *ddl.FileUpload, fileHeader *multipart.FileHeader) *response.Error {
	// バリデーション
	if err := s.v.S3UploadValidator(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	fileName, err := s.redis.Get(ctx, req.HashKey, static.REDIS_S3_NAME)
	if err != nil {
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}

	// S3 Upload
	objName := req.NamePre + "_" + *fileName + "." + req.Extension
	if err := s.a.S3Upload(objName, fileHeader); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := s.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 書類登録状況更新
	if req.NamePre == "resume" {
		if err := s.r.UpdateDocument(tx, &ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.HashKey,
			},
			Resume:          objName,
			CurriculumVitae: "",
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
	if req.NamePre == "curriculum_vitae" {
		if err := s.r.UpdateDocument(tx, &ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.HashKey,
			},
			Resume:          "",
			CurriculumVitae: objName,
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

	if err := s.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 書類ダウンロード(S3)
func (s *ApplicantService) S3Download(req *ddl.FileDownload) ([]byte, *string, *response.Error) {
	// バリデーション
	if err := s.v.S3DownloadValidator(req); err != nil {
		log.Printf("%v", err)
		return nil, nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ファイル名取得
	applicant, err := s.r.GetByHashKey(&ddl.Applicant{
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
	if applicant.Resume != "" && req.NamePre == "resume" {
		file, err := s.a.S3Download(applicant.Resume)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &applicant.Resume, nil
	}
	if applicant.CurriculumVitae != "" && req.NamePre == "curriculum_vitae" {
		file, err := s.a.S3Download(applicant.CurriculumVitae)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &applicant.CurriculumVitae, nil
	}

	return nil, nil, nil
}

// 面接希望日登録
func (s *ApplicantService) InsertDesiredAt(req *ddl.ApplicantDesired) *response.Error {
	// バリデーション
	if err := s.v.InsertDesiredAtValidator(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// カレンダーID取得
	calendar, err := s.u.GetSchedule(&ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.CalendarHashKey,
		},
	})
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := s.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接希望日登録
	if err := s.r.UpdateDesiredAt(tx, &ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
		CalendarID: calendar.ID,
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

// Google Meet Url 発行
func (s *ApplicantService) GetGoogleMeetUrl(req *ddl.ApplicantAndUser) (*ddl.Applicant, *response.Error) {
	// バリデーション
	if err := s.v.HashKeyValidate(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.UserHashKey,
		},
	}); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	applicantHashKey, err := s.redis.Get(ctx, req.UserHashKey, static.REDIS_APPLICANT_HASH_KEY)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 面接情報取得
	schedule, err := s.r.GetDesiredAt(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *applicantHashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー取得
	user, err := s.u.Get(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.UserHashKey,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	accessToken, err := s.g.GetAccessToken(&user.RefreshToken, &req.Code)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Google Meet Url 発行
	googleMeetUrl, err := s.g.GetGoogleMeetUrl(
		accessToken,
		user.Name,
		schedule.Start,
		schedule.End,
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := s.d.TxStart()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Google Meet Url 格納
	if err := s.r.UpdateGoogleMeet(tx, &ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *applicantHashKey,
		},
		GoogleMeetURL: *googleMeetUrl,
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

	return &ddl.Applicant{
		GoogleMeetURL: *googleMeetUrl,
	}, nil
}
