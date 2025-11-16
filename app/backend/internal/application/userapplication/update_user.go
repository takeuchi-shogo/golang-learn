package userapplication

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/service"
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
	//   2. ドメインモデルを更新
	//   3. ドメインサービスで同期更新(Cognito + DB)
	// 注意事項: ユーザーが存在しない場合はエラーを返す
	Run(ctx context.Context, id uuid.UUID, name *model.Name, email *model.Email) (*model.User, error)
}

// updateUser はUpdateUserApplicationの実装
type updateUser struct {
	queryUser       query.UserQuery
	userSyncService service.UserSyncService
}

// NewUpdateUser はUpdateUserApplicationのコンストラクタ
// 引数:
//   - queryUser: ユーザー取得用のクエリサービス
//   - userSyncService: ユーザー同期用のドメインサービス
// 戻り値: UpdateUserApplicationの実装
// 実装: 依存性注入により、必要なサービスを外部から受け取る
// 注意事項: CognitoとDBの同期はUserSyncServiceに委譲
func NewUpdateUser(queryUser query.UserQuery, userSyncService service.UserSyncService) UpdateUserApplication {
	return &updateUser{
		queryUser:       queryUser,
		userSyncService: userSyncService,
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
//   2. ドメインモデルを更新(nilでない場合のみ)
//   3. UserSyncServiceで同期更新
// 注意事項:
//   - ユーザーが存在しない場合はエラーを返す
//   - Cognito更新が失敗した場合、DB更新も行われない
func (u *updateUser) Run(ctx context.Context, id uuid.UUID, name *model.Name, email *model.Email) (*model.User, error) {
	// 既存のユーザーを取得
	user, err := u.queryUser.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// ロールバック用に元のメールアドレスを保存
	oldEmail := user.GetEmail()

	// ユーザー情報を更新(nilでない場合のみ)
	if name != nil {
		user.UpdateName(*name)
	}
	if email != nil {
		user.UpdateEmail(*email)
	}

	// UserSyncServiceで同期更新(Cognito + DB)
	// ロールバック用に元のメールアドレスを渡す
	updatedUser, err := u.userSyncService.SyncUserUpdate(ctx, user, oldEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user update: %w", err)
	}

	return updatedUser, nil
}
