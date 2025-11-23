package authcontroller

import (
	"context"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/authapplication"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/cognito"
)

type SignupController struct {
	cognitoClient     cognito.Cognito
	jwtManager        jwt.JwtManager
	signupApplication authapplication.SignupApplication
}

func NewSignupController(cognitoClient cognito.Cognito, jwtManager jwt.JwtManager, signupApplication authapplication.SignupApplication) *SignupController {
	return &SignupController{
		cognitoClient:     cognitoClient,
		jwtManager:        jwtManager,
		signupApplication: signupApplication,
	}
}

func (c *SignupController) Signup(ctx context.Context, email, password string) error {
	signUpResult, err := c.cognitoClient.SignUp(ctx, "magnito-client-name", email, password)
	if err != nil {
		return err
	}

	// TODO: ユーザー情報をデータベースに保存
	if err := c.signupApplication.Run(ctx, email, signUpResult.UserSub); err != nil {
		return err
	}

	return nil
}

func (c *SignupController) ConfirmSignup(ctx context.Context, email, code string) error {
	err := c.cognitoClient.ConfirmSignUp(ctx, email, code)
	if err != nil {
		return err
	}
	return nil
}
