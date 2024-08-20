package query

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
)

type Queries interface {
	GetUserByID(ctx context.Context, query *GetUserByIDQuery) error
}

type GetUserByIDQuery struct {
	ID int
}

type GetUserByIDQueryResult struct {
	User *entity.User
}

func NewGetUserByIDQuery(id int) *GetUserByIDQuery {
	return &GetUserByIDQuery{
		ID: id,
	}
}

type UserQueries struct {
	GetUserByID GetUserByIDQueryHandler
}

func NewUserQueries(
	getUserByID GetUserByIDQueryHandler,
) *UserQueries {
	return &UserQueries{
		GetUserByID: getUserByID,
	}
}
