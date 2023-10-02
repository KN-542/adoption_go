package service

import (
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"
)

type IApplicantService interface {
	/*
		OAuth2.0用(削除予定)
	*/
	// 認証URL作成
	GetOauthURL() (*model.GetOauthURLResponse, error)
	// シート取得
	GetSheets(search model.ApplicantSearch) (*[]model.ApplicantResponse, error)
	/*
		txt、csvダウンロード用
	*/
	// 応募者ダウンロード
	Download(d *model.ApplicantsDownload) (*model.ApplicantsDownloadResponse, error)
}

type applicantService struct {
	r repository.IApplicantRepository
	m repository.IMasterRepository
}

func NewApplicantService(r repository.IApplicantRepository, m repository.IMasterRepository) IApplicantService {
	return &applicantService{r, m}
}

// 認証URL作成
func (s *applicantService) GetOauthURL() (*model.GetOauthURLResponse, error) {
	resp, err := s.r.GetOauthURL()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return resp, nil
}

// シート取得
func (s *applicantService) GetSheets(search model.ApplicantSearch) (*[]model.ApplicantResponse, error) {
	refreshToken, _ := s.r.GetRefreshToken()

	accessToken, err := s.r.GetAccessToken(refreshToken, &search.Code)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	resp, err := s.r.GetSheets(search, accessToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return resp, nil
}

/*
	txt、csvダウンロード用
*/
// 応募者ダウンロード
func (s *applicantService) Download(d *model.ApplicantsDownload) (*model.ApplicantsDownloadResponse, error) {
	// STEP1 サイトIDチェック
	_, err := s.m.SelectByPrimaryKey(d.Site)
	if err != nil {
		log.Fatal(err)
		return nil, err
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
				log.Fatal(err)
				return nil, err
			}
			fmt.Println(m.Name)
			if *count == int64(0) {
				// STEP2-2 登録
				if err := s.r.Insert(&m); err != nil {
					log.Fatal(err)
					return nil, err
				}
			}
		}
	}

	// STEP3 検索
	applicants, err := s.r.Search()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	res := model.ApplicantsDownloadResponse{
		Applicants: applicants,
	}
	return &res, nil
}
