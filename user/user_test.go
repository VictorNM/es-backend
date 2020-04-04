package user

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/victornm/es-backend/store"
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
