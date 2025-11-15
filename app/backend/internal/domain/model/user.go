package model

import (
	"github.com/google/uuid"
)

type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypeCognito
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeNone:
		return "none"
	case AuthTypeCognito:
		return "cognito"
	}
	return "unknown"
}

type (
	User struct {
		ID          uuid.UUID `json:"id"`
		UserIDToken string    `json:"user_id_token"`
		Name        Name      `json:"name"`
		Email       Email     `json:"email"`
	}

	Name  string
	Email string
)

func newUser(id uuid.UUID, name Name, email Email, userIDToken string) *User {
	return &User{
		ID:          id,
		Name:        name,
		Email:       email,
		UserIDToken: userIDToken,
	}
}

func NewUser(name Name, email Email, userIDToken string) *User {
	id := uuid.Must(uuid.NewV7())
	return newUser(id, name, email, userIDToken)
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) GetName() Name {
	return u.Name
}

func (u *User) GetEmail() Email {
	return u.Email
}
