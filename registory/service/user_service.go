package service

import (
	"errors"
	"fmt"

	"github.com/takeuchi-shogo/golang-learn/registory/repository"
)

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{repository: repository}
}

func (s *userService) Run() error {
	user, err := s.repository.Store().GetUserById(1)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if user == nil {
		fmt.Println("User not found")
		return errors.New("user not found")
	}
	fmt.Println("User name: ", user.Name)
	return nil
}
