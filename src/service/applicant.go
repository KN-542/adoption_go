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
	"strings"
	"time"
	"unicode/utf8"
)

type IApplicantService interface {
	/*
		OAuth2.0用(削除予定)
	*/
	// 認証URL作成
	GetOauthURL() (*model.GetOauthURLResponse, *model.ErrorResponse)
	// シート取得
	GetSheets(search model.ApplicantSearch) (*[]model.ApplicantResponse, *model.ErrorResponse)
	/*
		txt、csvダウンロード用
	*/
	// 応募者ダウンロード
	Download(d *model.ApplicantsDownload) *model.ErrorResponse
	// 検索
	Search() (*model.ApplicantsDownloadResponse, *model.ErrorResponse)
	// 書類アップロード(S3)
	S3Upload(req *model.FileUpload, fileHeader *multipart.FileHeader) *model.ErrorResponse
	// 面接希望日登録
	InsertDesiredAt(req *model.ApplicantDesired) *model.ErrorResponse
}

type ApplicantService struct {
	r     repository.IApplicantRepository
	m     repository.IMasterRepository
	a     repository.IAWSRepository
	redis repository.IRedisRepository
	v     validator.IApplicantValidator
}

func NewApplicantService(
	r repository.IApplicantRepository,
	m repository.IMasterRepository,
	a repository.IAWSRepository,
	redis repository.IRedisRepository,
	v validator.IApplicantValidator,
) IApplicantService {
	return &ApplicantService{r, m, a, redis, v}
}

// 認証URL作成
func (s *ApplicantService) GetOauthURL() (*model.GetOauthURLResponse, *model.ErrorResponse) {
	res, err := s.r.GetOauthURL()
	if err != nil {
		log.Printf("%v", err)
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
			_, hashKey, err := generateRandomStr(1, 25)
			if err != nil {
				log.Printf("%v", err)
				return &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}

			m := model.Applicant{
				ID:        values[enum.RECRUIT_ID],
				HashKey:   *hashKey,
				SiteID:    int(enum.RECRUIT),
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
				// STEP2-2 登録
				if err := s.r.Insert(&m); err != nil {
					log.Printf("%v", err)
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
func (s *ApplicantService) Search() (*model.ApplicantsDownloadResponse, *model.ErrorResponse) {
	applicants, err := s.r.Search()
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

	// 書類登録状況更新
	if req.NamePre == "resume" {
		if err := s.r.UpdateDocument(&model.Applicant{
			HashKey:         req.HashKey,
			Resume:          objName,
			CurriculumVitae: "",
		}); err != nil {
			log.Printf("%v", err)
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}
	if req.NamePre == "curriculum_vitae" {
		if err := s.r.UpdateDocument(&model.Applicant{
			HashKey:         req.HashKey,
			Resume:          "",
			CurriculumVitae: objName,
		}); err != nil {
			log.Printf("%v", err)
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	return nil
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

	if err := s.r.UpdateDesiredAt(&model.Applicant{
		HashKey:   req.HashKey,
		DesiredAt: strings.Join(req.DesiredAt, ","),
	}); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
