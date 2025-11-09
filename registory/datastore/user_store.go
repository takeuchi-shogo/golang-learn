package datastore

import (
	"database/sql"
	"errors"

	"github.com/takeuchi-shogo/golang-learn/registory/model"
)

type userStore struct {
	db *sql.DB
}

type UserStore interface {
	GetUserById(id int) (*model.User, error)
}

func NewUserStore(db *sql.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) GetUserById(id int) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	user := model.NewUser(id, "test", "test@example.com")
	return user, nil
}
