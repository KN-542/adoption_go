package service

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/enum"
	"api/src/model/response"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type IApplicantService interface {
	// Google 認証URL作成
	GetOauthURL(req *ddl.ApplicantAndUser) (*ddl.GetOauthURLResponse, *response.Error)
	// 応募者ダウンロード
	Download(d *ddl.ApplicantsDownload) *response.Error
	// 応募者取得(1件)
	Get(req *ddl.Applicant) (*ddl.Applicant, *response.Error)
	// 検索
	Search(req *ddl.ApplicantSearchRequest) (*ddl.ApplicantsDownloadResponse, *response.Error)
	// 書類アップロード(S3)
	S3Upload(req *ddl.FileUpload, fileHeader *multipart.FileHeader) *response.Error
	// 書類ダウンロード(S3)
	S3Download(req *ddl.FileDownload) ([]byte, *string, *response.Error)
	// 面接希望日登録
	InsertDesiredAt(req *ddl.ApplicantDesired) *response.Error
	// 応募者ステータス一覧取得
	GetApplicantStatus() (*ddl.ApplicantStatusList, *response.Error)
	// サイト一覧取得
	GetSites() (*ddl.Sites, *response.Error)
	// Google Meet Url 発行
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

/*
	txt、csvダウンロード用
*/
// 応募者ダウンロード
func (s *ApplicantService) Download(d *ddl.ApplicantsDownload) *response.Error {
	// STEP1 サイトIDチェック
	_, err := s.m.SelectSiteByPrimaryKey(d.Site)
	if err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// STEP2 登録
	if d.Site == int(enum.RECRUIT) {
		for _, values := range d.Values {
			_, size := utf8.DecodeLastRuneInString(values[enum.RECRUIT_AGE])
			_, err := strconv.ParseInt(
				values[enum.RECRUIT_AGE][:len(values[enum.RECRUIT_AGE])-size],
				10,
				64,
			)
			if err != nil {
				// TODO
			}

			// ハッシュキー生成
			_, hashKey, err := GenerateHash(1, 25)
			if err != nil {
				log.Printf("%v", err)
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			m := ddl.Applicant{
				OuterID: values[enum.RECRUIT_ID],
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey:   *hashKey,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				SiteID: uint(enum.RECRUIT),
				Status: uint(enum.SCHEDULE_UNANSWERED),
				Name:   values[enum.RECRUIT_NAME],
				Email:  values[enum.RECRUIT_EMAIL],
			}

			// STEP2-1 重複チェック
			count, err := s.r.CountByPrimaryKey(&m.OuterID)
			if err != nil {
				log.Printf("%v", err)
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			if *count == int64(0) {
				tx, err := s.d.TxStart()
				if err != nil {
					return &response.Error{
						Status: http.StatusInternalServerError,
					}
				}

				// STEP2-2 登録
				if err := s.r.Insert(tx, &m); err != nil {
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
			}
		}
	}

	return nil
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

// 検索
func (s *ApplicantService) Search(req *ddl.ApplicantSearchRequest) (*ddl.ApplicantsDownloadResponse, *response.Error) {
	// バリデーション
	if err := s.v.SearchValidator(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	applicants, err := s.r.Search(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	for i, row := range applicants {
		var list []string
		user, err := s.u.GetUserScheduleAssociationByScheduleID(
			&ddl.UserScheduleAssociation{
				UserScheduleID: row.CalendarID,
			},
		)
		if err != nil {
			log.Printf("%v", err)
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		for _, r := range user {
			list = append(list, r.Name)
		}
		applicants[i].UserNames = strings.Join(list, ",")
	}

	return &ddl.ApplicantsDownloadResponse{
		Applicants: applicants,
	}, nil
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
		CalendarID: uint(calendar.ID),
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

// 応募者ステータス一覧取得
func (s *ApplicantService) GetApplicantStatus() (*ddl.ApplicantStatusList, *response.Error) {
	applicantStatus, err := s.m.SelectApplicantStatus()
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.ApplicantStatusList{List: applicantStatus}, nil
}

// サイト一覧取得
func (s *ApplicantService) GetSites() (*ddl.Sites, *response.Error) {
	sites, err := s.m.SelectSite()
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.Sites{List: sites}, nil
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
