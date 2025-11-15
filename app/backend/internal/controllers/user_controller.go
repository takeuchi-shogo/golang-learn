package controllers

import (
	"context"

	"github.com/google/uuid"
	userapplication "github.com/takeuchi-shogo/golang-learn/app/backend/internal/application/userapplication"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

// UserController はユーザー関連のHTTPリクエストを処理するコントローラー
type UserController struct {
	userApplication userapplication.UserApplication
}

// NewUserController はUserControllerのコンストラクタ
// 依存性注入により、userApplicationを外部から受け取る
func NewUserController(userApplication userapplication.UserApplication) *UserController {
	return &UserController{
		userApplication: userApplication,
	}
}

// Get は指定されたIDのユーザーを取得する
func (c *UserController) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return c.userApplication.Run(ctx, id)
}
