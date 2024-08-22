package query

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
)

type Queries interface {
	GetUserByID(ctx context.Context, query *GetUserByIDQuery) error
	SearchUser(ctx context.Context, query *SearchUsersQuery) (*SearchUsersQueryResult, error)
}

type GetUserByIDQuery struct {
	ID int
}

type GetUserByIDQueryResult struct {
	User *aggregate.User
}

func NewGetUserByIDQuery(id int) *GetUserByIDQuery {
	return &GetUserByIDQuery{
		ID: id,
	}
}

type SearchUsersQuery struct {
	page    int
	perPage int
	orderBy string
	sortBy  string
	term    string
}

func NewSearchUsersQuery(page, perPage int, orderBy, sortBy, term string) *SearchUsersQuery {
	return &SearchUsersQuery{
		page:    page,
		perPage: perPage,
		orderBy: orderBy,
		sortBy:  sortBy,
		term:    term,
	}
}

type SearchUsersQueryResult struct {
	Users      []*aggregate.User
	Pagination *util.Pagination
}

type UserQueries struct {
	GetUserByID GetUserByIDQueryHandler
	SearchUsers SearchUsersQueryHandler
}

func NewUserQueries(
	getUserByID GetUserByIDQueryHandler,
	searchUsers SearchUsersQueryHandler,
) *UserQueries {
	return &UserQueries{
		GetUserByID: getUserByID,
		SearchUsers: searchUsers,
	}
}
