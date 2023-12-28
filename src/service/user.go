package service

import (
	"api/resources/static"
	"api/src/model"
	"api/src/repository"
	"api/src/validator"
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	// 一覧
	List() (*model.UsersResponse, *model.ErrorResponse)
	// 登録
	Create(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// 取得
	Get(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// ロール一覧
	RoleList() (*model.UserRoles, *model.ErrorResponse)
	// 検索(グループ)
	SearchGroups() (*model.UserGroupsResponse, *model.ErrorResponse)
	// グループ登録
	CreateGroup(req *model.UserGroup) *model.ErrorResponse
}

type UserService struct {
	r repository.IUserRepository
	m repository.IMasterRepository
	v validator.IUserValidator
}

func NewUserService(
	r repository.IUserRepository,
	m repository.IMasterRepository,
	v validator.IUserValidator,
) IUserService {
	return &UserService{r, m, v}
}

// 一覧
func (u *UserService) List() (*model.UsersResponse, *model.ErrorResponse) {
	users, err := u.r.List()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UsersResponse{
		Users: users,
	}, nil
}

// 登録
func (u *UserService) Create(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	// バリデーション
	if err := u.v.CreateValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// メールアドレス重複チェック
	if err := u.r.EmailDuplCheck(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusConflict,
			Code:   static.CODE_USER_EMAIL_DUPL,
		}
	}

	// 初回パスワード発行
	password, hashPassword, err := generateRandomStr(8, 16)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := generateRandomStr(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	user := model.User{
		HashKey:      *hashKey,
		Name:         req.Name,
		Email:        req.Email,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
		RoleID:       req.RoleID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.r.Insert(&user); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	res := model.UserResponse{
		Email:        user.Email,
		InitPassword: *password,
	}
	return &res, nil
}

// 取得
func (u *UserService) Get(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	if err := u.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	res, err := u.r.Get(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &*model.ConvertUser(res), nil
}

// ロール一覧
func (u *UserService) RoleList() (*model.UserRoles, *model.ErrorResponse) {
	roles, err := u.m.SelectRole()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UserRoles{Roles: *roles}, nil
}

// 検索(グループ)
func (u *UserService) SearchGroups() (*model.UserGroupsResponse, *model.ErrorResponse) {
	userGroups, err := u.r.SearchGroup()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー存在確認
	for index, userGroup := range userGroups {
		var l []string
		users, err := u.r.ConfirmUserByHashKeys(strings.Split(userGroup.Users, ","))
		if err != nil {
			if err != nil {
				return nil, &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}
		}
		if users == nil || len(users) == 0 {
			userGroups[index].Users = ""
			continue
		}

		for _, user := range users {
			l = append(l, user.Name)
		}
		userGroups[index].Users = strings.Join(l, ",")
	}

	return &model.UserGroupsResponse{
		UserGroups: userGroups,
	}, nil
}

// グループ登録
func (u *UserService) CreateGroup(req *model.UserGroup) *model.ErrorResponse {
	// バリデーション
	if err := u.v.CreateGroupValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー存在確認
	users, err := u.r.ConfirmUserByHashKeys(strings.Split(req.Users, ","))
	if err != nil {
		if err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	var l []string
	for _, row := range users {
		l = append(l, row.HashKey)
	}
	req.Users = strings.Join(l, ",")

	// ハッシュキー生成
	_, hashKey, err := generateRandomStr(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// グループ登録
	req.HashKey = *hashKey
	if err := u.r.InsertGroup(req); err != nil {
		if err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	return nil
}

func generateRandomStr(minLength, maxLength int) (*string, *string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if err != nil {
		return nil, nil, err
	}
	strLength := minLength + int(length.Int64())

	buffer := make([]byte, strLength)
	_, err = rand.Read(buffer)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < strLength; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	str := string(buffer)

	buffer2, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return nil, nil, err
	}
	hash := string(buffer2)

	return &str, &hash, nil
}
