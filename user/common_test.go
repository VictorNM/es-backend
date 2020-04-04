package user

import (
	"errors"
	"github.com/victornm/es-backend/store"
	"log"
)

type mockUserDAO struct {
	currentID int
	users     []*store.UserRow
}

func (dao *mockUserDAO) FindUserByEmail(email string) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *mockUserDAO) FindUserByID(id int) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *mockUserDAO) CreateUser(u *store.UserRow) (int, error) {
	for _, row := range dao.users {
		if u.Email == row.Email {
			return 0, ErrEmailExisted
		}
	}

	dao.currentID++
	u.ID = dao.currentID
	dao.users = append(dao.users, u)

	return u.ID, nil
}

func newMockUserDao() *mockUserDAO {
	return &mockUserDAO{currentID: 0}
}

func (dao *mockUserDAO) seed(users []*store.UserRow) {
	for _, u := range users {
		_, err := dao.CreateUser(u)
		if err != nil {
			panic(err)
		}
	}
}

func mustHashPassword(password string) string {
	hashed, err := hashPassword(password)
	if err != nil {
		log.Panic(err)
	}

	return hashed
}
