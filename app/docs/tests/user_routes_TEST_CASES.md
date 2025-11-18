# エンドポイント統合テスト仕様

## 概要

ユーザー関連エンドポイント(GET/PUT /users/{id})の統合テスト仕様書

## テスト対象エンドポイント

1. **GET /users/{id}** - ユーザー情報取得
2. **PUT /users/{id}** - ユーザー情報更新

## テストファイル

- `user_routes_test.go` - GET /users/{id} のテスト
- `update_user_routes_test.go` - PUT /users/{id} のテスト
- `test_helpers.go` - テストヘルパー関数

---

# 1. GET /users/{id} エンドポイント

## テストケース一覧

### 1.1 正常系テスト

#### 1.1.1 認証あり本人データ取得

- **テスト名**: `TestGetUserEndpoint_Success_WithAuth_OwnData`
- **目的**: 有効なJWTトークンで本人のデータを取得できることを検証
- **前提条件**:
  - 認証済み(Contextにユーザー情報あり)
  - パスパラメータのIDとログインユーザーIDが一致
- **入力**:
  - HTTPメソッド: GET
  - パス: `/users/{id}`
  - ユーザーID: `uuid.New()`
  - Context: 認証済みユーザー情報
- **期待される出力**:
  - ステータスコード: `200 OK`
  - レスポンスボディ: ユーザー情報のJSON
- **検証項目**:
  - [x] ステータスコードが200か
  - [x] レスポンスボディが正しいユーザー情報を含むか

---

### 1.2 異常系テスト

#### 1.2.1 認証なし

- **テスト名**: `TestGetUserEndpoint_Unauthorized_NoAuth`
- **目的**: Authorizationヘッダーがない場合、401を返すことを検証
- **入力**:
  - HTTPメソッド: GET
  - パス: `/users/{id}`
  - Authorizationヘッダー: なし
- **期待される出力**:
  - ステータスコード: `401 Unauthorized`
  - エラーメッセージ: "Authorization header is required"
- **検証項目**:
  - [x] ステータスコードが401か
  - [x] エラーレスポンスが返されるか

#### 1.2.2 不正なID形式

- **テスト名**: `TestGetUserEndpoint_BadRequest_InvalidID`
- **目的**: パスパラメータが不正なUUIDの場合、400を返すことを検証
- **入力**:
  - HTTPメソッド: GET
  - パス: `/users/invalid-uuid`
  - パスパラメータ "id": "invalid-uuid"
- **期待される出力**:
  - ステータスコード: `400 Bad Request`
  - エラーメッセージ: "invalid user id"
- **検証項目**:
  - [x] ステータスコードが400か
  - [x] エラーレスポンスが返されるか

---

# 2. PUT /users/{id} エンドポイント

## テストケース一覧

### 2.1 正常系テスト

#### 2.1.1 認証あり本人データ更新

- **テスト名**: `TestUpdateUserEndpoint_Success_WithAuth_OwnData`
- **目的**: 有効なJWTトークンで本人のデータを更新できることを検証
- **前提条件**:
  - 認証済み(Contextにユーザー情報あり)
  - パスパラメータのIDとログインユーザーIDが一致
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/{id}`
  - ユーザーID: `uuid.New()`
  - リクエストボディ:

    ```json
    {
      "name": "Updated Name"
    }
    ```

- **期待される出力**:
  - ステータスコード: `200 OK`
  - レスポンスボディ: 更新されたユーザー情報のJSON
- **検証項目**:
  - [x] ステータスコードが200か
  - [x] レスポンスボディが更新されたユーザー情報を含むか
  - [x] リファクタリング後も既存の動作が保持されるか

---

### 2.2 異常系テスト

#### 2.2.1 認証あり他人データ更新試行

- **テスト名**: `TestUpdateUserEndpoint_Forbidden_OthersData`
- **目的**: 有効なJWTトークンで他人のデータ更新を試みた場合、403を返すことを検証
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/{targetUserId}`
  - ログインユーザーID: `uuid.New()`
  - ターゲットユーザーID: 別の`uuid.New()`
  - リクエストボディ:

    ```json
    {
      "name": "Hacked Name"
    }
    ```

