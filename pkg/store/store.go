package store

import (
	"time"
)

// TODO: Set mapper func
type UserRow struct {
	ID             int       `db:"id"`
	Email          string    `db:"email"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hashed_password"`
	FullName       string    `db:"full_name"`
	Phone          string    `db:"phone"`
	YearOfBirth    int       `db:"year_of_birth"`
	Country        string    `db:"country"`
	Gender         string    `db:"gender"`
	Language       string    `db:"language"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	IsActive       bool      `db:"is_active"`
	IsSuperAdmin   bool      `db:"is_super_admin"`
	ActivationKey  string    `db:"activation_key"`
	OAuth2Provider string    `db:"oauth2_provider"`
}

type CourseRow struct {
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryID  int
}

type CategoryRow struct {
	ID   int
	Name string
}
