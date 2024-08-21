package auth

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth/command"
	tokendomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/service"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type AuthUseCase interface {
	Login(ctx context.Context, cmd *command.LoginCommand) (*command.LoginCommandResult, error)
	Logout(ctx context.Context, cmd *command.LogoutCommand) (*command.LogoutCommandResult, error)
	VerifyToken(ctx context.Context, cmd *command.VerifyTokenCommand) (*command.VerifyTokenCommandResult, error)
}

type AuthApplicationService struct {
	Commands *command.AuthCommands
}

var _ AuthUseCase = (*AuthApplicationService)(nil)

func NewAuthApplicationService(tokenDomainService *tokendomainservice.TokenDomainService, userDomainService *userdomainservice.UserDomainService) *AuthApplicationService {
	loginHandler := command.NewLoginHandler(tokenDomainService, userDomainService)
	logoutHandler := command.NewLogoutHandler()
	verifyTokenHandler := command.NewVerifyTokenHandler(tokenDomainService)

	authCommands := command.NewAuthCommands(loginHandler, logoutHandler, verifyTokenHandler)

	return &AuthApplicationService{
		Commands: authCommands,
	}
}

func (svc *AuthApplicationService) Login(ctx context.Context, cmd *command.LoginCommand) (*command.LoginCommandResult, error) {
	return svc.Commands.Login.Handle(ctx, cmd)
}

func (svc *AuthApplicationService) Logout(ctx context.Context, cmd *command.LogoutCommand) (*command.LogoutCommandResult, error) {
	panic("unimplemented")
}

func (svc *AuthApplicationService) VerifyToken(ctx context.Context, cmd *command.VerifyTokenCommand) (*command.VerifyTokenCommandResult, error) {
	return svc.Commands.VerifyToken.Handle(ctx, cmd)
}
