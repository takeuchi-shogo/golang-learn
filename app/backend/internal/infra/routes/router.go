package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/authapplication"
	userapplication "github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/userapplication"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/controllers"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/controllers/authcontroller"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/controllers/dto"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/service"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/factory"
	infracognito "github.com/takeuchi-shogo/golang-learn/app/backend/internal/infra/cognito"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/middleware"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	httputil "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/http"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/cognito"
)

// router は依存関係を組み立て、HTTPルーティングを設定する
// 実装:
//   1. 各エンドポイントを定義
//   2. 認証が必要なエンドポイントにはJwtVerifyミドルウェアを適用
// 注意事項:
//   - PUT /users/{id}は認証必須(JwtVerifyミドルウェア適用)
//   - 認証エンドポイント(signup, login)は認証不要
func router() *http.ServeMux {
	mux := http.NewServeMux()

	// JWT Manager の初期化(ミドルウェア用)
	jwtManager, err := jwt.NewJwtManager("ap-northeast-1", "ap-northeast-1_local")
	if err != nil {
		// 初期化失敗は致命的エラー
		panic("failed to initialize jwt manager: " + err.Error())
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// ルーティング設定
	mux.HandleFunc("GET /users/{id}", userRouter)

	// 認証が必要なエンドポイント: ユーザー情報更新
	// JwtVerifyミドルウェアを適用してContextにユーザー情報を追加
	mux.Handle("PUT /users/{id}", middleware.JwtVerify(jwtManager)(http.HandlerFunc(updateUserRouter)))

	// 認証不要なエンドポイント
	mux.HandleFunc("POST /auth/signup", authSignupRouter)
	mux.HandleFunc("POST /auth/login", authLoginRouter)
	return mux
}

// userRouter はユーザー取得のルーティングハンドラー
// 引数:
//   - w: HTTPレスポンスライター
//   - r: HTTPリクエスト
// 実装:
//   1. パスパラメータからユーザーIDを取得
//   2. コントローラーを初期化
//   3. ユーザー情報を取得
//   4. レスポンスを返却
func userRouter(w http.ResponseWriter, r *http.Request) {
	userController := controllers.NewUserController(userapplication.NewGetUser(factory.NewFactory().GetUserRegistory().UserQuery()))
	// リクエストパラメータの取得
	idStr := r.PathValue("id")

	// IDを文字列からintに変換
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// コントローラー経由でビジネスロジックを実行
	user, err := userController.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスの返却
	json.NewEncoder(w).Encode(user)
}

// updateUserRouter はユーザー更新のルーティングハンドラー
// 引数:
//   - w: HTTPレスポンスライター
//   - r: HTTPリクエスト
// 実装:
//   1. Contextからログインユーザー情報を取得(JWT認証済み)
//   2. パスパラメータからユーザーIDを取得
//   3. 本人確認(ログインユーザーIDと更新対象ユーザーIDが一致するか)
//   4. リクエストボディから更新情報を取得してバリデーション
//   5. 依存関係を組み立て
//   6. コントローラーを初期化
//   7. ユーザー情報を更新
//   8. レスポンスを返却
// 注意事項:
//   - JWT認証が必須(ミドルウェアで事前に検証)
//   - 他人のユーザー情報は更新できない(403エラー)
//   - name, emailはnilの場合は更新しない
//   - バリデーションエラーは400を返す
//   - ユーザーが見つからない場合は404を返す
func updateUserRouter(w http.ResponseWriter, r *http.Request) {
	// Contextからログインユーザー情報を取得
	// 注意: JwtVerifyミドルウェアで既に検証済みのため、存在が保証されている
	userInfo, ok := middleware.GetUserInfoFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, "user info not found in context", http.StatusUnauthorized)
		return
	}

	// ファクトリーから依存関係を取得
	f := factory.NewFactory()
	userRegistory := f.GetUserRegistory()

	// Cognitoクライアントの初期化
	cognitoClient := cognito.New()
	cognitoAdapter := infracognito.NewCognitoAdapter(cognitoClient)

	// UserSyncServiceの初期化
	userSyncService := service.NewUserSyncService(cognitoAdapter, userRegistory.UserCommand())

	// コントローラーを初期化
	userController := controllers.NewUserControllerWithUpdate(
		userapplication.NewGetUser(userRegistory.UserQuery()),
		userapplication.NewUpdateUser(userRegistory.UserQuery(), userSyncService),
	)

	// リクエストパラメータの取得
	idStr := r.PathValue("id")

	// IDを文字列からUUIDに変換
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// 本人確認: ログインユーザーIDと更新対象ユーザーIDが一致するか
	// ビジネスルール: ユーザーは自分自身の情報のみ更新可能
	loginUserID, err := uuid.Parse(userInfo.Sub)
	if err != nil {
		httputil.WriteError(w, "invalid user id in token", http.StatusUnauthorized)
		return
	}

	if loginUserID != id {
		httputil.WriteError(w, "forbidden: cannot update other user's information", http.StatusForbidden)
		return
	}

	// リクエストボディから更新情報を取得
	var updateUserRequest dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		httputil.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// バリデーション
	if err := updateUserRequest.Validate(); err != nil {
		httputil.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// DTOからドメインモデルの型に変換
	var name *model.Name
	var email *model.Email

	if updateUserRequest.Name != nil {
		n := model.Name(*updateUserRequest.Name)
		name = &n
	}
	if updateUserRequest.Email != nil {
		e := model.Email(*updateUserRequest.Email)
		email = &e
	}

	// コントローラー経由でビジネスロジックを実行
	user, err := userController.Update(r.Context(), id, name, email)
	if err != nil {
		// エラーの種類に応じてステータスコードを変更
		if errors.Is(err, repository.ErrUserNotFound) {
			httputil.WriteError(w, "user not found", http.StatusNotFound)
			return
		}
		httputil.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// レスポンスの返却
	httputil.WriteJSON(w, user, http.StatusOK)
}

func authSignupRouter(w http.ResponseWriter, r *http.Request) {
	cognitoClient := cognito.New()
	jwtManager, err := jwt.NewJwtManager("ap-northeast-1", "ap-northeast-1_local")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	signupController := authcontroller.NewSignupController(cognitoClient, jwtManager, authapplication.NewSignupApplication(factory.NewFactory().GetUserRegistory().UserCommand()))
	type SignupRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var signupRequest SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = signupController.Signup(r.Context(), signupRequest.Email, signupRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスの返却
	json.NewEncoder(w).Encode(map[string]string{"message": "signup successful"})
}

func authLoginRouter(w http.ResponseWriter, r *http.Request) {
	cognitoClient := cognito.New()
	jwtManager, err := jwt.NewJwtManager("ap-northeast-1", "ap-northeast-1_local")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginController := authcontroller.NewLoginController(cognitoClient, jwtManager)
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var loginRequest LoginRequest
	err = json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authTokens, err := loginController.Login(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスの返却
	json.NewEncoder(w).Encode(authTokens)
}

// NewRouter はルーターを初期化する
func NewRouter() *http.ServeMux {
	return router()
}
