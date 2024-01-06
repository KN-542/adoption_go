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
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ILoginService interface {
	// ログイン認証
	Login(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// JWTトークン作成
	JWT(hashKey *string, name string, secret string) (*http.Cookie, *model.ErrorResponse)
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
	// ログイン(応募者)
	LoginApplicant(req *model.Applicant) (*model.Applicant, *model.ErrorResponse)
	// MFA 認証コード生成(応募者)
	CodeGenerateApplicant(req *model.Applicant) *model.ErrorResponse
	// ユーザー存在確認(応募者)
	UserCheckApplicant(req *model.Applicant) (*model.Applicant, *model.ErrorResponse)
	// セッション存在確認(応募者)
	SessionConfirmApplicant(req *model.Applicant) *model.ErrorResponse
	// ログアウト(応募者)
	LogoutApplicant(req *model.Applicant) (*http.Cookie, *model.ErrorResponse)
	// S3 Name Redisに登録
	S3NamePreInsert(req *model.Applicant) *model.ErrorResponse
}

type LoginService struct {
	login     repository.IUserRepository
	applicant repository.IApplicantRepository
	redis     repository.IRedisRepository
	v         validator.IUserValidator
	d         repository.IDBRepository
}

func NewLoginService(
	login repository.IUserRepository,
	applicant repository.IApplicantRepository,
	redis repository.IRedisRepository,
	v validator.IUserValidator,
	d repository.IDBRepository,
) ILoginService {
	return &LoginService{login, applicant, redis, v, d}
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
		24*time.Hour,
	); err != nil {
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
func (l *LoginService) JWT(hashKey *string, name string, secret string) (*http.Cookie, *model.ErrorResponse) {
	// Token作成
	if name == "jwt_token" || name == "jwt_token3" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": hashKey,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
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
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,  // JavaScriptからのアクセスを禁止する場合はtrueに
			Secure:   false, // HTTPSでのみ送信 開発環境ではfalseに
			SameSite: http.SameSiteDefaultMode,
			Path:     "/",
		}

		return &cookie, nil
	} else if name == "jwt_token2" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": hashKey,
			"exp":     time.Now().Add(24 * 30 * 3 * time.Hour).Unix(),
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
			Expires:  time.Now().Add(24 * 30 * 3 * time.Hour),
			HttpOnly: true,  // JavaScriptからのアクセスを禁止する場合はtrueに
			Secure:   false, // HTTPSでのみ送信 開発環境ではfalseに
			SameSite: http.SameSiteDefaultMode,
			Path:     "/",
		}

		return &cookie, nil
	}

	return nil, nil
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
			24*time.Hour,
		); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}

		return &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_EXPIRED,
		}
	} else if err != nil {
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
			24*time.Hour,
		); err != nil {
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

	tx, err := l.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// パスワード変更
	req.UpdatedAt = time.Now()
	if err := l.login.PasswordChange(tx, req); err != nil {
		if err := l.d.TxRollback(tx); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := l.d.TxCommit(tx); err != nil {
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
		24*time.Hour,
	); err != nil {
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

// ログイン(応募者)
func (l *LoginService) LoginApplicant(req *model.Applicant) (*model.Applicant, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.LoginApplicantValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログイン認証
	applicants, err := l.applicant.GetByEmail(req)
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	if len(applicants) == 0 {
		return nil, &model.ErrorResponse{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_AUTH,
		}
	}

	applicant := applicants[0]

	// Redisに保存
	ctx := context.Background()
	if err := l.redis.Set(
		ctx,
		applicant.HashKey,
		static.REDIS_USER_HASH_KEY,
		&applicant.HashKey,
		24*time.Hour,
	); err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &applicant, nil
}

// MFA 認証コード生成(応募者)
func (l *LoginService) CodeGenerateApplicant(req *model.Applicant) *model.ErrorResponse {
	// バリデーション
	if err := l.v.HashKeyValidateApplicant(req); err != nil {
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
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// メール送信 TODO
	log.Print(code)

	return nil
}

// ユーザー存在確認(応募者)
func (l *LoginService) UserCheckApplicant(req *model.Applicant) (*model.Applicant, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	res, err := l.applicant.GetByHashKey(req)
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return res, nil
}

// セッション存在確認(応募者)
func (l *LoginService) SessionConfirmApplicant(req *model.Applicant) *model.ErrorResponse {
	// バリデーション
	if err := l.v.HashKeyValidateApplicant(req); err != nil {
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
		24*time.Hour,
	); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// ログアウト(応募者)
func (l *LoginService) LogoutApplicant(req *model.Applicant) (*http.Cookie, *model.ErrorResponse) {
	// バリデーション
	if err := l.v.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// Redis破棄
	ctx := context.Background()
	if err := l.redis.Delete(ctx, req.HashKey); err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// Cookie破棄(使えないcookieに更新)
	cookie := http.Cookie{
		Name:     "jwt_token3",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return &cookie, nil
}

// S3 Name Redisに登録
func (l *LoginService) S3NamePreInsert(req *model.Applicant) *model.ErrorResponse {
	res, err := l.applicant.GetByHashKey(req)
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	s3Name := res.Name + "_" + strings.Replace(res.Email, ".", "", -1)

	ctx := context.Background()
	if err := l.redis.Set(
		ctx,
		req.HashKey,
		static.REDIS_S3_NAME,
		&s3Name,
		24*time.Hour,
	); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
