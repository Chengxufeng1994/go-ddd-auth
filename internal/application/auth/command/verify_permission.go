package command

import (
	"context"

	iamservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
)

type VerifyPermissionCommandHandler interface {
	Handle(ctx context.Context, cmd *VerifyPermissionCommand) (*VerifyPermissionCommandResult, error)
}

type verifyPermissionHandler struct {
	roleDomainService *iamservice.RBACDomainService
}

var _ VerifyPermissionCommandHandler = (*verifyPermissionHandler)(nil)

func NewVerifyPermissionCommandHandler(
	roleDomainService *iamservice.RBACDomainService,
) *verifyPermissionHandler {
	return &verifyPermissionHandler{
		roleDomainService: roleDomainService,
	}
}

func (h *verifyPermissionHandler) Handle(ctx context.Context, cmd *VerifyPermissionCommand) (*VerifyPermissionCommandResult, error) {
	ok, err := h.roleDomainService.VerifyPermission(ctx, cmd.RoleID, cmd.Resource, cmd.Action)
	if err != nil || !ok {
		return nil, err
	}
	return &VerifyPermissionCommandResult{}, nil
}
