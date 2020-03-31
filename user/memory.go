package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var fixedUsers = []*DTO{
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

type memoryDAO struct {
	users []*DTO
}

func (dao *memoryDAO) FindUserByEmail(email string) (*DTO, error) {
	for _, u := range dao.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func NewMemoryDAO() *memoryDAO {
	return &memoryDAO{users: fixedUsers}
}
