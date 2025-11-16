package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypeCognito
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeNone:
		return "none"
	case AuthTypeCognito:
		return "cognito"
	}
	return "unknown"
}

// JwtManager はJWTトークンの検証とユーザー情報取得を行うインターフェース
// 実装: jwtManager構造体
type JwtManager interface {
	// VerifyToken はJWTトークンを検証する
	// 引数:
	//   - ctx: コンテキスト
	//   - tokenString: JWT トークン文字列
	// 戻り値:
	//   - *jwt.Token: 検証済みトークン
	//   - error: エラー情報
	VerifyToken(ctx context.Context, tokenString string) (*jwt.Token, error)

	// GetUserInfo はトークンからユーザー情報を取得する
	// 引数:
	//   - token: 検証済みJWTトークン
	// 戻り値:
	//   - *UserInfo: ユーザー情報
	//   - error: エラー情報
	GetUserInfo(token *jwt.Token) (*UserInfo, error)
}

type jwtManager struct {
	region     string
	userPoolID string
	authType   AuthType
	issuer     string
	keys       map[string]*rsa.PublicKey
}

func NewJwtManager(region, userPoolID string) (JwtManager, error) {
	jwtManager := &jwtManager{
		region:     region,
		userPoolID: userPoolID,
		// ローカル開発環境用
		issuer: fmt.Sprintf("http://localhost:5050/%s", userPoolID),
		// 本番環境用
		// issuer:     fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID),
		authType: AuthTypeCognito,
		keys:     make(map[string]*rsa.PublicKey),
	}
	if err := jwtManager.fetchJWKS(); err != nil {
		return nil, err
	}
	return jwtManager, nil
}

func (v *jwtManager) fetchJWKS() error {
	// ローカル開発環境の場合はmagnitoを使用
	url := fmt.Sprintf("http://localhost:5050/%s/.well-known/jwks.json", v.userPoolID)
	// 本番環境の場合は以下を使用
	// url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
	// 	v.region, v.userPoolID)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
			Kty string `json:"kty"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// RSA 公開鍵を生成
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" {
			continue
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			return fmt.Errorf("failed to decode N: %w", err)
		}

		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			return fmt.Errorf("failed to decode E: %w", err)
		}

		var e int
		for _, b := range eBytes {
			e = e<<8 + int(b)
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}

		v.keys[key.Kid] = pubKey
	}

	return nil
}

func (v *jwtManager) VerifyToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	log.Println("VerifyToken", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズムチェック
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// kid (Key ID) を取得
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		// 対応する公開鍵を取得
		key, ok := v.keys[kid]
		if !ok {
			return nil, fmt.Errorf("key not found for kid: %s", kid)
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// クレームの検証
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Issuer チェック
	// ローカル開発環境用
	expectedIssuer := fmt.Sprintf("http://localhost:5050/%s", v.userPoolID)
	// 本番環境用
	// expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s",
	// 	v.region, v.userPoolID)
	if claims["iss"] != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	// Token Use チェック (ID Token であることを確認)
	if claims["token_use"] != "id" {
		return nil, fmt.Errorf("token is not an id token")
	}

	// 有効期限チェック
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return token, nil
}

// GetUserInfo はトークンからユーザー情報を取得する
// 引数:
//   - token: 検証済みJWTトークン
// 戻り値:
//   - *UserInfo: ユーザー情報(Sub, Email)
//   - error: エラー情報
// 実装: トークンのクレームからsub(ユーザーID)とemailを抽出
// 注意事項: クレームが不正な場合はエラーを返す
func (v *jwtManager) GetUserInfo(token *jwt.Token) (*UserInfo, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("sub claim not found or invalid")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email claim not found or invalid")
	}

	return &UserInfo{
		Sub:   sub,
		Email: email,
	}, nil
}

// UserInfo はJWTトークンから取得したユーザー情報
// 意味: 認証済みユーザーの基本情報を保持
type UserInfo struct {
	// Sub はユーザーの一意識別子(Cognito User Pool内のユーザーID)
	Sub string
	// Email はユーザーのメールアドレス
	Email string
}
