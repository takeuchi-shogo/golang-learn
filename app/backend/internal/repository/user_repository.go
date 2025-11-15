package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/domain/model"
)

// userRepository はUserQueryインターフェースの実装
type userRepository struct {
	// TODO: データベース接続やキャッシュなどの依存関係をここに追加
	db *sql.DB
	// cache cache.UserCache
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// NewUserRepository はUserRepositoryのコンストラクタ
// query.UserQueryインターフェースを実装した実体を返す
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (id, name, email, user_id_token, auth_type) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UserIDToken, model.AuthTypeCognito.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if errors.Is(err, sql.ErrConnDone) {
			return nil, errors.New("database connection is already closed")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserById は指定されたIDのユーザーをデータストアから取得する
func (r *userRepository) GetUserById(_ context.Context, id uuid.UUID) (*model.User, error) {
	// TODO: 実際のデータベースクエリやキャッシュからの取得処理を実装
	return &model.User{ID: id, Name: "test", Email: "test@example.com"}, nil
}
