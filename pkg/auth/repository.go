package auth

import (
	"github.com/victornm/es-backend/pkg/auth/internal"
)

type ReadUserRepository interface {
	FindUserByID(id int) (*internal.User, error)
	FindUserByEmail(email string) (*internal.User, error)
	FindUserByUsername(username string) (*internal.User, error)
}

type WriteUserRepository interface {
	CreateUser(u *internal.User) (int, error)
}

type UserRepository interface {
	ReadUserRepository
	WriteUserRepository
}
