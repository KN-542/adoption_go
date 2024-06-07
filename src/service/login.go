package service

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ILoginService interface {
	// ログイン認証
	Login(req *request.Login) (*response.Login, *response.Error)
	// MFA 認証コード生成
	CodeGenerate(req *request.CodeGenerate) *response.Error
	// MFA
	MFA(req *request.MFA) (*response.MFA, *response.Error)
	// JWTトークン作成
	JWT(hashKey *string, name string, secret string) (*http.Cookie, *response.Error)
	// JWT検証
	JWTDecode(cookie *http.Cookie, secret string) *response.Error
	// ユーザー存在確認
	UserCheck(req *request.JWTDecode) *response.Error
	// パスワード変更
	PasswordChange(req *request.PasswordChange) *response.Error
	// パスワード変更 必要性
	PasswordChangeCheck(req *ddl.User) (*int8, *response.Error)
	// セッション存在確認
	SessionConfirm(req *ddl.User) *response.Error
	// ログアウト
	Logout(req *request.Logout, token string) (*http.Cookie, *response.Error)
	// ログイン種別取得
	GetLoginType(req *request.GetLoginType) (*response.GetLoginType, *response.Error)
	// ログイン(応募者)*
	LoginApplicant(req *ddl.Applicant) (*entity.Applicant, *response.Error)
	// MFA 認証コード生成(応募者)*
	CodeGenerateApplicant(req *ddl.Applicant) *response.Error
	// ユーザー存在確認(応募者)*
	UserCheckApplicant(req *ddl.Applicant) (*ddl.Applicant, *response.Error)
	// セッション存在確認(応募者)*
	SessionConfirmApplicant(req *ddl.Applicant) *response.Error
	// ログアウト(応募者)*
	LogoutApplicant(req *ddl.Applicant) (*http.Cookie, *response.Error)
	// S3 Name Redisに登録*
	S3NamePreInsert(req *ddl.Applicant) *response.Error
}

type LoginService struct {
	login     repository.IUserRepository
	applicant repository.IApplicantRepository
	redis     repository.IRedisRepository
	v         validator.ILoginValidator
	v_0       validator.IUserValidator
	d         repository.IDBRepository
}

func NewLoginService(
	login repository.IUserRepository,
	applicant repository.IApplicantRepository,
	redis repository.IRedisRepository,
	v validator.ILoginValidator,
	v_0 validator.IUserValidator,
	d repository.IDBRepository,
) ILoginService {
	return &LoginService{login, applicant, redis, v, v_0, d}
}

