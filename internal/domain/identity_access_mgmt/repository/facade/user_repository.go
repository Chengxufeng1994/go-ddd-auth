package facade

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
)

type UserRepository interface {
	Save(ctx context.Context, aggregate *aggregate.User) error
	Update(ctx context.Context, aggregate *aggregate.User) error
	Remove(ctx context.Context, userID int) error

	GetUserByID(ctx context.Context, userID int) (*aggregate.User, error)
	GetUserByUsername(ctx context.Context, username string) (*aggregate.User, error)
	SearchUsers(ctx context.Context, query *aggregate.SearchUserOpts) ([]*aggregate.User, *util.Pagination, error)
}
