package user

import "github.com/victornm/es-backend/store"

type Finder interface {
	FindUserByID(id int) (*store.UserRow, error)
	FindUserByEmail(email string) (*store.UserRow, error)
}

type Creater interface {
	CreateUser(*store.UserRow) (int, error)
}

type FindCreater interface {
	Finder
	Creater
}
