package aggregate

import (
	"time"
)

type RoleID int

const (
	AdminRoleID RoleID = 1
	UserRoleID  RoleID = 2
)

type Role struct {
	ID        RoleID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
