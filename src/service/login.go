package service

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ILoginService interface {
	// JWTトークン作成
	JWT(email *string) (*http.Cookie, error)
}

type LoginService struct {
}

func NewLoginService() ILoginService {
	return &LoginService{}
}

// JWTトークン作成
func (l *LoginService) JWT(email *string) (*http.Cookie, error) {
	// Token作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": email,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	// 署名
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Cookie作成
	cookie := http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,  // JavaScriptからのアクセスを禁止する場合はtrueに
		Secure:   false, // HTTPSでのみ送信 開発環境ではfalseに
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return &cookie, nil
}
