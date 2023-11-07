package service

import (
	"api/src/model"
	"api/src/repository"
	"log"
)

type IUserService interface {
	// 一覧
	List() (*model.UsersResponse, error)
}

type UserService struct {
	r repository.IUserRepository
}

func NewUserService(r repository.IUserRepository) IUserService {
	return &UserService{r}
}

// 一覧
func (u *UserService) List() (*model.UsersResponse, error) {
	user, err := u.r.List()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &model.UsersResponse{
		Users: *model.ConvertUser(user),
	}, nil
}
