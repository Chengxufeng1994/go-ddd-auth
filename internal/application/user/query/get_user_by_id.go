package query

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type GetUserByIDQueryHandler interface {
	Handle(ctx context.Context, query *GetUserByIDQuery) (*GetUserByIDQueryResult, error)
}

type getUserByIDHandler struct {
	userDomainService *service.UserDomainService
}

var _ GetUserByIDQueryHandler = (*getUserByIDHandler)(nil)

func NewGetUserByIDHandler(userDomainService *service.UserDomainService) *getUserByIDHandler {
	return &getUserByIDHandler{
		userDomainService: userDomainService,
	}
}

func (h *getUserByIDHandler) Handle(ctx context.Context, query *GetUserByIDQuery) (*GetUserByIDQueryResult, error) {
	ent, err := h.userDomainService.GetUserByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	var result GetUserByIDQueryResult
	result.User = ent
	return &result, nil
}
