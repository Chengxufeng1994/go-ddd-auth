package command

import (
	"context"

	identityaccessmgmtservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
)

type VerifyTokenCommandHandler interface {
	Handle(ctx context.Context, cmd *VerifyTokenCommand) (*VerifyTokenCommandResult, error)
}

type verifyTokenHandler struct {
	tokenDomainService *identityaccessmgmtservice.TokenDomainService
}

var _ VerifyTokenCommandHandler = (*verifyTokenHandler)(nil)

func NewVerifyTokenHandler(tokenDomainService *identityaccessmgmtservice.TokenDomainService) VerifyTokenCommandHandler {
	return &verifyTokenHandler{
		tokenDomainService: tokenDomainService,
	}
}

func (h *verifyTokenHandler) Handle(ctx context.Context, cmd *VerifyTokenCommand) (*VerifyTokenCommandResult, error) {
	_, authDetails, err := h.tokenDomainService.VerifyToken(ctx, cmd.TokenValue)
	if err != nil {
		return nil, err
	}
	return &VerifyTokenCommandResult{
		UserID:   authDetails.User.UserID,
		RoleID:   authDetails.User.RoleID,
		Username: authDetails.User.Username,
	}, nil
}
