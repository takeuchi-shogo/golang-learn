package command

import (
	"context"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

type UserCommand interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
}
