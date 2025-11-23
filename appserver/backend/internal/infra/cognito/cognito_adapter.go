package cognito

import (
	"context"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	"github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/cognito"
)

// cognitoAdapter はCognitoClientインターフェースの実装
// 意味: pkg/aws/cognitoのラッパーとして、Repository層のインターフェースを実装
type cognitoAdapter struct {
	client cognito.Cognito
}

// NewCognitoAdapter はCognitoAdapterのコンストラクタ
// 引数:
//   - client: AWS Cognitoクライアント
// 戻り値: CognitoClientインターフェースの実装
// 実装: 依存性注入により、Cognitoクライアントを外部から受け取る
func NewCognitoAdapter(client cognito.Cognito) repository.CognitoClient {
	return &cognitoAdapter{
		client: client,
	}
}

// UpdateUserAttributes はCognitoのユーザー属性を更新する
// 引数:
//   - ctx: コンテキスト
//   - userID: 更新対象のユーザーID (Cognito Username)
//   - attributes: 更新する属性のマップ
// 戻り値: エラー情報
// 実装: pkg/aws/cognitoのAdminUpdateUserAttributesを呼び出す
// 注意事項: エラーはそのまま上位層に伝播する
func (a *cognitoAdapter) UpdateUserAttributes(ctx context.Context, userID string, attributes map[string]string) error {
	return a.client.AdminUpdateUserAttributes(ctx, userID, attributes)
}
