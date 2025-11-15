package userapplication

import (
	"context"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
	userquery "github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/query"
)

// UserApplication はユーザー取得のユースケースを定義するインターフェース
type UserApplication interface {
	Run(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// getUser はUserApplicationの実装
type getUser struct {
	queryUser userquery.UserQuery
}

// NewGetUser はUserApplicationのコンストラクタ
// 依存性注入により、queryUserを外部から受け取る
func NewGetUser(queryUser userquery.UserQuery) UserApplication {
	return &getUser{
		queryUser: queryUser,
	}
}

// Run は指定されたIDのユーザーを取得するユースケースを実行する
func (u *getUser) Run(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return u.queryUser.GetUserById(ctx, id)
}
