package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	jwtpkg "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
)

// MockJwtManager はテスト用のJWT Managerモック
// 意味: テストでJWT検証の振る舞いを制御するためのモック実装
type MockJwtManager struct {
	VerifyTokenFunc func(ctx context.Context, tokenString string) (*jwt.Token, error)
	GetUserInfoFunc func(token *jwt.Token) (*jwtpkg.UserInfo, error)
}

// VerifyToken はトークン検証のモック実装
func (m *MockJwtManager) VerifyToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	if m.VerifyTokenFunc != nil {
		return m.VerifyTokenFunc(ctx, tokenString)
	}
	return nil, errors.New("not implemented")
}

// GetUserInfo はユーザー情報取得のモック実装
func (m *MockJwtManager) GetUserInfo(token *jwt.Token) (*jwtpkg.UserInfo, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(token)
	}
	return nil, errors.New("not implemented")
}

// TestJwtVerify_Success は正常系のテスト
// 実装: 正しいトークンが渡された場合、Contextにユーザー情報が追加されることを検証
func TestJwtVerify_Success(t *testing.T) {
	// テスト用のモック設定
	mockManager := &MockJwtManager{
		VerifyTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Token, error) {
			return &jwt.Token{Valid: true}, nil
		},
		GetUserInfoFunc: func(token *jwt.Token) (*jwtpkg.UserInfo, error) {
			return &jwtpkg.UserInfo{
				Sub:   "test-user-id",
				Email: "test@example.com",
			}, nil
		},
	}

	// テスト用のハンドラー
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Contextからユーザー情報を取得
		userInfo, ok := GetUserInfoFromContext(r.Context())
		if !ok {
			t.Error("Expected user info in context, got none")
			return
		}

		// ユーザー情報の検証
		if userInfo.Sub != "test-user-id" {
			t.Errorf("Expected Sub 'test-user-id', got '%s'", userInfo.Sub)
		}
		if userInfo.Email != "test@example.com" {
			t.Errorf("Expected Email 'test@example.com', got '%s'", userInfo.Email)
		}

		w.WriteHeader(http.StatusOK)
	})

	// ミドルウェアを適用
	handler := JwtVerify(mockManager)(nextHandler)

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

// TestJwtVerify_MissingAuthHeader はAuthorizationヘッダーがない場合のテスト
// 実装: Authorizationヘッダーがない場合、401を返すことを検証
func TestJwtVerify_MissingAuthHeader(t *testing.T) {
	mockManager := &MockJwtManager{}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	handler := JwtVerify(mockManager)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestJwtVerify_InvalidBearerFormat はBearerプレフィックスがない場合のテスト
// 実装: "Bearer "プレフィックスがない場合、401を返すことを検証
func TestJwtVerify_InvalidBearerFormat(t *testing.T) {
	mockManager := &MockJwtManager{}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	handler := JwtVerify(mockManager)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestJwtVerify_InvalidToken はトークン検証が失敗する場合のテスト
// 実装: トークン検証が失敗した場合、401を返すことを検証
func TestJwtVerify_InvalidToken(t *testing.T) {
	mockManager := &MockJwtManager{
		VerifyTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Token, error) {
			return nil, errors.New("invalid token")
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	handler := JwtVerify(mockManager)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestJwtVerify_GetUserInfoFails はユーザー情報取得が失敗する場合のテスト
// 実装: ユーザー情報取得が失敗した場合、500を返すことを検証
func TestJwtVerify_GetUserInfoFails(t *testing.T) {
	mockManager := &MockJwtManager{
		VerifyTokenFunc: func(ctx context.Context, tokenString string) (*jwt.Token, error) {
			return &jwt.Token{Valid: true}, nil
		},
		GetUserInfoFunc: func(token *jwt.Token) (*jwtpkg.UserInfo, error) {
			return nil, errors.New("failed to get user info")
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	handler := JwtVerify(mockManager)(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

// TestGetUserInfoFromContext_Success は正常系のテスト
// 実装: Contextからユーザー情報を正しく取得できることを検証
func TestGetUserInfoFromContext_Success(t *testing.T) {
	expectedUserInfo := &jwtpkg.UserInfo{
		Sub:   "test-user-id",
		Email: "test@example.com",
	}

	ctx := context.WithValue(context.Background(), UserInfoKey, expectedUserInfo)

	userInfo, ok := GetUserInfoFromContext(ctx)
	if !ok {
		t.Error("Expected user info to be found")
	}
	if userInfo == nil {
		t.Fatal("Expected user info, got nil")
	}
	if userInfo.Sub != expectedUserInfo.Sub {
		t.Errorf("Expected Sub '%s', got '%s'", expectedUserInfo.Sub, userInfo.Sub)
	}
	if userInfo.Email != expectedUserInfo.Email {
		t.Errorf("Expected Email '%s', got '%s'", expectedUserInfo.Email, userInfo.Email)
	}
}

// TestGetUserInfoFromContext_NotFound はユーザー情報が存在しない場合のテスト
// 実装: Contextにユーザー情報が存在しない場合、falseを返すことを検証
func TestGetUserInfoFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	userInfo, ok := GetUserInfoFromContext(ctx)
	if ok {
		t.Error("Expected user info not to be found")
	}
	if userInfo != nil {
		t.Error("Expected nil user info")
	}
}

// TestGetUserInfoFromContext_InvalidType は型アサーションが失敗する場合のテスト
// 実装: Contextの値が不正な型の場合、falseを返すことを検証
func TestGetUserInfoFromContext_InvalidType(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserInfoKey, "invalid-type")

	userInfo, ok := GetUserInfoFromContext(ctx)
	if ok {
		t.Error("Expected user info not to be found due to type mismatch")
	}
	if userInfo != nil {
		t.Error("Expected nil user info")
	}
}
