package po

import (
	"time"
)

type User struct {
	ID        int `gorm:"primaryKey"`
	Username  string
	Password  string
	RoleID    int
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
