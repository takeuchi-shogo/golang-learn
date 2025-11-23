package factory

import (
	"database/sql"

	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/factory/userregistory"
	"github.com/takeuchi-shogo/golang-learn/app/backend/internal/infra/database"
)

type Factory interface {
	GetUserRegistory() userregistory.UserRegistory
}

type factory struct {
	db *sql.DB
}

func NewFactory() Factory {
	db := database.Connect()
	return &factory{db: db}
}

func (f *factory) GetUserRegistory() userregistory.UserRegistory {
	return userregistory.NewUserRegistory(f.db)
}
