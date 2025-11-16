package service

import (
	"context"
	"fmt"
	"log"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/command"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/query"
)

// UserSyncServiceWithSaga はSagaパターンを使用したユーザー同期サービス
// 意味: 分散トランザクションをSagaパターンで管理し、障害時のロールバックを実現
// 実装: userSyncServiceWithSaga構造体
type UserSyncServiceWithSaga interface {
	// SyncUserUpdate はユーザー情報をCognitoとDBに同期して更新する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 更新するユーザー情報
	// 戻り値:
	//   - *model.User: 更新されたユーザー情報
	//   - error: エラー情報
	// 実装:
	//   1. 元のユーザー情報を保存(ロールバック用)
	//   2. Cognitoのユーザー属性を更新
	//   3. DBのユーザー情報を更新
	//   4. DB更新失敗時はCognitoをロールバック
	// 注意事項: Sagaパターンにより、障害時の補償トランザクションを実行
	SyncUserUpdate(ctx context.Context, user *model.User) (*model.User, error)
}

// userSyncServiceWithSaga はUserSyncServiceWithSagaの実装
type userSyncServiceWithSaga struct {
	cognitoClient repository.CognitoClient
	userCommand   command.UserCommand
	userQuery     query.UserQuery
}

// NewUserSyncServiceWithSaga はUserSyncServiceWithSagaのコンストラクタ
// 引数:
//   - cognitoClient: Cognito操作用のクライアント
//   - userCommand: ユーザー更新用のコマンドサービス
//   - userQuery: ユーザー取得用のクエリサービス(ロールバック用)
// 戻り値: UserSyncServiceWithSagaの実装
// 実装: 依存性注入により、必要なサービスを外部から受け取る
// 注意事項: userQueryはロールバック時の元データ取得に使用
func NewUserSyncServiceWithSaga(
	cognitoClient repository.CognitoClient,
	userCommand command.UserCommand,
	userQuery query.UserQuery,
) UserSyncServiceWithSaga {
	return &userSyncServiceWithSaga{
		cognitoClient: cognitoClient,
		userCommand:   userCommand,
		userQuery:     userQuery,
	}
}

// SyncUserUpdate はユーザー情報をCognitoとDBに同期して更新する(Sagaパターン)
// 引数:
//   - ctx: コンテキスト
//   - user: 更新するユーザー情報
// 戻り値:
//   - *model.User: 更新されたユーザー情報
//   - error: エラー情報
// 実装:
//   1. 元のユーザー情報を取得(ロールバック用)
//   2. Cognitoのユーザー属性を更新
//   3. DBのユーザー情報を更新
//   4. DB更新失敗時は補償トランザクション(Cognitoロールバック)を実行
// 注意事項:
//   - 補償トランザクション失敗時はログに記録し、手動リカバリが必要
//   - べき等性を保証するため、ロールバックは元の値に戻す
func (s *userSyncServiceWithSaga) SyncUserUpdate(ctx context.Context, user *model.User) (*model.User, error) {
	// ステップ1: 元のユーザー情報を取得(ロールバック用)
	// これにより、Cognito更新後にDB更新が失敗した場合、元の状態に戻せる
	originalUser, err := s.userQuery.GetUserById(ctx, user.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to get original user for rollback: %w", err)
	}

	// ステップ2: Cognitoのユーザー属性を更新
	cognitoAttributes := map[string]string{
		"email": string(user.GetEmail()),
	}
	if err := s.cognitoClient.UpdateUserAttributes(ctx, user.GetUserIDToken(), cognitoAttributes); err != nil {
		return nil, fmt.Errorf("failed to update cognito user attributes: %w", err)
	}

	// ステップ3: DBのユーザー情報を更新
	updatedUser, err := s.userCommand.UpdateUser(ctx, user)
	if err != nil {
		// DB更新失敗 → 補償トランザクション(Cognitoロールバック)を実行
		log.Printf("[Saga] DB update failed, starting compensating transaction (rollback Cognito)")

		// 補償トランザクション: Cognitoを元の状態に戻す
		originalAttributes := map[string]string{
			"email": string(originalUser.GetEmail()),
		}
		if rollbackErr := s.cognitoClient.UpdateUserAttributes(ctx, user.GetUserIDToken(), originalAttributes); rollbackErr != nil {
			// ロールバックも失敗 → 重大なエラー
			// この場合、手動リカバリが必要
			log.Printf("[Saga] CRITICAL: Compensating transaction failed! Manual recovery required. User ID: %s, Error: %v", user.GetID(), rollbackErr)
			return nil, fmt.Errorf("failed to update user in database and rollback failed (CRITICAL - manual recovery required): db error=%w, rollback error=%v", err, rollbackErr)
		}

		log.Printf("[Saga] Compensating transaction succeeded: Cognito rolled back to original state")
		return nil, fmt.Errorf("failed to update user in database (cognito rolled back successfully): %w", err)
	}

	return updatedUser, nil
}
