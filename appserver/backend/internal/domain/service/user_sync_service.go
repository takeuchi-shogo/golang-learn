package service

import (
	"context"
	"fmt"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/command"
)

// UserSyncService はユーザー情報をCognitoとDBで同期するドメインサービスインターフェース
// 意味: ユーザー情報の同期に関するビジネスルールを定義
// 実装: userSyncService構造体
type UserSyncService interface {
	// SyncUserUpdate はユーザー情報をCognitoとDBに同期して更新する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 更新するユーザー情報
	//   - oldEmail: 元のメールアドレス(ロールバック用)
	// 戻り値:
	//   - *model.User: 更新されたユーザー情報
	//   - error: エラー情報
	// 実装:
	//   1. Cognito のユーザー属性を更新
	//   2. DB のユーザー情報を更新
	//   3. DB更新失敗時はCognitoをロールバック
	// 注意事項:
	//   - Cognito更新が失敗した場合、DB更新は行わない
	//   - DB更新が失敗した場合、Cognitoを元の状態にロールバック
	//   - これによりCognitoとDBの整合性を保つ
	SyncUserUpdate(ctx context.Context, user *model.User, oldEmail model.Email) (*model.User, error)
}

// userSyncService はUserSyncServiceの実装
type userSyncService struct {
	cognitoClient repository.CognitoClient
	userCommand   command.UserCommand
}

// NewUserSyncService はUserSyncServiceのコンストラクタ
// 引数:
//   - cognitoClient: Cognito操作用のクライアント
//   - userCommand: ユーザー更新用のコマンドサービス
// 戻り値: UserSyncServiceの実装
// 実装: 依存性注入により、必要なサービスを外部から受け取る
func NewUserSyncService(cognitoClient repository.CognitoClient, userCommand command.UserCommand) UserSyncService {
	return &userSyncService{
		cognitoClient: cognitoClient,
		userCommand:   userCommand,
	}
}

// SyncUserUpdate はユーザー情報をCognitoとDBに同期して更新する
// 引数:
//   - ctx: コンテキスト
//   - user: 更新するユーザー情報
//   - oldEmail: 元のメールアドレス(ロールバック用)
// 戻り値:
//   - *model.User: 更新されたユーザー情報
//   - error: エラー情報
// 実装:
//   1. ユーザー情報からCognito更新用の属性マップを作成
//   2. Cognitoのユーザー属性を更新
//   3. DBのユーザー情報を更新
//   4. DB更新失敗時はCognitoをロールバック(Saga Pattern)
// 注意事項:
//   - Cognito更新が失敗した場合、DB更新は行わない(整合性維持)
//   - DB更新が失敗した場合、Cognitoを元の状態にロールバック
//   - ロールバック失敗時はCRITICALエラーとして記録
//   - 現在はemailのみCognitoに同期(nameは標準属性にないためDBのみ)
func (s *userSyncService) SyncUserUpdate(ctx context.Context, user *model.User, oldEmail model.Email) (*model.User, error) {
	// Cognito更新用の属性マップを作成
	// 注意: 現在の実装ではemailのみをCognitoに同期
	// nameはCognitoの標準属性にないため、カスタム属性として保存する場合は別途実装が必要
	cognitoAttributes := map[string]string{
		"email": string(user.GetEmail()),
	}

	// Cognitoのユーザー属性を更新
	// ビジネスルール: Cognito更新が失敗した場合、DB更新も行わない
	if err := s.cognitoClient.UpdateUserAttributes(ctx, user.GetUserIDToken(), cognitoAttributes); err != nil {
		return nil, fmt.Errorf("failed to update cognito user attributes: %w", err)
	}

	// DBのユーザー情報を更新
	updatedUser, err := s.userCommand.UpdateUser(ctx, user)
	if err != nil {
		// DB更新失敗時、Cognitoをロールバック(Saga Pattern)
		rollbackAttributes := map[string]string{
			"email": string(oldEmail),
		}

		if rollbackErr := s.cognitoClient.UpdateUserAttributes(ctx, user.GetUserIDToken(), rollbackAttributes); rollbackErr != nil {
			// ロールバック失敗は致命的エラー
			// この場合、CognitoとDBでデータ不整合が発生している
			return nil, fmt.Errorf("CRITICAL: failed to rollback cognito update - data inconsistency may occur (original error: %w, rollback error: %v)", err, rollbackErr)
		}

		// ロールバック成功時は通常のエラーを返す
		return nil, fmt.Errorf("failed to update user in database (cognito rolled back successfully): %w", err)
	}

	return updatedUser, nil
}
