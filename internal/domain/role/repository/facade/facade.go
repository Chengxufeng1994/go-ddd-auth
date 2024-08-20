package facade

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/entity"
)

type RoleRepository interface {
	GetByID(ctx context.Context, roleID int) (*entity.Role, error)
}
