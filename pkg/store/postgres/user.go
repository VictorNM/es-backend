package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/victornm/es-backend/pkg/store"
	"time"
)

type DB interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type UserGateway struct {
	db DB
}

func NewUserGateway(db DB) *UserGateway {
	return &UserGateway{db: db}
}

func (gw *UserGateway) CreateUser(u *store.UserRow) (int, error) {
	u.CreatedAt = time.Now()

	stmt, err := gw.db.PrepareNamed(
		`INSERT INTO users (email, username, hashed_password, full_name) VALUES(:email, :username, :hashed_password, :full_name) RETURNING id;`,
	)

	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.Get(&id, u)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
