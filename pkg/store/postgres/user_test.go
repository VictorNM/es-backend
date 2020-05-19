// +build database_test

package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/victornm/es-backend/pkg/store"
	"testing"
)

var db *sqlx.DB

func init() {
	db, _ = sqlx.Open("postgres", "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable&search_path=public")
}

func TestConnect(t *testing.T) {
	u := NewUserGateway(db)
	id, err := u.CreateUser(&store.UserRow{
		Email:          "victor.nguyenmau@gmail.com",
		Username:       "victornm",
		HashedPassword: "1923ashdjasd918239213",
		FullName:       "Nguyen Mau Vinh",
	})

	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
}
