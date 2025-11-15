package cognito

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type Cognito interface {
	SignUp(ctx context.Context, userID, email, password string) (*SignUpResult, error)
	ConfirmSignUp(ctx context.Context, email, code string) error
	SignIn(ctx context.Context, email, password string) (*AuthTokens, error)
	AdminCreateUser(ctx context.Context, clientID, userID, email string) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error)
}

type cognito struct {
	client     *cognitoidentityprovider.Client
	clientID   string
	userPoolID string
}

func New() Cognito {
	// ローカル開発用のmagnito設定
	// magnitoはローカルCognitoエミュレータ
	// 重要: AWS SDK v2はリージョン名にバリデーションがあるため、
	//       標準のAWSリージョン名を使用し、エンドポイントのみカスタマイズする
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		// 標準のAWSリージョン名を使用（SDK v2の要件）
		awsconfig.WithRegion("ap-northeast-1"),
		// magnitoのアクセスキーとシークレットキーを設定
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"magnito-access-key", // Access Key ID (compose.ymlのCOGNITO_ACCESS_KEY)
			"magnito-secret-key", // Secret Access Key (compose.ymlのCOGNITO_SECRET_KEY)
			"",                   // Session Token (不要)
		)),
	)
	if err != nil {
		panic(err)
	}
	return &cognito{
		client: cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
			// magnitoのローカルエンドポイントを指定
			// これにより、実際のAWSではなくローカルのmagnitoに接続される
			o.BaseEndpoint = aws.String("http://localhost:5050")
		}),
		clientID:   "magnito-client-name",  // compose.ymlのCOGNITO_USER_POOL_CLIENT_ID
		userPoolID: "ap-northeast-1_local", // compose.ymlのCOGNITO_USER_POOL_ID
	}
}

func newSignUpInput(clientID, email, password string) *cognitoidentityprovider.SignUpInput {
	return &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientID),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
}

type SignUpResult struct {
	UserSub       string `json:"user_sub"`       // Cognito User ID
	UserConfirmed bool   `json:"user_confirmed"` // 確認済みかどうか
}

func (c *cognito) SignUp(ctx context.Context, userID, email, password string) (*SignUpResult, error) {
	input := newSignUpInput(c.clientID, email, password)
	log.Printf("SignUp request: ClientID=%s, Username=%s, Email=%s, Endpoint=%s",
		c.clientID, email, email, "http://localhost:5050")

	opt, err := c.client.SignUp(ctx, input)
	if err != nil {
		log.Printf("SignUp error details: %+v", err)

		// 特定のエラータイプを判定
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			log.Printf("Invalid password error: %v", aws.ToString(invalidPassword.Message))
			return nil, fmt.Errorf("invalid password: %s", aws.ToString(invalidPassword.Message))
		}

		var userNameExists *types.UsernameExistsException
		if errors.As(err, &userNameExists) {
			log.Printf("Username exists error: %v", aws.ToString(userNameExists.Message))
			return nil, fmt.Errorf("user name already exists: %s", aws.ToString(userNameExists.Message))
		}

		// その他のエラー（403など）の場合、magnitoやエンドポイントの問題の可能性
		log.Printf("SignUp failed - Check magnito is running: docker compose ps")
		log.Printf("SignUp failed - Check magnito config: ClientID=%s, Endpoint=%s", c.clientID, "http://localhost:5050")
		return nil, fmt.Errorf("failed to sign up (check magnito status): %w", err)
	}

	log.Printf("SignUp success: UserSub=%s, UserConfirmed=%t", aws.ToString(opt.UserSub), opt.UserConfirmed)
	return &SignUpResult{
		UserSub:       aws.ToString(opt.UserSub),
		UserConfirmed: opt.UserConfirmed,
	}, nil
}

func (c *cognito) ConfirmSignUp(ctx context.Context, email, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	if _, err := c.client.ConfirmSignUp(ctx, input); err != nil {
		return fmt.Errorf("failed to confirm sign up: %w", err)
	}
	return nil
}

func (c *cognito) SignIn(ctx context.Context, email, password string) (*AuthTokens, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	resp, err := c.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}

	return &AuthTokens{
		IDToken:      aws.ToString(resp.AuthenticationResult.IdToken),
		AccessToken:  aws.ToString(resp.AuthenticationResult.AccessToken),
		RefreshToken: aws.ToString(resp.AuthenticationResult.RefreshToken),
		ExpiresIn:    resp.AuthenticationResult.ExpiresIn,
	}, nil
}

type AuthTokens struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

func newAdminCreateUserInput(userPoolID, userID, email string) *cognitoidentityprovider.AdminCreateUserInput {
	return &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:    aws.String(userPoolID),
		Username:      aws.String(userID),
		MessageAction: types.MessageActionTypeSuppress,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	}
}

func (c *cognito) AdminCreateUser(ctx context.Context, clientID, userID, email string) (*cognitoidentityprovider.AdminCreateUserOutput, error) {
	// UserPoolIDを使用（clientIDではない）
	input := newAdminCreateUserInput(c.userPoolID, userID, email)
	log.Printf("AdminCreateUser input: UserPoolID=%s, Username=%s, Email=%s", c.userPoolID, userID, email)

	opt, err := c.client.AdminCreateUser(ctx, input)
	if err != nil {
		log.Printf("AdminCreateUser error: %+v", err)
		return nil, fmt.Errorf("failed to admin create user: %w", err)
	}
	log.Printf("AdminCreateUser success: Username=%s", aws.ToString(opt.User.Username))
	return opt, nil
}

func (c *cognito) GetUser(ctx context.Context, accessToken string) (*cognitoidentityprovider.GetUserOutput, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: &accessToken,
	}
	return c.client.GetUser(ctx, input)
}
