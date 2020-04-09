package store

import (
	"time"
)

type UserRow struct {
	ID             int
	Email          string
	Username       string
	HashedPassword string
	FullName       string
	Phone          string
	YearOfBirth    int
	Country        string
	Gender         string
	Language       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsActive       bool
	IsSuperAdmin   bool
	ActivationKey  string
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
