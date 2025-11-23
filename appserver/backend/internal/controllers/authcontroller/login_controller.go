package authcontroller

import (
	"context"
	"log"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/cognito"
)

type LoginController struct {
	cognitoClient cognito.Cognito
	jwtManager    jwt.JwtManager
}

func NewLoginController(cognitoClient cognito.Cognito, jwtManager jwt.JwtManager) *LoginController {
	return &LoginController{
		cognitoClient: cognitoClient,
		jwtManager:    jwtManager,
	}
}

type Response struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

func (c *LoginController) Login(ctx context.Context, email, password string) (*Response, error) {
	log.Println("LoginController.Login", email, password)
	authTokens, err := c.cognitoClient.SignIn(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return &Response{
		IDToken:      authTokens.IDToken,
		AccessToken:  authTokens.AccessToken,
		RefreshToken: authTokens.RefreshToken,
		ExpiresIn:    authTokens.ExpiresIn,
	}, nil
}
