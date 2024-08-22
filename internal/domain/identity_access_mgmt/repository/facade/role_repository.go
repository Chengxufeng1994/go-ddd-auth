package facade

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
)

type RoleRepository interface {
	GetByID(ctx context.Context, roleID int) (*aggregate.Role, error)
}
