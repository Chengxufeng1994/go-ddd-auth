package facade

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
)

type UpdateFn func(*entity.User) (*entity.User, error)

type UserRepository interface {
	SaveUser(ctx context.Context, aggregate *entity.User) error
	UpdateUser(ctx context.Context, id int, updateFn UpdateFn) error
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	DeleteUserByID(ctx context.Context, userID int) error
}
