# 認可ミドルウェア テストケース仕様

## 概要

本人確認を行う認可ミドルウェア(`RequireSameUser`)の包括的なテストケース仕様書

## テスト対象

- **ファイル**: `authorization_middleware.go`
- **関数**: `RequireSameUser(next http.Handler) http.Handler`

## テストファイル

- **ファイル**: `authorization_middleware_test.go`

---

## テストケース一覧

### 1. 正常系テスト

#### 1.1 本人アクセス成功

- **テスト名**: `TestRequireSameUser_Success`
- **目的**: JWT内のユーザーIDとパスパラメータのIDが一致する場合、次のハンドラーが実行されることを検証
- **入力**:
  - ユーザーID: `uuid.New()`
  - JWT UserInfo.Sub: 同じユーザーID
  - パスパラメータ "id": 同じユーザーID
- **期待される出力**:
  - ステータスコード: `200 OK`
  - 次のハンドラーが実行される
- **検証項目**:
  - [x] 次のハンドラーが実行されたか
  - [x] ステータスコードが200か

---

### 2. 異常系テスト

#### 2.1 他人のリソースアクセス

- **テスト名**: `TestRequireSameUser_Forbidden_DifferentUser`
- **目的**: JWT内のユーザーIDとパスパラメータのIDが異なる場合、403を返すことを検証
- **入力**:
  - ログインユーザーID: `uuid.New()`
  - JWT UserInfo.Sub: ログインユーザーID
  - パスパラメータ "id": 別のユーザーID
- **期待される出力**:
  - ステータスコード: `403 Forbidden`
  - エラーメッセージ: "forbidden: cannot access other user's information"
- **検証項目**:
  - [x] ステータスコードが403か
  - [x] 次のハンドラーが実行されないか
  - [x] エラーメッセージが適切か

#### 2.2 Contextにユーザー情報なし

- **テスト名**: `TestRequireSameUser_Unauthorized_NoUserInfoInContext`
- **目的**: Contextが空の場合、401を返すことを検証
- **入力**:
  - Context: 空のContext
  - パスパラメータ "id": 有効なUUID
- **期待される出力**:
  - ステータスコード: `401 Unauthorized`
  - エラーメッセージ: "user info not found in context"
- **検証項目**:
  - [x] ステータスコードが401か
  - [x] 次のハンドラーが実行されないか

#### 2.3 無効なUUID(パスパラメータ)

- **テスト名**: `TestRequireSameUser_BadRequest_InvalidUUID`
- **目的**: パスパラメータが不正なUUID形式の場合、400を返すことを検証
- **入力**:
  - JWT UserInfo.Sub: 有効なUUID
  - パスパラメータ "id": "invalid-uuid"
- **期待される出力**:
  - ステータスコード: `400 Bad Request`
  - エラーメッセージ: "invalid user id"
- **検証項目**:
  - [x] ステータスコードが400か
  - [x] 次のハンドラーが実行されないか

#### 2.4 無効なUUID(Token内)

- **テスト名**: `TestRequireSameUser_Unauthorized_InvalidUserIDInToken`
- **目的**: JWT内のSubが不正なUUID形式の場合、401を返すことを検証
- **入力**:
  - JWT UserInfo.Sub: "invalid-uuid"
  - パスパラメータ "id": 有効なUUID
- **期待される出力**:
  - ステータスコード: `401 Unauthorized`
  - エラーメッセージ: "invalid user id in token"
- **検証項目**:
  - [x] ステータスコードが401か
  - [x] 次のハンドラーが実行されないか

#### 2.5 空のパスID

- **テスト名**: `TestRequireSameUser_BadRequest_EmptyPathID`
- **目的**: パスパラメータが空の場合、400を返すことを検証
- **入力**:
  - JWT UserInfo.Sub: 有効なUUID
  - パスパラメータ "id": "" (空文字列)
- **期待される出力**:
  - ステータスコード: `400 Bad Request`
  - エラーメッセージ: "user id is required"
- **検証項目**:
  - [x] ステータスコードが400か
  - [x] 次のハンドラーが実行されないか

---

## エッジケースのカバレッジ

### カバーされているエッジケース

1. ✅ UUID parseエラー(パスパラメータ)
2. ✅ UUID parseエラー(Token内のSub)
3. ✅ Contextにユーザー情報なし
4. ✅ 空のパスパラメータ
5. ✅ ユーザーID不一致

### 境界値テスト

- UUID形式: 有効/無効なUUID文字列
- ユーザーID比較: 一致/不一致

---

## テスト実行方法

```bash
# 認可ミドルウェアのテストのみ実行
go test -v ./internal/middleware -run TestRequireSameUser

# カバレッジ計測
go test -cover ./internal/middleware

# 詳細なカバレッジレポート
go test -coverprofile=coverage.out ./internal/middleware
go tool cover -html=coverage.out
```

---

## モック設計

### 使用するモック

- なし(ミドルウェア単体のテストのため、モックは不要)

### Context設定

- `UserInfoKey`: ユーザー情報をContextに保存するためのキー
- `jwtpkg.UserInfo`: JWT認証済みのユーザー情報

---

## 注意事項

1. **前提条件**:
   - このミドルウェアは`JwtVerify`ミドルウェアの後に適用する必要がある
   - Contextに`UserInfo`が存在することが前提

2. **テストの独立性**:
   - 各テストケースは独立して実行可能
   - テスト間で状態を共有しない

3. **エラーハンドリング**:
   - すべてのエラーケースで適切なHTTPステータスコードを返す
   - エラーメッセージは統一されたフォーマット(httputil.WriteError)を使用

---

## カバレッジ目標

- **行カバレッジ**: 100%
- **分岐カバレッジ**: 100%

すべてのコードパス(正常系、異常系、エッジケース)をカバーするテストケースを作成しています。
