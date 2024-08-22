package command

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/entity"
	identityaccessmgmtservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
)

type LoginCommandHandler interface {
	Handle(ctx context.Context, cmd *LoginCommand) (*LoginCommandResult, error)
}

var _ LoginCommandHandler = (*LoginHandler)(nil)

type LoginHandler struct {
	tokenDomainService    *identityaccessmgmtservice.TokenDomainService
	userDomainService     *identityaccessmgmtservice.UserDomainService
	passwordDomainService *identityaccessmgmtservice.PasswordDomainService
}

func NewLoginHandler(
	tokenDomainService *identityaccessmgmtservice.TokenDomainService,
	userDomainService *identityaccessmgmtservice.UserDomainService,
	passwordDomainService *identityaccessmgmtservice.PasswordDomainService,
) *LoginHandler {
	return &LoginHandler{
		tokenDomainService: tokenDomainService,
		userDomainService:  userDomainService,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd *LoginCommand) (*LoginCommandResult, error) {
	user, err := h.userDomainService.GetUserByUsername(ctx, cmd.Username)
	if user == nil || err != nil {
		return nil, err
	}

	userDetails := &entity.UserDetails{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   int(*user.RoleID),
	}

	err = h.passwordDomainService.VerifyPassword(cmd.Password, user.Password)
	if err != nil {
		return nil, err
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
