package service

import (
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"log"
	"net/http"
	"strconv"
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
}

type ApplicantService struct {
	r repository.IApplicantRepository
	m repository.IMasterRepository
}

func NewApplicantService(r repository.IApplicantRepository, m repository.IMasterRepository) IApplicantService {
	return &ApplicantService{r, m}
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
			m := model.Applicant{
				ID:     values[enum.RECRUIT_ID],
				SiteID: int(enum.RECRUIT),
				Name:   values[enum.RECRUIT_NAME],
				Email:  values[enum.RECRUIT_EMAIL],
				Tel:    values[enum.RECRUIT_TEL],
				Age:    int(age),
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
