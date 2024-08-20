package entity

import "time"

type Role struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
