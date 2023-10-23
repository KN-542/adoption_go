package service

import (
	"api/src/model"
	"api/src/repository"
	"log"
)

type IUserService interface {
	// 一覧
	List() (*[]model.UserResponse, error)
}

type userService struct {
	r repository.IUserRepository
}

func NewUserService(r repository.IUserRepository) IUserService {
	return &userService{r}
}

// 一覧
func (u *userService) List() (*[]model.UserResponse, error) {
	user, err := u.r.List()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return model.ConvertUser(user), nil
}
