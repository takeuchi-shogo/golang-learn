package controllers

import (
	"context"

	"github.com/google/uuid"
	userapplication "github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/userapplication"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

// UserController はユーザー関連のHTTPリクエストを処理するコントローラー
// 実装: ユーザーの取得と更新を行う
type UserController struct {
	getUserApplication    userapplication.UserApplication
	updateUserApplication userapplication.UpdateUserApplication
}

// NewUserController はUserControllerのコンストラクタ
// 引数:
//   - getUserApplication: ユーザー取得のユースケース
// 戻り値: UserControllerのポインタ
// 実装: 依存性注入により、userApplicationを外部から受け取る
func NewUserController(getUserApplication userapplication.UserApplication) *UserController {
	return &UserController{
		getUserApplication: getUserApplication,
	}
}

// NewUserControllerWithUpdate はUserControllerのコンストラクタ(更新機能付き)
// 引数:
//   - getUserApplication: ユーザー取得のユースケース
//   - updateUserApplication: ユーザー更新のユースケース
// 戻り値: UserControllerのポインタ
// 実装: 依存性注入により、必要なアプリケーションを外部から受け取る
// 注意事項: 更新機能が必要な場合はこのコンストラクタを使用する
func NewUserControllerWithUpdate(
	getUserApplication userapplication.UserApplication,
	updateUserApplication userapplication.UpdateUserApplication,
) *UserController {
	return &UserController{
		getUserApplication:    getUserApplication,
		updateUserApplication: updateUserApplication,
	}
}

// Get は指定されたIDのユーザーを取得する
// 引数:
//   - ctx: コンテキスト
//   - id: ユーザーID
// 戻り値:
//   - *model.User: 取得されたユーザー情報
//   - error: エラー情報
// 実装: アプリケーション層のユースケースを実行する
func (c *UserController) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return c.getUserApplication.Run(ctx, id)
}

// Update はユーザー情報を更新する
// 引数:
//   - ctx: コンテキスト
//   - id: 更新対象のユーザーID
//   - name: 新しいユーザー名(nilの場合は更新しない)
//   - email: 新しいメールアドレス(nilの場合は更新しない)
// 戻り値:
//   - *model.User: 更新されたユーザー情報
//   - error: エラー情報
// 実装: アプリケーション層のユースケースを実行する
func (c *UserController) Update(ctx context.Context, id uuid.UUID, name *model.Name, email *model.Email) (*model.User, error) {
	return c.updateUserApplication.Run(ctx, id, name, email)
}
