package memory

import (
	"errors"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"sync"
	"time"
)

var UserStore *userStore

func init() {
	UserStore = NewUserStore()
}

var fixedUsers = []*store.UserRow{
	{
		Email:          "admin@es.com",
		Username:       "admin",
		HashedPassword: genPassword("admin"),
		IsSuperAdmin:   true,
		IsActive:       true,
	},
	{
		Email:          "admin2@es.com",
		HashedPassword: genPassword("admin"),
		IsActive:       true,
	},
	{
		Email:          "admin3@es.com",
		HashedPassword: genPassword("admin"),
		IsActive:       true,
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
	mu        sync.RWMutex
	currentID int
	users     []*store.UserRow
}

func (dao *userStore) FindUserByEmail(email string) (*store.UserRow, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, u := range dao.users {
		if strings.ToLower(u.Email) == strings.ToLower(email) {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *userStore) FindUserByUsername(username string) (*store.UserRow, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, u := range dao.users {
		if strings.ToLower(u.Username) == strings.ToLower(username) {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *userStore) FindUserByID(id int) (*store.UserRow, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, u := range dao.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (dao *userStore) CreateUser(u *store.UserRow) (int, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	for _, row := range dao.users {
		if u.Email == row.Email {
			return 0, errors.New("email existed")
		}
	}

	dao.currentID++
	u.ID = dao.currentID
	u.CreatedAt = time.Now()
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
