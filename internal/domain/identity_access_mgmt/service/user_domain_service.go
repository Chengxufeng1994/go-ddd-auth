package service

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
)

type UserDomainService struct {
	repository facade.UserRepository
}

func NewUserDomainService(
	repository facade.UserRepository,
) *UserDomainService {
	return &UserDomainService{
		repository: repository,
	}
}

func (s *UserDomainService) DeleteUserByID(ctx context.Context, id int) error {
	return s.repository.Remove(ctx, id)
}

func (s *UserDomainService) GetUserByID(ctx context.Context, id int) (*aggregate.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

func (s *UserDomainService) GetUserByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	return s.repository.GetUserByUsername(ctx, username)
}

func (s *UserDomainService) SearchUsers(ctx context.Context, opts *aggregate.SearchUserOpts) ([]*aggregate.User, *util.Pagination, error) {
	return s.repository.SearchUsers(ctx, opts)
}
