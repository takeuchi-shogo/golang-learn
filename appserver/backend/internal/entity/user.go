package entity

import "github.com/google/uuid"

type user struct {
	id          uuid.UUID
	userIDToken string
	name        string
	email       string
}

func toEntity(id uuid.UUID, name string, email string, userIDToken string) *user {
	return &user{
		id:          id,
		name:        name,
		email:       email,
		userIDToken: userIDToken,
	}
}
