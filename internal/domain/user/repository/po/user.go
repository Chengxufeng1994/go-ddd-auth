package po

import (
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/po"
)

type User struct {
	ID        int `gorm:"primaryKey"`
	Username  string
	Password  string
	RoleID    int
	Role      po.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
