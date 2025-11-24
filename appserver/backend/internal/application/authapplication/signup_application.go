package authapplication

import (
	"context"
	"errors"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	usercommand "github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/command"
)

type SignupApplication interface {
	Run(ctx context.Context, email, userSub string) error
}

type signupApplication struct {
	userRepository usercommand.UserCommand
}

var _ SignupApplication = (*signupApplication)(nil)

func NewSignupApplication(userRepository usercommand.UserCommand) SignupApplication {
	return &signupApplication{userRepository: userRepository}
}

func (a *signupApplication) Run(ctx context.Context, email, userSub string) error {
	user := model.NewUser(model.Name(email), model.Email(email), userSub)
	createdUser, err := a.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	if createdUser == nil {
		return errors.New("failed to create user")
	}
	return nil
}
