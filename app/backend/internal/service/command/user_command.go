package command

import (
	"context"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

// UserCommand はユーザーの書き込み操作を定義するインターフェース
// 実装: リポジトリ層で実装される
type UserCommand interface {
	// CreateUser はユーザーを新規作成する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 作成するユーザー情報
	// 戻り値:
	//   - *model.User: 作成されたユーザー情報
	//   - error: エラー情報
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// UpdateUser はユーザー情報を更新する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 更新するユーザー情報
	// 戻り値:
	//   - *model.User: 更新されたユーザー情報
	//   - error: エラー情報
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
}
