package command

import (
	"context"
	"fmt"

	roledomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/service"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/entity/valueobject"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type UpdateUserCommandHandler interface {
	Handle(ctx context.Context, cmd *UpdateUserCommand) error
}

type updateUserHandler struct {
	roleDomainService *roledomainservice.RoleDomainService
	userDomainService *userdomainservice.UserDomainService
}

var _ UpdateUserCommandHandler = (*updateUserHandler)(nil)

func NewUpdateUserHandler(
	roleDomainService *roledomainservice.RoleDomainService,
	userDomainService *userdomainservice.UserDomainService,
) *updateUserHandler {
	return &updateUserHandler{
		roleDomainService: roleDomainService,
		userDomainService: userDomainService,
	}
}

func (h *updateUserHandler) Handle(ctx context.Context, cmd *UpdateUserCommand) error {
	var roleVo *valueobject.Role
	if cmd.RoleID != 0 {
		isExists, err := h.roleDomainService.Exists(ctx, cmd.RoleID)
		if err != nil || isExists == false {
			return fmt.Errorf("role %d not exists", cmd.RoleID)
		}
		roleVo = valueobject.NewRole(cmd.RoleID)
	}

	user := entity.NewUser(cmd.Username, cmd.Password, roleVo)
	user.ID = cmd.UserID
	return h.userDomainService.UpdateUser(ctx, user)
}
