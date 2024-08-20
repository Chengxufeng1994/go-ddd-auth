package entity

import (
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity/valueobject"
)

// Aggregate root
type User struct {
	ID        int
	Username  string
	Password  string
	Role      *valueobject.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username, password string, role *valueobject.Role) *User {
	return &User{
		Username: username,
		Password: password,
		Role:     role,
	}
}

func (u *User) Create() {
	u.CreatedAt = time.Now().UTC()
}

func (u *User) Update() {
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) UpdateUsername(username string) {
	// TODO: validate
	u.Username = username
}

func (u *User) UpdatePassword(password string) {
	// TODO: validate
	u.Password = password
}

func (u *User) UpdateRole(role *valueobject.Role) {
	u.Role = role
}
