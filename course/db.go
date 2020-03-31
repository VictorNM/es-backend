package course

import "time"

type Row struct {
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
