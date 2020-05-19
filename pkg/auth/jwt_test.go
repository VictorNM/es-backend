package auth_test

import (
	"testing"

	"github.com/victornm/es-backend/pkg/auth"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	t.Run("receive valid token", func(t *testing.T) {
		s := auth.NewJWTService("#12345", 24)

		tokenString, err := auth.GenerateToken(s, &auth.User{ID: 1})
		if err != nil {
			t.FailNow()
		}

		u, err := s.ParseToken(tokenString)
		assert.NoError(t, err)

		assert.Equal(t, 1, u.UserID)
	})
}
