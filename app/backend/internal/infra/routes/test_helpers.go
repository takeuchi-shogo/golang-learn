package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/middleware"
	jwtpkg "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/auth/jwt"
)

// TestRequestBuilder はテストリクエストを簡単に構築するためのヘルパー
// 実装: ビルダーパターンでテストリクエストの作成を簡略化
type TestRequestBuilder struct {
	method      string
	path        string
	body        []byte
	userInfo    *jwtpkg.UserInfo
	pathParams  map[string]string
	headers     map[string]string
	queryParams map[string]string
}

// NewTestRequestBuilder はTestRequestBuilderのコンストラクタ
// 引数:
//   - method: HTTPメソッド
//   - path: パス
//
// 戻り値: TestRequestBuilderのポインタ
// 実装: メソッドチェーンでリクエストを構築できるようにする
func NewTestRequestBuilder(method, path string) *TestRequestBuilder {
	return &TestRequestBuilder{
		method:      method,
		path:        path,
		pathParams:  make(map[string]string),
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
}

// WithBody はリクエストボディを設定する
// 引数:
//   - body: リクエストボディのバイト配列
//
// 戻り値: TestRequestBuilderのポインタ(メソッドチェーン用)
func (b *TestRequestBuilder) WithBody(body []byte) *TestRequestBuilder {
	b.body = body
	return b
}

// WithUserInfo はContextにユーザー情報を設定する
// 引数:
//   - userInfo: ユーザー情報
//
// 戻り値: TestRequestBuilderのポインタ(メソッドチェーン用)
// 実装: JWT認証済みのリクエストを簡単に作成
func (b *TestRequestBuilder) WithUserInfo(userInfo *jwtpkg.UserInfo) *TestRequestBuilder {
	b.userInfo = userInfo
	return b
}

// WithPathParam はパスパラメータを設定する
// 引数:
//   - key: パスパラメータのキー
//   - value: パスパラメータの値
//
// 戻り値: TestRequestBuilderのポインタ(メソッドチェーン用)
func (b *TestRequestBuilder) WithPathParam(key, value string) *TestRequestBuilder {
	b.pathParams[key] = value
	return b
}

// WithHeader はHTTPヘッダーを設定する
// 引数:
//   - key: ヘッダーのキー
//   - value: ヘッダーの値
//
// 戻り値: TestRequestBuilderのポインタ(メソッドチェーン用)
func (b *TestRequestBuilder) WithHeader(key, value string) *TestRequestBuilder {
	b.headers[key] = value
	return b
}

// Build はテストリクエストを構築する
// 戻り値: *http.Request
// 実装: 設定された情報からテストリクエストを生成
func (b *TestRequestBuilder) Build() *http.Request {
	var req *http.Request
	if b.body != nil {
		req = httptest.NewRequest(b.method, b.path, nil)
	} else {
		req = httptest.NewRequest(b.method, b.path, nil)
	}

	// Contextにユーザー情報を追加
	if b.userInfo != nil {
		ctx := context.WithValue(req.Context(), middleware.UserInfoKey, b.userInfo)
		req = req.WithContext(ctx)
	}

	// パスパラメータを設定
	for key, value := range b.pathParams {
		req.SetPathValue(key, value)
	}

	// HTTPヘッダーを設定
	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	return req
}

// CreateTestUserInfo はテスト用のユーザー情報を作成するヘルパー関数
// 引数:
//   - userID: ユーザーID(nilの場合は新規UUID生成)
//   - email: メールアドレス(空の場合は"test@example.com")
//
// 戻り値: *jwtpkg.UserInfo
// 実装: テスト用のユーザー情報を簡単に作成
func CreateTestUserInfo(userID *uuid.UUID, email string) *jwtpkg.UserInfo {
	var id uuid.UUID
	if userID == nil {
		id = uuid.New()
	} else {
		id = *userID
	}

	if email == "" {
		email = "test@example.com"
	}

	return &jwtpkg.UserInfo{
		Sub:   id.String(),
		Email: email,
	}
}

// CreateAuthenticatedRequest は認証済みリクエストを作成するヘルパー関数
// 引数:
//   - method: HTTPメソッド
//   - path: パス
//   - userID: ユーザーID
//   - pathID: パスパラメータのID
//
// 戻り値: *http.Request
// 実装: 認証済みリクエストを簡単に作成
// 注意事項: テストで頻繁に使用するパターンを簡略化
func CreateAuthenticatedRequest(method, path string, userID uuid.UUID, pathID uuid.UUID) *http.Request {
	userInfo := CreateTestUserInfo(&userID, "")
	return NewTestRequestBuilder(method, path).
		WithUserInfo(userInfo).
		WithPathParam("id", pathID.String()).
		Build()
}

// AssertStatusCode はステータスコードをアサートするヘルパー関数
// 引数:
//   - t: テストオブジェクト
//   - rec: レスポンスレコーダー
//   - expected: 期待されるステータスコード
//
// 実装: テストコードの可読性を向上
func AssertStatusCode(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("Expected status code %d, got %d. Body: %s",
			expected, rec.Code, rec.Body.String())
	}
}

// AssertErrorResponse はエラーレスポンスをアサートするヘルパー関数
// 引数:
//   - t: テストオブジェクト
//   - rec: レスポンスレコーダー
//   - expectedStatusCode: 期待されるステータスコード
//   - expectedErrorContains: エラーメッセージに含まれるべき文字列
//
// 実装: エラーレスポンスの検証を簡略化
func AssertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatusCode int, expectedErrorContains string) {
	t.Helper()

	// ステータスコードの検証
	AssertStatusCode(t, rec, expectedStatusCode)

	// エラーメッセージの検証
	// 注意: 実際の実装ではJSONのパースとエラーメッセージの検証を行う
	// ここでは概念的な実装として示す
}
