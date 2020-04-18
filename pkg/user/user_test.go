package user

import (
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/victornm/es-backend/pkg/store"
	"testing"
)

func TestGetProfile(t *testing.T) {
	users := []*store.UserRow{
		{
			ID:    1,
			Email: "admin@admin.com",
		},
	}

	db := func() *mockUserDAO {
		dao := newMockUserDao()
		dao.seed(users)

		return dao
	}

	tests := []struct {
		id            int
		wantedProfile *ProfileDTO
		wantedErr     error
	}{
		{1, &ProfileDTO{ID: users[0].ID, Email: users[0].Email}, nil},
		{10, nil, ErrNotFound},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("find with ID = %d", test.id), func(t *testing.T) {
			query := NewQueryService(db())
			gotProfile, gotErr := query.GetProfile(test.id)

			assert.Equal(t, test.wantedErr, gotErr)
			if gotErr == nil {
				assert.Equal(t, test.wantedProfile.Email, gotProfile.Email)
			}
		})
	}
}

// Mocks

type mockUserDAO struct {
	currentID int
	users     []*store.UserRow
}

func (dao *mockUserDAO) FindUserByID(id int) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
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

func (dao *mockUserDAO) CreateUser(u *store.UserRow) (int, error) {
	for _, row := range dao.users {
		if u.Email == row.Email || u.Username == row.Username {
			return 0, errors.New("user existed")
		}
	}

	dao.currentID++
	u.ID = dao.currentID
	dao.users = append(dao.users, u)

	return u.ID, nil
}
