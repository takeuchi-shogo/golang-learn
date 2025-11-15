package userregistory

import (
	"database/sql"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/repository"
	usercommand "github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/command"
	userquery "github.com/takeuchi-shogo/golang-learn/app/backend/internal/service/query"
)

type userRegistory struct {
	db *sql.DB
}

type UserRegistory interface {
	UserQuery() userquery.UserQuery
	UserCommand() usercommand.UserCommand
}

func NewUserRegistory(db *sql.DB) UserRegistory {
	return &userRegistory{db: db}
}

func (f *userRegistory) UserQuery() userquery.UserQuery {
	return repository.NewUserRepository(f.db)
}

func (f *userRegistory) UserCommand() usercommand.UserCommand {
	return repository.NewUserRepository(f.db)
}
