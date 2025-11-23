package repository

import "context"

// CognitoClient はCognito操作を行うインターフェース
// 意味: Infrastructure層のCognitoクライアントを抽象化し、Domain層からの依存を防ぐ
// 実装: internal/infra/cognito/cognito_adapter.go
type CognitoClient interface {
	// UpdateUserAttributes はCognitoのユーザー属性を更新する
	// 引数:
	//   - ctx: コンテキスト
	//   - userID: 更新対象のユーザーID (Cognito Username)
	//   - attributes: 更新する属性のマップ (例: {"email": "new@example.com"})
	// 戻り値: エラー情報
	// 実装: CognitoのAdminUpdateUserAttributes APIを使用
	// 注意事項:
	//   - ユーザーが存在しない場合はエラーを返す
	//   - email更新時は自動的にemail_verifiedもtrueに設定される
	UpdateUserAttributes(ctx context.Context, userID string, attributes map[string]string) error
}
