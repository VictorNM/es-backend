package user

import "github.com/victornm/es-backend/store"

type Finder interface {
	FindUserByID(id int) (*store.UserRow, error)
	FindUserByEmail(email string) (*store.UserRow, error)		// Note: case-insensitive query
	FindUserByUsername(username string) (*store.UserRow, error) // Note: case-insensitive query
}

type Creator interface {
	CreateUser(*store.UserRow) (int, error)
}

type FindCreator interface {
	Finder
	Creator
}
