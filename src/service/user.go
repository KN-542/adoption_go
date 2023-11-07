package service

import (
	"api/src/model"
	"api/src/repository"
	"api/src/validator"
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	// 一覧
	List() (*model.UsersResponse, error)
	// 登録
	Create(req *model.User) (*model.UserResponse, error)
}

type UserService struct {
	r repository.IUserRepository
	v validator.IUserValidator
}

func NewUserService(r repository.IUserRepository, v validator.IUserValidator) IUserService {
	return &UserService{r, v}
}

// 一覧
func (u *UserService) List() (*model.UsersResponse, error) {
	user, err := u.r.List()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &model.UsersResponse{
		Users: *model.ConvertUser(user),
	}, nil
}

// 登録
func (u *UserService) Create(req *model.User) (*model.UserResponse, error) {
	// バリデーション
	if err := u.v.CreateValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// 初回パスワード発行
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minLength := 8
	maxLength := 16

	length, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	passLength := minLength + int(length.Int64())

	buffer := make([]byte, passLength)
	_, err = rand.Read(buffer)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	for i := 0; i < passLength; i++ {
		buffer[i] = passwordChars[int(buffer[i])%len(passwordChars)]
	}
	password := string(buffer)

	buffer2, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	hashPassword := string(buffer2)

	// 登録
	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashPassword,
		InitPassword: hashPassword,
		RoleID:       req.RoleID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.r.Insert(&user); err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	res := model.UserResponse{
		Email:        user.Email,
		InitPassword: password,
	}
	return &res, nil
}
