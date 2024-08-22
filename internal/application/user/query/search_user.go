package query

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	iamservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
)

type SearchUsersQueryHandler interface {
	Handle(ctx context.Context, query *SearchUsersQuery) (*SearchUsersQueryResult, error)
}

type searchUsersHandler struct {
	userDomainService *iamservice.UserDomainService
}

var _ SearchUsersQueryHandler = (*searchUsersHandler)(nil)

func NewSearchUsersHandler(
	userDomainService *iamservice.UserDomainService,
) *searchUsersHandler {
	return &searchUsersHandler{
		userDomainService: userDomainService,
	}
}

func (s *searchUsersHandler) Handle(ctx context.Context, query *SearchUsersQuery) (*SearchUsersQueryResult, error) {
	opts := &aggregate.SearchUserOpts{
		Page:    query.page,
		PerPage: query.perPage,
		OrderBy: query.orderBy,
		SortBy:  query.sortBy,
		Term:    query.term,
	}
	users, pagination, err := s.userDomainService.SearchUsers(ctx, opts)
	return &SearchUsersQueryResult{
		Users:      users,
		Pagination: pagination,
	}, err
}
