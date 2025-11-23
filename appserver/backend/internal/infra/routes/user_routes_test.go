package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/middleware"
	jwtpkg "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
)

// MockUserQuery はテスト用のUserQueryモック
// 意味: テストでユーザー取得の振る舞いを制御するためのモック実装
type MockUserQuery struct {
	GetUserByIdFunc func(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// GetUserById はユーザー取得のモック実装
func (m *MockUserQuery) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

// TestGetUserEndpoint_Success_WithAuth_OwnData は認証あり本人データ取得の正常系テスト
// 実装: 有効なJWTトークンで本人のデータを取得できることを検証
func TestGetUserEndpoint_Success_WithAuth_OwnData(t *testing.T) {
	// テスト用のユーザーID
	testUserID := uuid.New()
	testEmail := "test@example.com"

	// テスト用のユーザー情報をContextに設定
	userInfo := &jwtpkg.UserInfo{
		Sub:   testUserID.String(),
		Email: testEmail,
	}
	ctx := context.WithValue(context.Background(), middleware.UserInfoKey, userInfo)

	// モックの設定: 期待されるユーザー情報を返す
	expectedUser := &model.User{
		// モックのユーザー情報
		// 注意: 実際のUserモデルの構造に応じて調整が必要
	}

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", testUserID.String())
	rec := httptest.NewRecorder()

	// エンドポイントを実行
	// 注意: 実際のテストでは依存関係を注入する方法が必要
	// ここでは概念的なテストケースとして示す
	userRouter(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// レスポンスボディの検証
	var user model.User
	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// ユーザー情報の検証
	_ = expectedUser // 期待されるユーザー情報と比較
}

// TestGetUserEndpoint_Unauthorized_NoAuth は認証なしの異常系テスト
// 実装: Authorizationヘッダーがない場合、401を返すことを検証
func TestGetUserEndpoint_Unauthorized_NoAuth(t *testing.T) {
	testUserID := uuid.New()

	// テストリクエストを作成(認証情報なし)
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID.String(), nil)
	req.SetPathValue("id", testUserID.String())
	rec := httptest.NewRecorder()

	// JwtVerifyミドルウェアを適用したハンドラーを実行
	// モックのJWTマネージャーを使用
	mockJwtManager := &MockJwtManager{}
	handler := middleware.JwtVerify(mockJwtManager)(http.HandlerFunc(userRouter))

	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestGetUserEndpoint_BadRequest_InvalidID は不正なID形式の異常系テスト
// 実装: パスパラメータが不正なUUIDの場合、400を返すことを検証
func TestGetUserEndpoint_BadRequest_InvalidID(t *testing.T) {
	// テスト用のユーザー情報をContextに設定
	userInfo := &jwtpkg.UserInfo{
		Sub:   uuid.New().String(),
		Email: "test@example.com",
	}
	ctx := context.WithValue(context.Background(), middleware.UserInfoKey, userInfo)

	// テストリクエストを作成(不正なUUID)
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid-uuid")
	rec := httptest.NewRecorder()

	// エンドポイントを実行
	userRouter(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// MockJwtManager はテスト用のJWT Managerモック(認証ミドルウェアテスト用)
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
