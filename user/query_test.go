package user_test

import (
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/victornm/es-backend/store"
	"github.com/victornm/es-backend/user"
	"testing"
)

type mockUserDAO struct {
	users []*store.UserRow
}

func (dao *mockUserDAO) FindUserByID(id int) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

var users = []*store.UserRow{
	{
		ID:    1,
		Email: "admin@admin.com",
	},
}

func fixedUserStore() *mockUserDAO {
	return &mockUserDAO{
		users: users,
	}
}

func TestGetProfile(t *testing.T) {
	tests := []struct {
		id            int
		wantedProfile *user.ProfileDTO
		wantedErr     error
	}{
		{1, &user.ProfileDTO{ID: users[0].ID, Email: users[0].Email}, nil},
		{10, nil, user.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("find with ID = %d", test.id), func(t *testing.T) {
			query := user.NewQueryService(fixedUserStore())
			gotProfile, gotErr := query.GetProfile(test.id)

			assert.Equal(t, test.wantedErr, gotErr)
			if gotErr == nil {
				assert.Equal(t, test.wantedProfile.Email, gotProfile.Email)
			}
		})
	}
}
