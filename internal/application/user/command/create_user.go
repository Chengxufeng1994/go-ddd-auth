package command

import (
	"context"
	"fmt"

	roledomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/service"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity/valueobject"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type CreateUserCommandHandler interface {
	Handle(ctx context.Context, cmd *CreateUserCommand) error
}

type createUserHandler struct {
	roleDomainService *roledomainservice.RoleDomainService
	userDomainService *userdomainservice.UserDomainService
}

var _ CreateUserCommandHandler = (*createUserHandler)(nil)

func NewCreateUserHandler(
	roleDomainService *roledomainservice.RoleDomainService,
	userDomainService *userdomainservice.UserDomainService,
) *createUserHandler {
	return &createUserHandler{
		roleDomainService: roleDomainService,
		userDomainService: userDomainService,
	}
}

func (h *createUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
	// TODO: check if role exists (roleDomainService)
	isExists, err := h.roleDomainService.Exists(ctx, cmd.RoleID)
	if err != nil || isExists == false {
		return fmt.Errorf("role %d not exists", cmd.RoleID)
	}

	role := valueobject.NewRole(cmd.RoleID)
	user := entity.NewUser(cmd.Username, cmd.Password, role)
	return h.userDomainService.CreateUser(ctx, user)
}
