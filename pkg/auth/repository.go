package auth

import "github.com/victornm/es-backend/pkg/store"

type ReadUserRepository interface {
	FindUserByID(id int) (*User, error)
	FindUserByEmail(email string) (*User, error)       // Note: case-insensitive query
	FindUserByUsername(username string) (*User, error) // Note: case-insensitive query
}

type WriteUserRepository interface {
	CreateUser(u *User) (int, error)
}

type UserRepository interface {
	ReadUserRepository
	WriteUserRepository
}

type MemoryRepository struct {
	dao Gateway
}

func (r *MemoryRepository) FindUserByID(id int) (*User, error) {
	row, err := r.dao.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return &User{UserRow: row}, nil
}

func (r *MemoryRepository) FindUserByEmail(email string) (*User, error) {
	row, err := r.dao.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{UserRow: row}, nil
}

func (r *MemoryRepository) FindUserByUsername(username string) (*User, error) {
	row, err := r.dao.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return &User{UserRow: row}, nil
}

func (r *MemoryRepository) CreateUser(u *User) (int, error) {
	id, err := r.dao.CreateUser(u.UserRow)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewRepository(dao Gateway) *MemoryRepository {
	return &MemoryRepository{
		dao: dao,
	}
}

type ReadGateway interface {
	FindUserByID(id int) (*store.UserRow, error)
	FindUserByEmail(email string) (*store.UserRow, error)       // Note: case-insensitive query
	FindUserByUsername(username string) (*store.UserRow, error) // Note: case-insensitive query
}

type WriteGateway interface {
	CreateUser(*store.UserRow) (int, error)
}

type Gateway interface {
	ReadGateway
	WriteGateway
}
