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
	// ログアウト
	Logout(req *request.Logout, token string) (*http.Cookie, *response.Error)
	// ログイン種別取得
	GetLoginType(req *request.GetLoginType) (*response.GetLoginType, *response.Error)
	// チーム存在確認(応募者)
	ConfirmTeamApplicant(req *request.ConfirmTeamApplicant) *response.Error
	// ログイン(応募者)
	LoginApplicant(req *request.LoginApplicant) (*response.LoginApplicant, *response.Error)
	// MFA 認証コード生成(応募者)
	CodeGenerateApplicant(req *request.CodeGenerateApplicant) *response.Error
	// MFA(応募者)
	MFAApplicant(req *request.MFAApplicant) *response.Error
	// ログアウト(応募者)
	LogoutApplicant(req *request.LogoutApplicant, token string) (*http.Cookie, *response.Error)
	// 応募者チェック
	CheckApplicant(req *request.CheckApplicant) *response.Error
}

type LoginService struct {
	login     repository.IUserRepository
	team      repository.ITeamRepository
	applicant repository.IApplicantRepository
	redis     repository.IRedisRepository
	v         validator.ILoginValidator
	v_0       validator.IUserValidator
	d         repository.IDBRepository
}

func NewLoginService(
	login repository.IUserRepository,
	team repository.ITeamRepository,
	applicant repository.IApplicantRepository,
	redis repository.IRedisRepository,
	v validator.ILoginValidator,
	v_0 validator.IUserValidator,
	d repository.IDBRepository,
) ILoginService {
	return &LoginService{login, team, applicant, redis, v, v_0, d}
}

// ログイン認証
func (l *LoginService) Login(req *request.Login) (*response.Login, *response.Error) {
	// バリデーション
	if err := l.v.Login(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
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
	teams, teamErr := l.team.ListTeamAssociation(&ddl.TeamAssociation{UserID: user.ID})
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

// ユーザー存在確認
func (l *LoginService) UserCheck(req *request.JWTDecode) *response.Error {
	// バリデーション
	if err := l.v.JWTDecode(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
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

// ログアウト
func (l *LoginService) Logout(req *request.Logout, token string) (*http.Cookie, *response.Error) {
	// バリデーション
	if err := l.v.Logout(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
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

// チーム存在確認(応募者)
func (l *LoginService) ConfirmTeamApplicant(req *request.ConfirmTeamApplicant) *response.Error {
	// バリデーション
	if err := l.v.ConfirmTeamApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// チーム取得
	_, err := l.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if err != nil {
		return &response.Error{
			Status: http.StatusNotFound,
			Code:   static.CODE_CONFIRM_TEAM_NOT_EXIST,
		}
	}

	return nil
}

// ログイン(応募者)
func (l *LoginService) LoginApplicant(req *request.LoginApplicant) (*response.LoginApplicant, *response.Error) {
	// バリデーション
	if err := l.v.LoginApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// チーム取得
	team, teamErr := l.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.TeamHashKey,
		},
	})
	if teamErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.TeamID = team.ID

	// ログイン認証
	applicant, err := l.applicant.GetByEmail(&req.Applicant)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_AUTH,
		}
	}

	// 応募者チェック
	if err := l.CheckApplicant(&request.CheckApplicant{
		Applicant: ddl.Applicant{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: applicant.HashKey,
			},
		},
	}); err != nil {
		return nil, &response.Error{
			Status: err.Status,
			Code:   err.Code,
		}
	}

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
	teamID := strconv.FormatUint(applicant.TeamID, 10)
	if err := l.redis.Set(
		ctx,
		applicant.HashKey,
		static.REDIS_USER_TEAM_ID,
		&teamID,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	s3Name := applicant.Name + "_" + strings.Replace(applicant.Email, ".", "", -1)
	if err := l.redis.Set(
		ctx,
		applicant.HashKey,
		static.REDIS_S3_NAME,
		&s3Name,
		24*time.Hour,
	); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.LoginApplicant{
		Applicant: entity.Applicant{
			Applicant: ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: applicant.HashKey,
				},
				Name:  applicant.Name,
				Email: applicant.Email,
			},
		},
	}, nil
}

// MFA 認証コード生成(応募者)
func (l *LoginService) CodeGenerateApplicant(req *request.CodeGenerateApplicant) *response.Error {
	// バリデーション
	if err := l.v.CodeGenerateApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
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

// MFA(応募者)
func (l *LoginService) MFAApplicant(req *request.MFAApplicant) *response.Error {
	// バリデーション
	if err := l.v.MFAApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
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
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}

			return &response.Error{
				Status: http.StatusUnauthorized,
				Code:   static.CODE_EXPIRED,
			}
		}
		return &response.Error{
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
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		return nil
	} else {
		log.Printf("invalid code")
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_INVALID_CODE,
		}
	}
}

// ログアウト(応募者)
func (l *LoginService) LogoutApplicant(req *request.LogoutApplicant, token string) (*http.Cookie, *response.Error) {
	// バリデーション
	if err := l.v.LogoutApplicant(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
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

// 応募者チェック
func (l *LoginService) CheckApplicant(req *request.CheckApplicant) *response.Error {
	// バリデーション
	if err := l.v.CheckApplicant(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 応募者取得
	applicant, applicantErr := l.applicant.Get(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if applicantErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 書類選考フラグチェック ＆ 過程チェック
	if applicant.DocumentPassFlg == static.DOCUMENT_FAIL ||
		applicant.ProcessingID == static.INTERVIEW_PROCESSING_PASS ||
		applicant.ProcessingID == static.INTERVIEW_PROCESSING_FAIL {
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_CHECK_APPLICANT_TEST_FINISHED,
		}
	}

	// 面接予定日チェック
	if applicant.ScheduleID > 0 && time.Now().After(applicant.Start.Add(-24*time.Hour)) {
		return &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_CHECK_APPLICANT_CANNOT_UPDATE_SCHEDULE,
		}
	}

	return nil
}
