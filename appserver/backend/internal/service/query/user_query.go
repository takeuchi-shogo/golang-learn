package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

type UserQuery interface {
	GetUserById(context.Context, uuid.UUID) (*model.User, error)
}
