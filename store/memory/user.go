package memory

import (
	"errors"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var fixedUsers = []*store.UserRow{
	{
		ID:             1,
		Email:          "admin1@es.com",
		HashedPassword: genPassword("admin"),
	},
	{
		ID:             2,
		Email:          "admin2@es.com",
		HashedPassword: genPassword("admin"),
	},
	{
		ID:             3,
		Email:          "admin3@es.com",
		HashedPassword: genPassword("admin"),
	},
}

func genPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hashed)
}

type userStore struct {
	users []*store.UserRow
}

func (dao *userStore) FindUserByEmail(email string) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *userStore) FindUserByID(id int) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func NewUserStore() *userStore {
	return &userStore{users: fixedUsers}
}
