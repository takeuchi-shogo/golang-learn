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

// UserRepository はユーザーのデータ永続化を行うインターフェース
// 実装: user_repository構造体
type UserRepository interface {
	// CreateUser はユーザーを新規作成する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 作成するユーザー情報
	// 戻り値:
	//   - *model.User: 作成されたユーザー情報
	//   - error: エラー情報
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// GetUserById は指定されたIDのユーザーを取得する
	// 引数:
	//   - ctx: コンテキスト
	//   - id: ユーザーID
	// 戻り値:
	//   - *model.User: 取得されたユーザー情報
	//   - error: エラー情報
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error)

	// UpdateUser はユーザー情報を更新する
	// 引数:
	//   - ctx: コンテキスト
	//   - user: 更新するユーザー情報
	// 戻り値:
	//   - *model.User: 更新されたユーザー情報
	//   - error: エラー情報
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
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
// 引数:
//   - ctx: コンテキスト
//   - id: ユーザーID
// 戻り値:
//   - *model.User: 取得されたユーザー情報
//   - error: エラー情報
// 実装: データベースから指定されたIDのユーザーを取得する
// 注意事項: ユーザーが見つからない場合はnilとエラーを返す
func (r *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := "SELECT id, name, email, user_id_token FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var user model.User
	var name string
	var email string
	var userIDToken string

	err := row.Scan(&user.ID, &name, &email, &userIDToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return model.ReconstructUser(user.ID, model.Name(name), model.Email(email), userIDToken), nil
}

// UpdateUser はユーザー情報を更新する
// 引数:
//   - ctx: コンテキスト
//   - user: 更新するユーザー情報
// 戻り値:
//   - *model.User: 更新されたユーザー情報
//   - error: エラー情報
// 実装: データベースのユーザー情報を更新する
// 注意事項: ユーザーIDは更新できない
func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}
