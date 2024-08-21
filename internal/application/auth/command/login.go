package command

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/entity"
	tokendomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/service"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type LoginCommandHandler interface {
	Handle(ctx context.Context, cmd *LoginCommand) (*LoginCommandResult, error)
}

var _ LoginCommandHandler = (*LoginHandler)(nil)

type LoginHandler struct {
	tokenDomainService *tokendomainservice.TokenDomainService
	userDomainService  *userdomainservice.UserDomainService
}

func NewLoginHandler(tokenDomainService *tokendomainservice.TokenDomainService, userDomainService *userdomainservice.UserDomainService) *LoginHandler {
	return &LoginHandler{
		tokenDomainService: tokenDomainService,
		userDomainService:  userDomainService,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd *LoginCommand) (*LoginCommandResult, error) {
	user, err := h.userDomainService.VerifyUser(ctx, cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}

	userDetails := &entity.UserDetails{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   user.Role.RoleID,
	}

	accessToken, refreshToken, err := h.tokenDomainService.CreateAccessToken(ctx, userDetails)
	if err != nil {
		return nil, err
	}

	return &LoginCommandResult{
			AccessToken:  accessToken.TokenValue,
			RefreshToken: refreshToken.TokenValue},
		nil
}
