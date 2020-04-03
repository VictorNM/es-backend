package auth

import "github.com/victornm/es-backend/store"

type UserFinder interface {
	FindUserByEmail(email string) (*store.UserRow, error)
}