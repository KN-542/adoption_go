package service

import (
	"api/resources/static"
	"api/src/model"
	enum "api/src/model/enum"
	"api/src/repository"
	"api/src/validator"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ILoginService interface {
	// ログイン認証
	Login(req *model.User) *model.ErrorResponse
	// JWTトークン作成
	JWT(email *string, exp float64, name string) (*http.Cookie, *model.ErrorResponse)
	// MFA 認証コード生成
	CodeGenerate() *model.ErrorResponse
	// MFA
	MFA(req *model.UserMFA) *model.ErrorResponse
	// ユーザー存在確認
	UserCheck() *model.ErrorResponse
	// パスワード変更 必要性
	PasswordChangeCheck() (*int8, *model.ErrorResponse)
	// パスワード変更
	PasswordChange(req *model.User) *model.ErrorResponse
}

type LoginService struct {
	login repository.IUserRepository
	redis repository.IRedisRepository
	v     validator.IUserValidator
}

func NewLoginService(
	login repository.IUserRepository,
	redis repository.IRedisRepository,
	v validator.IUserValidator,
) ILoginService {
	return &LoginService{login, redis, v}
}

// ログイン認証
func (l *LoginService) Login(req *model.User) *model.ErrorResponse {
	// バリデーション
	if err := l.v.LoginValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  err,
		}
	}

	// ログイン認証
	user, err := l.login.Login(req)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	// Redisに保存
	ctx := context.Background()
	id := strconv.FormatUint(user.ID, 10)
	if err := l.redis.Set(
		ctx,
		static.REDIS_USER_ID,
		&id,
		24*time.Hour,
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	if err := l.redis.Set(
		ctx,
		static.REDIS_EMAIL,
		&user.Email,
		24*time.Hour,
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

// JWTトークン作成
func (l *LoginService) JWT(email *string, exp float64, name string) (*http.Cookie, *model.ErrorResponse) {
	userId := email
	if userId == nil {
		ctx := context.Background()
		value, err := l.redis.Get(
			ctx,
			static.REDIS_EMAIL,
		)
		if err != nil {
			log.Printf("%v", err)
			return nil, &model.ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  err,
			}
		}
		userId = value
	}

	// Token作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Duration(exp) * time.Hour).Unix(),
	})

	// 署名
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // 複数のトークンがあるため別のキーを使用するかも
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	// Cookie作成
	cookie := http.Cookie{
		Name:     name,
		Value:    tokenString,
		Expires:  time.Now().Add(time.Duration(exp) * time.Hour),
		HttpOnly: true,  // JavaScriptからのアクセスを禁止する場合はtrueに
		Secure:   false, // HTTPSでのみ送信 開発環境ではfalseに
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return &cookie, nil
}

// MFA 認証コード生成
func (l *LoginService) CodeGenerate() *model.ErrorResponse {
	// 認証コード生成
	code := fmt.Sprintf("%06d", rand.Intn(10^6))

	// Redisに保存
	ctx := context.Background()
	if err := l.redis.Set(
		ctx,
		static.REDIS_CODE,
		&code,
		5*time.Minute,
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	// メール送信 TODO
	log.Print(code)

	return nil
}

// MFA
func (l *LoginService) MFA(req *model.UserMFA) *model.ErrorResponse {
	// バリデーション
	if err := l.v.MFAValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  err,
		}
	}

	// コード取得
	ctx := context.Background()
	code, err := l.redis.Get(
		ctx,
		static.REDIS_CODE,
	)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	if req.Code == *code {
		return nil
	} else {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("invalid code"),
		}
	}
}

// パスワード変更 必要性
func (l *LoginService) PasswordChangeCheck() (*int8, *model.ErrorResponse) {
	ctx := context.Background()
	value, err := l.redis.Get(
		ctx,
		static.REDIS_USER_ID,
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	userId, err := strconv.ParseUint(*value, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	res, err := l.login.ConfirmInitPassword(
		&model.User{ID: userId},
	)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	return res, nil
}

// ユーザー存在確認
func (l *LoginService) UserCheck() *model.ErrorResponse {
	ctx := context.Background()
	value, err := l.redis.Get(
		ctx,
		static.REDIS_USER_ID,
	)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	userId, err := strconv.ParseUint(*value, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusForbidden,
			Error:  err,
		}
	}

	if err := l.login.UserCheck(&model.User{ID: userId}); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}

// パスワード変更
func (l *LoginService) PasswordChange(req *model.User) *model.ErrorResponse {
	// バリデーション
	if err := l.v.LoginValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  err,
		}
	}

	ctx := context.Background()
	value, err := l.redis.Get(
		ctx,
		static.REDIS_USER_ID,
	)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}
	userId, err := strconv.ParseUint(*value, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusForbidden,
			Error:  err,
		}
	}

	req.ID = userId

	// 初期パスワード一致確認
	status, err := l.login.ConfirmInitPassword(req)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	if *status == int8(enum.PASSWORD_CHANGE_REQUIRED) {
		return &model.ErrorResponse{
			Status: http.StatusConflict,
			Error:  err,
		}
	}

	// パスワード変更
	if err := l.login.PasswordChange(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  err,
		}
	}

	return nil
}