- **期待される出力**:
  - ステータスコード: `403 Forbidden`
  - エラーメッセージ: "forbidden: cannot update other user's information"
- **検証項目**:
  - [x] ステータスコードが403か
  - [x] エラーメッセージに"forbidden"が含まれるか

#### 2.2.2 認証なし

- **テスト名**: `TestUpdateUserEndpoint_Unauthorized_NoAuth`
- **目的**: Authorizationヘッダーがない場合、401を返すことを検証
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/{id}`
  - Authorizationヘッダー: なし
- **期待される出力**:
  - ステータスコード: `401 Unauthorized`
- **検証項目**:
  - [x] ステータスコードが401か

#### 2.2.3 不正なID形式

- **テスト名**: `TestUpdateUserEndpoint_BadRequest_InvalidID`
- **目的**: パスパラメータが不正なUUIDの場合、400を返すことを検証
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/invalid-uuid`
  - パスパラメータ "id": "invalid-uuid"
- **期待される出力**:
  - ステータスコード: `400 Bad Request`
- **検証項目**:
  - [x] ステータスコードが400か

#### 2.2.4 不正なリクエストボディ

- **テスト名**: `TestUpdateUserEndpoint_BadRequest_InvalidRequestBody`
- **目的**: JSONのパースに失敗した場合、400を返すことを検証
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/{id}`
  - リクエストボディ: `{"name": invalid}` (不正なJSON)
- **期待される出力**:
  - ステータスコード: `400 Bad Request`
  - エラーメッセージ: "invalid request body"
- **検証項目**:
  - [x] ステータスコードが400か

#### 2.2.5 バリデーションエラー

- **テスト名**: `TestUpdateUserEndpoint_BadRequest_ValidationError`
- **目的**: バリデーションに失敗した場合、400を返すことを検証
- **入力**:
  - HTTPメソッド: PUT
  - パス: `/users/{id}`
  - リクエストボディ:

    ```json
    {
      "name": ""
    }
    ```

- **期待される出力**:
  - ステータスコード: `400 Bad Request`
  - エラーメッセージ: "name cannot be empty"
- **検証項目**:
  - [x] ステータスコードが400か
  - [x] バリデーションエラーメッセージが適切か

---

# 3. テストヘルパー関数

## 3.1 TestRequestBuilder

ビルダーパターンでテストリクエストを構築

### メソッド

- `NewTestRequestBuilder(method, path string)`: コンストラクタ
- `WithBody(body []byte)`: リクエストボディを設定
- `WithUserInfo(userInfo *jwtpkg.UserInfo)`: Contextにユーザー情報を設定
- `WithPathParam(key, value string)`: パスパラメータを設定
- `WithHeader(key, value string)`: HTTPヘッダーを設定
- `Build()`: テストリクエストを構築

### 使用例

```go
req := NewTestRequestBuilder(http.MethodPut, "/users/123").
    WithUserInfo(userInfo).
    WithPathParam("id", "123").
    WithHeader("Content-Type", "application/json").
    Build()
