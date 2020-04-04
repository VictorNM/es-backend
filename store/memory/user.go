package memory

import (
	"errors"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var fixedUsers = []*store.UserRow{
	{
		Email:          "admin1@es.com",
		HashedPassword: genPassword("admin"),
	},
	{
		Email:          "admin2@es.com",
		HashedPassword: genPassword("admin"),
	},
	{
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
	currentID int
	users     []*store.UserRow
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

func (dao *userStore) CreateUser(u *store.UserRow) (int, error) {
	for _, row := range dao.users {
		if u.Email == row.Email {
			return 0, errors.New("email existed")
		}
	}

	dao.currentID++
	u.ID = dao.currentID
	dao.users = append(dao.users, u)

	return u.ID, nil
}

func NewUserStore() *userStore {
	s := &userStore{currentID: 0}
	for _, u := range fixedUsers {
		_, err := s.CreateUser(u)
		if err != nil {
			log.Fatalf("init memory user store failed: %v", err)
		}
	}

	return s
}
