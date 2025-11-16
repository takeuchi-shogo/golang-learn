package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
	httputil "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/http"
)

// contextKey はContext内でキーとして使用するための型
// 意味: 外部パッケージからのキー衝突を防ぐための専用型
type contextKey string

// UserInfoKey はContextにユーザー情報を保存する際のキー
// 意味: Context.Value()でユーザー情報を取得する際に使用
const UserInfoKey contextKey = "userInfo"

// GetUserInfoFromContext はContextからユーザー情報を取得するヘルパー関数
// 引数:
//   - ctx: コンテキスト
// 戻り値:
//   - *jwt.UserInfo: ユーザー情報(存在する場合)
//   - bool: ユーザー情報が存在するかどうか
// 実装: Context.Value()でUserInfoKeyを使ってユーザー情報を取得
// 注意事項: ユーザー情報が存在しない、または型アサーションに失敗した場合はfalseを返す
func GetUserInfoFromContext(ctx context.Context) (*jwt.UserInfo, bool) {
	userInfo, ok := ctx.Value(UserInfoKey).(*jwt.UserInfo)
	return userInfo, ok
}

// JwtVerify はJWTトークンを検証し、ユーザー情報をContextに追加するミドルウェア
// 引数:
//   - jwtManager: JWT検証を行うマネージャー
//   - next: 次のハンドラー
// 戻り値: http.Handler
// 実装:
//   1. Authorizationヘッダーからトークンを取得
//   2. "Bearer "プレフィックスを除去
//   3. JWTトークンを検証
//   4. ユーザー情報を取得してContextに追加
//   5. 次のハンドラーを実行
// 注意事項:
//   - トークンが無効な場合は401を返す
//   - トークンがない場合は401を返す
func JwtVerify(jwtManager jwt.JwtManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorizationヘッダーからトークンを取得
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.WriteError(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// "Bearer "プレフィックスを除去
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// プレフィックスが存在しない場合
				httputil.WriteError(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
				return
			}

			// JWTトークンを検証
			token, err := jwtManager.VerifyToken(r.Context(), tokenString)
			if err != nil {
				log.Printf("Failed to verify token: %v", err)
				httputil.WriteError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// ユーザー情報を取得
			userInfo, err := jwtManager.GetUserInfo(token)
			if err != nil {
				log.Printf("Failed to get user info: %v", err)
				httputil.WriteError(w, "Failed to get user info", http.StatusInternalServerError)
				return
			}

			// ContextにユーザーIDを追加
			ctx := context.WithValue(r.Context(), UserInfoKey, userInfo)

			// 次のハンドラーを実行
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
