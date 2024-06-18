package repository

import (
	"api/src/model/static"
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type IGoogleRepository interface {
	// refresh_token取得
	GetRefreshToken() (*string, error)
	// 認証クライアント作成
	GetOauthClient() (*oauth2.Config, error)
	// 認証URL作成
	GetOauthURL() (*string, error)
	// access_token取得
	GetAccessToken(refreshToken *string, code *string) (*oauth2.Token, error)
	// Google Meet Url 取得
	GetGoogleMeetUrl(token *oauth2.Token, title string, start, end time.Time) (*string, error)
}

type GoogleRepository struct {
	redis *redis.Client
}

func NewGoogleRepository(redis *redis.Client) IGoogleRepository {
	return &GoogleRepository{redis}
}

// refresh_token取得
func (g *GoogleRepository) GetRefreshToken() (*string, error) {
	var ctx = context.Background()

	value, err := g.redis.Get(ctx, static.REDIS_OAUTH_REFRESH_TOKEN).Result()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &value, nil
}

// 認証クライアント作成
func (g *GoogleRepository) GetOauthClient() (*oauth2.Config, error) {
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
func (g *GoogleRepository) GetOauthURL() (*string, error) {
	config, err := g.GetOauthClient()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return &authURL, nil
}

// access_token取得
func (g *GoogleRepository) GetAccessToken(refreshToken *string, code *string) (*oauth2.Token, error) {
	config, err := g.GetOauthClient()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// 認証が必要な場合
	if code != nil && *code != "" {
		var ctx = context.Background()
		tok, err := config.Exchange(ctx, *code)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}

		return tok, nil
	}

	// refresh_tokenがある場合
	if refreshToken != nil && *refreshToken != "" {
		var ctx = context.Background()
		tok := &oauth2.Token{
			RefreshToken: *refreshToken,
		}
		res, err := config.TokenSource(ctx, tok).Token()
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}
		return res, nil
	}

	return nil, nil
}

// Google Meet Url 取得
func (g *GoogleRepository) GetGoogleMeetUrl(token *oauth2.Token, title string, start, end time.Time) (*string, error) {
	ctx := context.Background()

	config, configErr := g.GetOauthClient()
	if configErr != nil {
		log.Printf("%v", configErr)
		return nil, configErr
	}
	client := config.Client(ctx, token)

	// Calendar APIクライアントの作成
	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// イベントの作成
	event := &calendar.Event{
		Summary: title,
		Start:   &calendar.EventDateTime{DateTime: start.Format("2006-01-02T15:04:05"), TimeZone: "Asia/Tokyo"},
		End:     &calendar.EventDateTime{DateTime: end.Format("2006-01-02T15:04:05"), TimeZone: "Asia/Tokyo"},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{RequestId: "unique-request-id"},
		},
	}

	event, err = calendarService.Events.Insert("primary", event).ConferenceDataVersion(1).Do()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// Google Meet URLを返す
	if event.HangoutLink != "" {
		return &event.HangoutLink, nil
	}

	return nil, nil
}
