package middleware

import (
	"net/http"

	"github.com/google/uuid"
	httputil "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/http"
)

// RequireSameUser は本人確認を行う認可ミドルウェア
// 引数:
//   - next: 次のハンドラー
// 戻り値: http.Handler
// 実装:
//   1. Contextからログインユーザー情報を取得
//   2. パスパラメータからリソースのユーザーIDを取得
//   3. ログインユーザーIDとリソースのユーザーIDが一致するか確認
//   4. 一致する場合は次のハンドラーを実行、一致しない場合は403を返す
// 注意事項:
//   - このミドルウェアの前にJwtVerifyミドルウェアを適用する必要がある
//   - パスパラメータのキーは"id"である必要がある
//   - ユーザー情報がContextに存在しない場合は401を返す
//   - UUIDのパースに失敗した場合は400を返す
func RequireSameUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Contextからログインユーザー情報を取得
		userInfo, ok := GetUserInfoFromContext(r.Context())
		if !ok {
			httputil.WriteError(w, "user info not found in context", http.StatusUnauthorized)
			return
		}

		// パスパラメータからリソースのユーザーIDを取得
		idStr := r.PathValue("id")
		if idStr == "" {
			httputil.WriteError(w, "user id is required", http.StatusBadRequest)
			return
		}

		// リソースのユーザーIDをUUIDに変換
		resourceUserID, err := uuid.Parse(idStr)
		if err != nil {
			httputil.WriteError(w, "invalid user id", http.StatusBadRequest)
			return
		}

		// ログインユーザーIDをUUIDに変換
		loginUserID, err := uuid.Parse(userInfo.Sub)
		if err != nil {
			httputil.WriteError(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		// 本人確認: ログインユーザーIDとリソースのユーザーIDが一致するか
		if loginUserID != resourceUserID {
			httputil.WriteError(w, "forbidden: cannot access other user's information", http.StatusForbidden)
			return
		}

		// 次のハンドラーを実行
		next.ServeHTTP(w, r)
	})
}
