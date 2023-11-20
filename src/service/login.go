package service

import (
	"api/resources/static"
	"api/src/model"
	"api/src/repository"
	"api/src/validator"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ILoginService interface {
	// ログイン認証
	Login(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// JWTトークン作成
	JWT(hashKey *string, exp time.Duration, name string, secret string) (*http.Cookie, *model.ErrorResponse)
	// MFA 認証コード生成
	CodeGenerate(req *model.User) *model.ErrorResponse
	// MFA
	MFA(req *model.UserMFA) *model.ErrorResponse
	// ユーザー存在確認
	UserCheck(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// パスワード変更 必要性
	PasswordChangeCheck(req *model.User) (*int8, *model.ErrorResponse)
	// パスワード変更
	PasswordChange(req *model.User) *model.ErrorResponse
	// セッション存在確認
	SessionConfirm(req *model.User) *model.ErrorResponse
	// ログアウト
	Logout(req *model.User) (*http.Cookie, *model.ErrorResponse)
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
func (l *LoginService) Login(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.LoginValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログイン認証
	users, err := l.login.Login(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	if len(users) != 1 {
		return nil, &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_AUTH,
		}
	}

	user := users[0]
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_AUTH,
		}
	}

	// Redisに保存
	ctx := context.Background()
	if err := l.redis.Set(
		ctx,
		user.HashKey,
		static.REDIS_USER_HASH_KEY,
		&user.HashKey,
		1*time.Hour,
	); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UserResponse{
		HashKey: user.HashKey,
		Name:    user.Name,
		Email:   user.Email,
		RoleID:  user.RoleID,
	}, nil
}

// JWTトークン作成
func (l *LoginService) JWT(hashKey *string, exp time.Duration, name string, secret string) (*http.Cookie, *model.ErrorResponse) {
	// Token作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": hashKey,
		"exp":     time.Now().Add(exp).Unix(),
	})

	// 署名
	tokenString, err := token.SignedString([]byte(os.Getenv(secret)))
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
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
func (l *LoginService) CodeGenerate(req *model.User) *model.ErrorResponse {
	// バリデーション
	if err := l.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}

	// 認証コード生成
	code := fmt.Sprintf("%06d", rand.Intn(int(math.Pow(10, 6))))

	// Redisに保存
	if err := l.redis.Set(
		ctx,
		req.HashKey,
		static.REDIS_CODE,
		&code,
		5*time.Minute,
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
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
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	code, err := l.redis.Get(
		ctx,
		req.HashKey,
		static.REDIS_CODE,
	)
	if err == redis.Nil {
		// Redisに保存
		if err := l.redis.Set(
			ctx,
			req.HashKey,
			static.REDIS_USER_HASH_KEY,
			&req.HashKey,
			1*time.Hour,
		); err != nil {
			log.Printf("%v", err)
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}

		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_EXPIRED,
		}
	} else if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if req.Code == *code {
		// 有効期限更新
		if err := l.redis.Set(
			ctx,
			req.HashKey,
			static.REDIS_USER_HASH_KEY,
			&req.HashKey,
			1*time.Hour,
		); err != nil {
			log.Printf("%v", err)
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}

		return nil
	} else {
		log.Printf("invalid code")
		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_INVALID_CODE,
		}
	}
}

// パスワード変更 必要性
func (l *LoginService) PasswordChangeCheck(req *model.User) (*int8, *model.ErrorResponse) {
	res, err := l.login.ConfirmInitPassword(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	return res, nil
}

// ユーザー存在確認
func (l *LoginService) UserCheck(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	res, err := l.login.Get(&model.User{HashKey: req.HashKey})
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return model.ConvertUser(res), nil
}

// パスワード変更
func (l *LoginService) PasswordChange(req *model.User) *model.ErrorResponse {
	// バリデーション
	if err := l.v.PasswordChangeValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 初期パスワード一致確認
	initPassword, err := l.login.ConfirmInitPassword2(req)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*initPassword),
		[]byte(req.InitPassword),
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_INIT_PASSWORD_INCORRECT,
		}
	}

	// パスワードハッシュ化
	buffer, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	req.Password = string(buffer)

	// パスワード変更
	req.UpdatedAt = time.Now()
	if err := l.login.PasswordChange(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// セッション存在確認
func (l *LoginService) SessionConfirm(req *model.User) *model.ErrorResponse {
	// バリデーション
	if err := l.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}

	// 有効期限更新
	if err := l.redis.Set(
		ctx,
		req.HashKey,
		static.REDIS_USER_HASH_KEY,
		&req.HashKey,
		1*time.Hour,
	); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// ログアウト
func (l *LoginService) Logout(req *model.User) (*http.Cookie, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// Redis破棄
	ctx := context.Background()
	if err := l.redis.Delete(ctx, req.HashKey); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// Cookie破棄(使えないcookieに更新)
	cookie := http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return &cookie, nil
}
