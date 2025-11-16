package userapplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/command"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/query"
)

// UpdateUserApplication はユーザー更新のユースケースを定義するインターフェース
// 実装: updateUser構造体
type UpdateUserApplication interface {
	// Run はユーザー情報を更新するユースケースを実行する
	// 引数:
	//   - ctx: コンテキスト
	//   - id: 更新対象のユーザーID
	//   - name: 新しいユーザー名(nilの場合は更新しない)
	//   - email: 新しいメールアドレス(nilの場合は更新しない)
	// 戻り値:
	//   - *model.User: 更新されたユーザー情報
	//   - error: エラー情報
	// 実装:
	//   1. 既存のユーザーを取得
	//   2. ユーザー情報を更新
	//   3. 更新されたユーザーを保存
	// 注意事項: ユーザーが存在しない場合はエラーを返す
	Run(ctx context.Context, id uuid.UUID, name *model.Name, email *model.Email) (*model.User, error)
}

// updateUser はUpdateUserApplicationの実装
type updateUser struct {
	queryUser   query.UserQuery
	commandUser command.UserCommand
}

// NewUpdateUser はUpdateUserApplicationのコンストラクタ
// 引数:
//   - queryUser: ユーザー取得用のクエリサービス
//   - commandUser: ユーザー更新用のコマンドサービス
// 戻り値: UpdateUserApplicationの実装
// 実装: 依存性注入により、必要なサービスを外部から受け取る
func NewUpdateUser(queryUser query.UserQuery, commandUser command.UserCommand) UpdateUserApplication {
	return &updateUser{
		queryUser:   queryUser,
		commandUser: commandUser,
	}
}

// Run はユーザー情報を更新するユースケースを実行する
// 引数:
//   - ctx: コンテキスト
//   - id: 更新対象のユーザーID
//   - name: 新しいユーザー名(nilの場合は更新しない)
//   - email: 新しいメールアドレス(nilの場合は更新しない)
// 戻り値:
//   - *model.User: 更新されたユーザー情報
//   - error: エラー情報
// 実装:
//   1. 既存のユーザーを取得
//   2. ユーザー情報を更新
//   3. 更新されたユーザーを保存
// 注意事項: ユーザーが存在しない場合はエラーを返す
func (u *updateUser) Run(ctx context.Context, id uuid.UUID, name *model.Name, email *model.Email) (*model.User, error) {
	// 既存のユーザーを取得
	user, err := u.queryUser.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// ユーザー情報を更新(nilでない場合のみ)
	if name != nil {
		user.UpdateName(*name)
	}
	if email != nil {
		user.UpdateEmail(*email)
	}

	// 更新されたユーザーを保存
	updatedUser, err := u.commandUser.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