```

## 3.2 CreateTestUserInfo

テスト用のユーザー情報を作成

### 引数

- `userID *uuid.UUID`: ユーザーID(nilの場合は新規生成)
- `email string`: メールアドレス(空の場合はデフォルト値)

### 使用例

```go
userInfo := CreateTestUserInfo(nil, "test@example.com")
```

## 3.3 CreateAuthenticatedRequest

認証済みリクエストを作成

### 引数

- `method string`: HTTPメソッド
- `path string`: パス
- `userID uuid.UUID`: ログインユーザーID
- `pathID uuid.UUID`: パスパラメータのID

### 使用例

```go
req := CreateAuthenticatedRequest(http.MethodGet, "/users/123", userID, pathID)
```

## 3.4 AssertStatusCode

ステータスコードをアサート

### 引数

- `t *testing.T`: テストオブジェクト
- `rec *httptest.ResponseRecorder`: レスポンスレコーダー
- `expected int`: 期待されるステータスコード

### 使用例

```go
AssertStatusCode(t, rec, http.StatusOK)
```

## 3.5 AssertErrorResponse

エラーレスポンスをアサート

### 引数

- `t *testing.T`: テストオブジェクト
- `rec *httptest.ResponseRecorder`: レスポンスレコーダー
- `expectedStatusCode int`: 期待されるステータスコード
- `expectedErrorContains string`: エラーメッセージに含まれるべき文字列

### 使用例

```go
AssertErrorResponse(t, rec, http.StatusForbidden, "forbidden")
```

---

# 4. モック設計

## 4.1 MockUserQuery

ユーザー取得のクエリサービスをモック

### メソッド

- `GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)`

### 実装

```go
type MockUserQuery struct {
    GetUserByIdFunc func(ctx context.Context, id uuid.UUID) (*model.User, error)
}
```

## 4.2 MockJwtManager

JWT検証サービスをモック

### メソッド

- `VerifyToken(ctx context.Context, tokenString string) (interface{}, error)`
- `GetUserInfo(token interface{}) (*jwtpkg.UserInfo, error)`

### 実装

```go
type MockJwtManager struct {
    VerifyTokenFunc func(ctx context.Context, tokenString string) (interface{}, error)
    GetUserInfoFunc func(token interface{}) (*jwtpkg.UserInfo, error)
}
```

---

# 5. テスト実行方法

## 5.1 全テスト実行

```bash
# routesパッケージのすべてのテストを実行
go test -v ./internal/infra/routes
```

## 5.2 特定のテストのみ実行

```bash
# GET /users/{id} のテストのみ実行
go test -v ./internal/infra/routes -run TestGetUserEndpoint

# PUT /users/{id} のテストのみ実行
go test -v ./internal/infra/routes -run TestUpdateUserEndpoint
```

## 5.3 カバレッジ計測

```bash
# カバレッジ計測
go test -cover ./internal/infra/routes

# 詳細なカバレッジレポート
go test -coverprofile=coverage.out ./internal/infra/routes
go tool cover -html=coverage.out
```

---

# 6. エッジケースのカバレッジ

## カバーされているエッジケース

1. ✅ UUID parseエラー(パスパラメータ)
2. ✅ UUID parseエラー(Token内のSub)
3. ✅ Contextにユーザー情報なし
4. ✅ 認証ヘッダーなし
5. ✅ 不正なJSONリクエストボディ
6. ✅ バリデーションエラー
7. ✅ 本人以外のデータアクセス
8. ✅ 空のリクエストフィールド

## 境界値テスト

- UUID形式: 有効/無効なUUID文字列
- ユーザーID比較: 一致/不一致
- リクエストボディ: 有効/無効なJSON
- バリデーション: 空文字列、長すぎる文字列、不正なメール形式

---

# 7. 注意事項

## 7.1 テストの前提条件

- **GET /users/{id}**: 現在は認可ミドルウェアなし(将来的に追加予定)
- **PUT /users/{id}**: JwtVerifyミドルウェアが必須

## 7.2 テストの独立性

- 各テストケースは独立して実行可能
- テスト間で状態を共有しない
- モックを使用して外部依存を排除

## 7.3 リファクタリング後の動作確認

PUT /users/{id}のテストは、リファクタリング後も既存の動作が保持されることを確認する目的も持つ

---

# 8. カバレッジ目標

- **行カバレッジ**: 80%以上
- **分岐カバレッジ**: 80%以上

統合テストのため、すべての外部依存をモック化することで、安定したテストを実現します。
