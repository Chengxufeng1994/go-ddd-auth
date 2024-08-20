package command

import (
	"context"

	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type DeleteUserCommandHandler interface {
	Handle(ctx context.Context, cmd *DeleteUserCommand) error
}

type DeleteUserHandler struct {
	userDomainService *userdomainservice.UserDomainService
}

var _ DeleteUserCommandHandler = (*DeleteUserHandler)(nil)

func NewDeleteUserHandler(userDomainService *userdomainservice.UserDomainService) *DeleteUserHandler {
	return &DeleteUserHandler{
		userDomainService: userDomainService,
	}
}

func (d *DeleteUserHandler) Handle(ctx context.Context, cmd *DeleteUserCommand) error {
	return d.userDomainService.DeleteUserByID(ctx, cmd.UserID)
}
