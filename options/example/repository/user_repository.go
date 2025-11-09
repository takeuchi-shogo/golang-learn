package repository

import (
	"database/sql"
	"errors"

	"github.com/takeuchi-shogo/golang-learn/options/example/model"
	"github.com/takeuchi-shogo/golang-learn/options/option"
)

type UserRepository interface {
	GetUserById(id int) (option.Option[model.User], error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserById(id int) (option.Option[model.User], error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return option.None[model.User](), nil
		}
		if errors.Is(err, sql.ErrConnDone) {
			return option.None[model.User](), errors.New("connection is already closed")
		}
		return option.None[model.User](), errors.New("failed to get user")
	}
	return option.Some(user), nil
}
