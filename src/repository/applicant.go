package repository

import (
	"api/resources/static"
	model "api/src/model"
	enum "api/src/model/enum"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	/*
		OAuth2.0用(削除予定)
	*/
	// refresh_token取得
	GetRefreshToken() (*string, error)
	// 認証クライアント作成
	GetOauthClient() (*oauth2.Config, error)
	// 認証URL作成
	GetOauthURL() (*model.GetOauthURLResponse, error)
	// access_token取得
	GetAccessToken(refreshToken *string, code *string) (*oauth2.Token, error)
	// シート取得
	GetSheets(search model.ApplicantSearch, token *oauth2.Token) (*[]model.ApplicantResponse, error)
	/*
		txt、csvダウンロード用
	*/
	// 登録
	Insert(applicant *model.Applicant) error
	// 検索
	Search() ([]model.Applicant, error)
	// PK検索(カウント)
	CountByPrimaryKey(key *string) (*int64, error)
}

type ApplicantRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewApplicantRepository(db *gorm.DB, redis *redis.Client) IApplicantRepository {
	return &ApplicantRepository{db, redis}
}

/*
	OAuth2.0用(削除予定)
*/
// refresh_token取得
func (a *ApplicantRepository) GetRefreshToken() (*string, error) {
	var ctx = context.Background()

	value, err := a.redis.Get(ctx, static.REDIS_OAUTH_REFRESH_TOKEN).Result()
	if err != nil {
		return nil, err
	}

	return &value, nil
}

// 認証クライアント作成
func (a *ApplicantRepository) GetOauthClient() (*oauth2.Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH_REDIRECT_URI"),
		Scopes:       []string{os.Getenv("AUTH_SCOPE_URI")},
		Endpoint:     google.Endpoint,
	}

	return &config, nil
}

// 認証URL作成
func (a *ApplicantRepository) GetOauthURL() (*model.GetOauthURLResponse, error) {
	_, err := a.GetRefreshToken()
	if err == nil {
		return nil, nil
	}

	config, err := a.GetOauthClient()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return &model.GetOauthURLResponse{Url: authURL}, nil
}

// access_token取得
func (a *ApplicantRepository) GetAccessToken(refreshToken *string, code *string) (*oauth2.Token, error) {
	config, err := a.GetOauthClient()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 認証が必要な場合
	if code != nil && *code != "" {
		var ctx = context.Background()
		tok, err := config.Exchange(ctx, *code)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		a.redis.Set(ctx, static.REDIS_OAUTH_REFRESH_TOKEN, tok.RefreshToken, 0)
		return tok, nil
	}

	// refresh_tokenがある場合
	if refreshToken != nil {
		var ctx = context.Background()
		tok := &oauth2.Token{
			RefreshToken: *refreshToken,
		}
		res, err := config.TokenSource(ctx, tok).Token()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return res, nil
	}

	return nil, nil
}

// シート取得
func (a *ApplicantRepository) GetSheets(search model.ApplicantSearch, token *oauth2.Token) (*[]model.ApplicantResponse, error) {
	ctx := context.Background()

	config, err := a.GetOauthClient()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	client := config.Client(ctx, token)

	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	readRange := fmt.Sprintf(
		"%s%s:%s%s",
		search.StartCellColumn,
		strconv.Itoa(search.StartCellRow),
		search.EndCellColumn,
		strconv.Itoa(search.EndCellRow))
	res, err := sheetsService.Spreadsheets.Values.Get(os.Getenv("SPREADSHEET_ID"), readRange).Do()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	respList := []model.ApplicantResponse{}
	for _, row := range res.Values {
		respList = append(
			respList,
			model.ApplicantResponse{
				Name:  row[enum.CELL_NAME].(string),
				Email: row[enum.CELL_EMAIL].(string),
			})
	}
	return &respList, nil
}

// 登録
func (a *ApplicantRepository) Insert(applicant *model.Applicant) error {
	if err := a.db.Create(applicant).Error; err != nil {
		return err
	}
	return nil
}

// 検索 TODO 検索仕様追加
func (a *ApplicantRepository) Search() ([]model.Applicant, error) {
	var l []model.Applicant
	if err := a.db.Find(&l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

// PK検索(カウント)
func (a *ApplicantRepository) CountByPrimaryKey(key *string) (*int64, error) {
	var count int64
	if err := a.db.Model(&model.Applicant{}).Where("id = ?", key).Count(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}
