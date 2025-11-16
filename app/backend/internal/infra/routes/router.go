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
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/factory"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	httputil "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/http"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/cognito"
)

// router は依存関係を組み立て、HTTPルーティングを設定する
func router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// ルーティング設定
	mux.HandleFunc("GET /users/{id}", userRouter)
	mux.HandleFunc("PUT /users/{id}", updateUserRouter)

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
//   1. パスパラメータからユーザーIDを取得
//   2. リクエストボディから更新情報を取得してバリデーション
//   3. コントローラーを初期化
//   4. ユーザー情報を更新
//   5. レスポンスを返却
// 注意事項:
//   - name, emailはnilの場合は更新しない
//   - バリデーションエラーは400を返す
//   - ユーザーが見つからない場合は404を返す
func updateUserRouter(w http.ResponseWriter, r *http.Request) {
	// ファクトリーから依存関係を取得
	f := factory.NewFactory()
	userRegistory := f.GetUserRegistory()

	// コントローラーを初期化
	userController := controllers.NewUserControllerWithUpdate(
		userapplication.NewGetUser(userRegistory.UserQuery()),
		userapplication.NewUpdateUser(userRegistory.UserQuery(), userRegistory.UserCommand()),
	)

	// リクエストパラメータの取得
	idStr := r.PathValue("id")

	// IDを文字列からUUIDに変換
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.WriteError(w, "invalid user id", http.StatusBadRequest)
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
