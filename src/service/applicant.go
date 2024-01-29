package service

import (
	"api/resources/static"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"
)

type IApplicantService interface {
	// 認証URL作成
	GetOauthURL() (*model.GetOauthURLResponse, *model.ErrorResponse)
	// シート取得
	GetSheets(search model.ApplicantSearch) (*[]model.ApplicantResponse, *model.ErrorResponse)
	// 応募者ダウンロード
	Download(d *model.ApplicantsDownload) *model.ErrorResponse
	// 検索
	Search(req *model.ApplicantSearchRequest) (*model.ApplicantsDownloadResponse, *model.ErrorResponse)
	// 書類アップロード(S3)
	S3Upload(req *model.FileUpload, fileHeader *multipart.FileHeader) *model.ErrorResponse
	// 書類ダウンロード(S3)
	S3Download(req *model.FileDownload) ([]byte, *string, *model.ErrorResponse)
	// 面接希望日登録
	InsertDesiredAt(req *model.ApplicantDesired) *model.ErrorResponse
	// 応募者ステータス一覧取得
	GetApplicantStatus() (*model.ApplicantStatusList, *model.ErrorResponse)
	// サイト一覧取得
	GetSites() (*model.Sites, *model.ErrorResponse)
	// Google Meet Url 発行
	GetGoogleMeetUrl(req *model.ApplicantAndUser) (*model.Applicant, *model.ErrorResponse)
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
func (s *ApplicantService) GetOauthURL() (*model.GetOauthURLResponse, *model.ErrorResponse) {
	res, err := s.g.GetOauthURL()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return res, nil
}

// シート取得
func (s *ApplicantService) GetSheets(search model.ApplicantSearch) (*[]model.ApplicantResponse, *model.ErrorResponse) {
	refreshToken, _ := s.r.GetRefreshToken()

	accessToken, err := s.r.GetAccessToken(refreshToken, &search.Code)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	res, err := s.r.GetSheets(search, accessToken)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
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
func (s *ApplicantService) Download(d *model.ApplicantsDownload) *model.ErrorResponse {
	// STEP1 サイトIDチェック
	_, err := s.m.SelectSiteByPrimaryKey(d.Site)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// STEP2 登録
	if d.Site == int(enum.RECRUIT) {
		for _, values := range d.Values {
			_, size := utf8.DecodeLastRuneInString(values[enum.RECRUIT_AGE])
			age, err := strconv.ParseInt(
				values[enum.RECRUIT_AGE][:len(values[enum.RECRUIT_AGE])-size],
				10,
				64,
			)
			if err != nil {
				age = -1
			}

			// ハッシュキー生成
			_, hashKey, err := GenerateHash(1, 25)
			if err != nil {
				log.Printf("%v", err)
				return &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}

			m := model.Applicant{
				ID:        values[enum.RECRUIT_ID],
				HashKey:   *hashKey,
				SiteID:    uint(enum.RECRUIT),
				Status:    uint(enum.SCHEDULE_UNANSWERED),
				Name:      values[enum.RECRUIT_NAME],
				Email:     values[enum.RECRUIT_EMAIL],
				Tel:       values[enum.RECRUIT_TEL],
				Age:       int(age),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// STEP2-1 重複チェック
			count, err := s.r.CountByPrimaryKey(&m.ID)
			if err != nil {
				log.Printf("%v", err)
				return &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}
			if *count == int64(0) {
				tx, err := s.d.TxStart()
				if err != nil {
					return &model.ErrorResponse{
						Status: http.StatusInternalServerError,
					}
				}

				// STEP2-2 登録
				if err := s.r.Insert(tx, &m); err != nil {
					if err := s.d.TxRollback(tx); err != nil {
						return &model.ErrorResponse{
							Status: http.StatusInternalServerError,
						}
					}
					return &model.ErrorResponse{
						Status: http.StatusInternalServerError,
					}
				}

				if err := s.d.TxCommit(tx); err != nil {
					return &model.ErrorResponse{
						Status: http.StatusInternalServerError,
					}
				}
			}
		}
	}

	return nil
}

// 検索
func (s *ApplicantService) Search(req *model.ApplicantSearchRequest) (*model.ApplicantsDownloadResponse, *model.ErrorResponse) {
	// バリデーション
	if err := s.v.SearchValidator(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	applicants, err := s.r.Search(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.ApplicantsDownloadResponse{
		Applicants: applicants,
	}, nil
}

// 書類アップロード(S3)
func (s *ApplicantService) S3Upload(req *model.FileUpload, fileHeader *multipart.FileHeader) *model.ErrorResponse {
	// バリデーション
	if err := s.v.S3UploadValidator(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	fileName, err := s.redis.Get(ctx, req.HashKey, static.REDIS_S3_NAME)
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}

	// S3 Upload
	objName := req.NamePre + "_" + *fileName + "." + req.Extension
	if err := s.a.S3Upload(objName, fileHeader); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := s.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 書類登録状況更新
	if req.NamePre == "resume" {
		if err := s.r.UpdateDocument(tx, &model.Applicant{
			HashKey:         req.HashKey,
			Resume:          objName,
			CurriculumVitae: "",
		}); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}
	if req.NamePre == "curriculum_vitae" {
		if err := s.r.UpdateDocument(tx, &model.Applicant{
			HashKey:         req.HashKey,
			Resume:          "",
			CurriculumVitae: objName,
		}); err != nil {
			if err := s.d.TxRollback(tx); err != nil {
				return &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 書類ダウンロード(S3)
func (s *ApplicantService) S3Download(req *model.FileDownload) ([]byte, *string, *model.ErrorResponse) {
	// バリデーション
	if err := s.v.S3DownloadValidator(req); err != nil {
		log.Printf("%v", err)
		return nil, nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ファイル名取得
	applicant, err := s.r.GetByHashKey(&model.Applicant{
		HashKey: req.HashKey,
	})
	if err != nil {
		log.Printf("%v", err)
		return nil, nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// S3からダウンロード
	if applicant.Resume != "" && req.NamePre == "resume" {
		file, err := s.a.S3Download(applicant.Resume)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &applicant.Resume, nil
	}
	if applicant.CurriculumVitae != "" && req.NamePre == "curriculum_vitae" {
		file, err := s.a.S3Download(applicant.CurriculumVitae)
		if err != nil {
			log.Printf("%v", err)
			return nil, nil, &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return file, &applicant.CurriculumVitae, nil
	}

	return nil, nil, nil
}

// 面接希望日登録
func (s *ApplicantService) InsertDesiredAt(req *model.ApplicantDesired) *model.ErrorResponse {
	// バリデーション
	if err := s.v.InsertDesiredAtValidator(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	tx, err := s.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.r.UpdateDesiredAt(tx, &model.Applicant{
		HashKey:   req.HashKey,
		DesiredAt: req.DesiredAt,
	}); err != nil {
		if err := s.d.TxRollback(tx); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.d.TxCommit(tx); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 応募者ステータス一覧取得
func (s *ApplicantService) GetApplicantStatus() (*model.ApplicantStatusList, *model.ErrorResponse) {
	applicantStatus, err := s.m.SelectApplicantStatus()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.ApplicantStatusList{List: *applicantStatus}, nil
}

// サイト一覧取得
func (s *ApplicantService) GetSites() (*model.Sites, *model.ErrorResponse) {
	sites, err := s.m.SelectSite()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.Sites{List: *sites}, nil
}

// Google Meet Url 発行
func (s *ApplicantService) GetGoogleMeetUrl(req *model.ApplicantAndUser) (*model.Applicant, *model.ErrorResponse) {
	// バリデーション
	if err := s.v.HashKeyValidate(&req.Applicant); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	userHashKey, err := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	user, err := s.u.Get(&model.User{
		HashKey: *userHashKey,
	})
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	accessToken, err := s.g.GetAccessToken(&user.RefreshToken, &req.Code)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	googleMeetUrl, err := s.g.GetGoogleMeetUrl(
		accessToken,
		user.Name,
		req.Applicant.DesiredAt,
		req.Applicant.DesiredAt.Add(time.Hour),
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.Applicant{
		GoogleMeetURL: *googleMeetUrl,
	}, nil
}
