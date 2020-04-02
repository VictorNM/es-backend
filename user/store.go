package user

import "github.com/victornm/es-backend/store"

type FindUserByID interface {
	FindUserByID(id int) (*store.UserRow, error)
}
