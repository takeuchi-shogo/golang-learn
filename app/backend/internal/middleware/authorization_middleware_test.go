package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	jwtpkg "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
)

// TestRequireSameUser_Success は本人アクセス成功のテスト
// 実装: JWT内のユーザーIDとパスパラメータのIDが一致する場合、次のハンドラーが実行されることを検証
func TestRequireSameUser_Success(t *testing.T) {
	// テスト用のユーザーID
	testUserID := uuid.New()

	// テスト用のContextにユーザー情報を追加
	userInfo := &jwtpkg.UserInfo{
		Sub:   testUserID.String(),
		Email: "test@example.com",
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// 次のハンドラーが実行されたかを確認するフラグ
	nextHandlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", testUserID.String())
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// 次のハンドラーが実行されたか検証
	if !nextHandlerCalled {
		t.Error("Expected next handler to be called")
	}

	// ステータスコードの検証
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

// TestRequireSameUser_Forbidden_DifferentUser は他人のリソースアクセスのテスト
// 実装: JWT内のユーザーIDとパスパラメータのIDが異なる場合、403を返すことを検証
func TestRequireSameUser_Forbidden_DifferentUser(t *testing.T) {
	// テスト用のユーザーID(ログインユーザー)
	loginUserID := uuid.New()
	// 別のユーザーID(リソース所有者)
	targetUserID := uuid.New()

	// テスト用のContextにユーザー情報を追加
	userInfo := &jwtpkg.UserInfo{
		Sub:   loginUserID.String(),
		Email: "login@example.com",
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// 次のハンドラーは実行されないはず
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成(他人のID)
	req := httptest.NewRequest(http.MethodGet, "/users/"+targetUserID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", targetUserID.String())
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, rec.Code)
	}
}

// TestRequireSameUser_Unauthorized_NoUserInfoInContext はContextにユーザー情報がない場合のテスト
// 実装: Contextが空の場合、401を返すことを検証
func TestRequireSameUser_Unauthorized_NoUserInfoInContext(t *testing.T) {
	testUserID := uuid.New()

	// 次のハンドラーは実行されないはず
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成(ユーザー情報なしのContext)
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID.String(), nil)
	req.SetPathValue("id", testUserID.String())
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestRequireSameUser_BadRequest_InvalidUUID は無効なUUIDのテスト
// 実装: パスパラメータが不正なUUID形式の場合、400を返すことを検証
func TestRequireSameUser_BadRequest_InvalidUUID(t *testing.T) {
	// テスト用のユーザーID
	testUserID := uuid.New()

	// テスト用のContextにユーザー情報を追加
	userInfo := &jwtpkg.UserInfo{
		Sub:   testUserID.String(),
		Email: "test@example.com",
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// 次のハンドラーは実行されないはず
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成(不正なUUID)
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "invalid-uuid")
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestRequireSameUser_Unauthorized_InvalidUserIDInToken はToken内のUUIDが不正な場合のテスト
// 実装: JWT内のSubが不正なUUID形式の場合、401を返すことを検証
func TestRequireSameUser_Unauthorized_InvalidUserIDInToken(t *testing.T) {
	testUserID := uuid.New()

	// テスト用のContextに不正なユーザー情報を追加
	userInfo := &jwtpkg.UserInfo{
		Sub:   "invalid-uuid", // 不正なUUID
		Email: "test@example.com",
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// 次のハンドラーは実行されないはず
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/users/"+testUserID.String(), nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", testUserID.String())
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestRequireSameUser_BadRequest_EmptyPathID は空のパスIDのテスト
// 実装: パスパラメータが空の場合、400を返すことを検証
func TestRequireSameUser_BadRequest_EmptyPathID(t *testing.T) {
	// テスト用のユーザーID
	testUserID := uuid.New()

	// テスト用のContextにユーザー情報を追加
	userInfo := &jwtpkg.UserInfo{
		Sub:   testUserID.String(),
		Email: "test@example.com",
	}
	ctx := context.WithValue(context.Background(), UserInfoKey, userInfo)

	// 次のハンドラーは実行されないはず
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	// ミドルウェアを適用
	handler := RequireSameUser(nextHandler)

	// テストリクエストを作成(空のパスID)
	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	req = req.WithContext(ctx)
	req.SetPathValue("id", "")
	rec := httptest.NewRecorder()

	// リクエストを実行
	handler.ServeHTTP(rec, req)

	// ステータスコードの検証
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
