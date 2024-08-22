package auth

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth/command"
	iamservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
	adapter "github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/rbac"
)

type AuthUseCase interface {
	Login(ctx context.Context, cmd *command.LoginCommand) (*command.LoginCommandResult, error)
	Logout(ctx context.Context, cmd *command.LogoutCommand) (*command.LogoutCommandResult, error)
	VerifyToken(ctx context.Context, cmd *command.VerifyTokenCommand) (*command.VerifyTokenCommandResult, error)
	VerifyPermission(ctx context.Context, cmd *command.VerifyPermissionCommand) (*command.VerifyPermissionCommandResult, error)
}

type AuthApplicationService struct {
	Commands *command.AuthCommands
}

var _ AuthUseCase = (*AuthApplicationService)(nil)

func NewAuthApplicationService(
	tokenDomainService *iamservice.TokenDomainService,
	userDomainService *iamservice.UserDomainService,
	passwordDomainService *iamservice.PasswordDomainService,
	rbacDomainService *iamservice.RBACDomainService,
	rbacAdapter adapter.RBACAdapter,
) *AuthApplicationService {
	loginHandler := command.NewLoginHandler(tokenDomainService, userDomainService, passwordDomainService)
	logoutHandler := command.NewLogoutHandler()
	verifyTokenHandler := command.NewVerifyTokenHandler(tokenDomainService)
	verifyPermissionHandler := command.NewVerifyPermissionCommandHandler(rbacDomainService)

	authCommands := command.NewAuthCommands(loginHandler, logoutHandler, verifyTokenHandler, verifyPermissionHandler)

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

func (svc *AuthApplicationService) VerifyPermission(ctx context.Context, cmd *command.VerifyPermissionCommand) (*command.VerifyPermissionCommandResult, error) {
	return svc.Commands.VerifyPermission.Handle(ctx, cmd)
}
