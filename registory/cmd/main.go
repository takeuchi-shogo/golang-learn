package main

import (
	"database/sql"
	"log"

	"github.com/takeuchi-shogo/golang-learn/registory/repository"
	"github.com/takeuchi-shogo/golang-learn/registory/service"
)

func main() {
	repository := repository.NewUserRepository(&sql.DB{})
	service := service.NewUserService(repository)
	if err := service.Run(); err != nil {
		log.Fatalf("failed to run service: %v", err)
	}
}