// ログイン認証
func (l *LoginService) Login(req *request.Login) (*response.Login, *response.Error) {
	// バリデーション
	if err := l.v.Login(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログイン認証
	users, loginErr := l.login.Login(&req.User)
	if loginErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(users) != 1 {
		return nil, &response.Error{
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
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_AUTH,
		}
	}

	// チーム一覧取得
	teams, teamErr := l.login.ListTeamAssociation(&ddl.TeamAssociation{UserID: user.ID})
	if teamErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
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
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	roleID := strconv.FormatUint(uint64(user.RoleID), 10)
	if err := l.redis.Set(
		ctx,
		user.HashKey,
		static.REDIS_USER_ROLE,
		&roleID,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	loginType := strconv.FormatUint(uint64(user.UserType), 10)
	if err := l.redis.Set(
		ctx,
		user.HashKey,
		static.REDIS_USER_LOGIN_TYPE,
		&loginType,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID := strconv.FormatUint(uint64(user.CompanyID), 10)
	if err := l.redis.Set(
		ctx,
		user.HashKey,
		static.REDIS_USER_COMPANY_ID,
		&companyID,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(teams) > 0 {
		teamID := strconv.FormatUint(teams[0].TeamID, 10)
		if err := l.redis.Set(
			ctx,
			user.HashKey,
			static.REDIS_USER_TEAM_ID,
			&teamID,
			24*time.Hour,
		); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	return &response.Login{
		User: entity.User{
			User: ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: user.HashKey,
				},
				Name:  user.Name,
				Email: user.Email,
			},
		},
	}, nil
}

// MFA 認証コード生成
func (l *LoginService) CodeGenerate(req *request.CodeGenerate) *response.Error {
	// バリデーション
	if err := l.v.CodeGenerate(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &response.Error{
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
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// メール送信 TODO
	log.Print(code)

	return nil
}

// MFA
func (l *LoginService) MFA(req *request.MFA) (*response.MFA, *response.Error) {
	// バリデーション
	if err := l.v.MFA(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー取得＆パスワード変更必要性チェック
	user, userErr := l.login.Get(&req.User)
	if userErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	ctx := context.Background()
	code, codeErr := l.redis.Get(
		ctx,
		req.HashKey,
		static.REDIS_CODE,
	)
	if codeErr != nil {
		if codeErr == redis.Nil {
			// Redisに保存
			if err := l.redis.Set(
				ctx,
				req.HashKey,
				static.REDIS_USER_HASH_KEY,
				&req.HashKey,
				24*time.Hour,
			); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			return nil, &response.Error{
				Status: http.StatusUnauthorized,
				Code:   static.CODE_EXPIRED,
			}
		}
		return nil, &response.Error{
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
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// ログイン種別取得
		login, loginTypeErr := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
		if loginTypeErr != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		loginType, loginTypeParseErr := strconv.ParseUint(*login, 10, 64)
		if loginTypeParseErr != nil {
			log.Printf("%v", loginTypeParseErr)
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		var path string
		if loginType == uint64(static.LOGIN_TYPE_ADMIN) {
			path = "admin"
		} else {
			path = "management"
		}

		return &response.MFA{
			Path:             path,
			IsPasswordChange: user.Password == user.InitPassword,
		}, nil
	} else {
		log.Printf("invalid code")
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_INVALID_CODE,
		}
	}
}

// JWTトークン作成
func (l *LoginService) JWT(hashKey *string, name string, secret string) (*http.Cookie, *response.Error) {
	// Token作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": hashKey,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// 署名
	tokenString, err := token.SignedString([]byte(os.Getenv(secret)))
	if err != nil {
		return nil, &response.Error{
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
}

// JWT検証
func (l *LoginService) JWTDecode(cookie *http.Cookie, secret string) *response.Error {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Print(static.MESSAGE_UNEXPECTED_COOKIE)
			return nil, fmt.Errorf(static.MESSAGE_UNEXPECTED_COOKIE)
		}
		return []byte(os.Getenv(secret)), nil
	})
	if err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusUnauthorized,
		}
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.Print(static.MESSAGE_UNEXPECTED_COOKIE)
		return &response.Error{
			Status: http.StatusUnauthorized,
		}
	}

	return nil
}

// パスワード変更 必要性
func (l *LoginService) PasswordChangeCheck(req *ddl.User) (*int8, *response.Error) {
	res, err := l.login.ConfirmInitPassword(req)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	return res, nil
}

// ユーザー存在確認
func (l *LoginService) UserCheck(req *request.JWTDecode) *response.Error {
	// バリデーション
	if err := l.v.JWTDecode(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	_, loginErr := l.login.Get(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if loginErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	hash, hashErr := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if hashErr != nil {
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}

	// 有効期限更新
	if err := l.redis.Set(
		ctx,
		req.HashKey,
		static.REDIS_USER_HASH_KEY,
		hash,
		24*time.Hour,
	); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// パスワード変更
func (l *LoginService) PasswordChange(req *request.PasswordChange) *response.Error {
	// バリデーション
	if err := l.v.PasswordChange(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 初期パスワード一致確認
	user, confirmErr := l.login.Get(&req.User)
	if confirmErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.InitPassword),
		[]byte(req.InitPassword),
	); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_INIT_PASSWORD_INCORRECT,
		}
	}

	// パスワードハッシュ化
	buffer, bufferErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if bufferErr != nil {
		log.Printf("%v", bufferErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.Password = string(buffer)

	tx, txErr := l.d.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// パスワード変更
	if err := l.login.Update(tx, &ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   req.HashKey,
			UpdatedAt: time.Now(),
		},
		Password: req.Password,
	}); err != nil {
		if err := l.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := l.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// セッション存在確認
func (l *LoginService) SessionConfirm(req *ddl.User) *response.Error {
	// バリデーション
	if err := l.v_0.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &response.Error{
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
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// ログアウト
func (l *LoginService) Logout(req *request.Logout, token string) (*http.Cookie, *response.Error) {
	// バリデーション
	if err := l.v.Logout(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// Redis破棄
	ctx := context.Background()
	if err := l.redis.Delete(ctx, req.HashKey); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// Cookie破棄(使えないcookieに更新)
	cookie := http.Cookie{
		Name:     token,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return &cookie, nil
}

// ログイン種別取得
func (l *LoginService) GetLoginType(req *request.GetLoginType) (*response.GetLoginType, *response.Error) {
	// バリデーション
	if err := l.v.GetLoginType(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	ctx := context.Background()
	loginType_0, loginTypeErr := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if loginTypeErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	loginType, loginTypeParseErr := strconv.ParseUint(*loginType_0, 10, 64)
	if loginTypeParseErr != nil {
		log.Printf("%v", loginTypeParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.GetLoginType{
		LoginType: uint(loginType),
	}, nil
}

// ログイン(応募者)
func (l *LoginService) LoginApplicant(req *ddl.Applicant) (*entity.Applicant, *response.Error) {
	// バリデーション
	if err := l.v_0.LoginApplicantValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログイン認証
	applicants, err := l.applicant.GetByEmail(req)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if len(applicants) == 0 {
		return nil, &response.Error{
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
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &applicant, nil
}

// MFA 認証コード生成(応募者)
func (l *LoginService) CodeGenerateApplicant(req *ddl.Applicant) *response.Error {
	// バリデーション
	if err := l.v_0.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &response.Error{
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
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// メール送信 TODO
	log.Print(code)

	return nil
}

// ユーザー存在確認(応募者)
func (l *LoginService) UserCheckApplicant(req *ddl.Applicant) (*ddl.Applicant, *response.Error) {
	// バリデーション
	if err := l.v_0.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	res, err := l.applicant.GetByHashKey(req)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return res, nil
}

// セッション存在確認(応募者)
func (l *LoginService) SessionConfirmApplicant(req *ddl.Applicant) *response.Error {
	// バリデーション
	if err := l.v_0.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ログインの一時的セッション存在確認
	ctx := context.Background()
	_, err := l.redis.Get(ctx, req.HashKey, static.REDIS_USER_HASH_KEY)
	if err != nil {
		return &response.Error{
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
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// ログアウト(応募者)
func (l *LoginService) LogoutApplicant(req *ddl.Applicant) (*http.Cookie, *response.Error) {
	// バリデーション
	if err := l.v_0.HashKeyValidateApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// Redis破棄
	ctx := context.Background()
	if err := l.redis.Delete(ctx, req.HashKey); err != nil {
		return nil, &response.Error{
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
func (l *LoginService) S3NamePreInsert(req *ddl.Applicant) *response.Error {
	res, err := l.applicant.GetByHashKey(req)
	if err != nil {
		return &response.Error{
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
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
