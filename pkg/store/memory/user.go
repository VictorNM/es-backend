package memory

import (
	"errors"
	"github.com/victornm/es-backend/pkg/store"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"sync"
	"time"
)

var GlobalUserStore *UserGateway

func init() {
	GlobalUserStore = NewUserGateway()
	GlobalUserStore.Seed(fixedUsers)
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

type UserGateway struct {
	mu        sync.RWMutex
	currentID int
	users     []*store.UserRow
}

func (gw *UserGateway) FindUserByEmail(email string) (*store.UserRow, error) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	for _, u := range gw.users {
		if strings.ToLower(u.Email) == strings.ToLower(email) {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (gw *UserGateway) FindUserByUsername(username string) (*store.UserRow, error) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	for _, u := range gw.users {
		if strings.ToLower(u.Username) == strings.ToLower(username) {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (gw *UserGateway) FindUserByID(id int) (*store.UserRow, error) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	for _, u := range gw.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (gw *UserGateway) CreateUser(u *store.UserRow) (int, error) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	for _, row := range gw.users {
		if u.Email == row.Email {
			return 0, errors.New("email existed")
		}
	}

	gw.currentID++
	u.ID = gw.currentID
	u.CreatedAt = time.Now()
	gw.users = append(gw.users, u)

	return u.ID, nil
}

func (gw *UserGateway) Seed(users []*store.UserRow) {
	for _, u := range users {
		_, err := gw.CreateUser(u)
		if err != nil {
			log.Panicf("seeding memory user failed: %v", err)
		}
	}

	return
}

func (gw *UserGateway) Clear() {
	gw.currentID = 0
	gw.users = nil
}

func NewUserGateway() *UserGateway {
	return &UserGateway{currentID: 0, mu: sync.RWMutex{}}
}
