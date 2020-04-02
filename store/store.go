package store

import "time"

type UserRow struct {
	ID             int
	Email          string
	HashedPassword string
	FirstName      string
	LastName       string
	Phone          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsActive       bool
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
