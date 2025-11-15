package routes

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/authapplication"
	userapplication "github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/userapplication"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/controllers"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/controllers/authcontroller"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/factory"
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

	mux.HandleFunc("POST /auth/signup", authSignupRouter)
	mux.HandleFunc("POST /auth/login", authLoginRouter)
	return mux
}

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
