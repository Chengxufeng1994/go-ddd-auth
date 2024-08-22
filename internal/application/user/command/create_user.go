package command

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/transaction"
)

type CreateUserCommandHandler interface {
	Handle(ctx context.Context, cmd *CreateUserCommand) error
}

type createUserHandler struct {
	userRepository        facade.UserRepository
	rbacDomainService     *userdomainservice.RBACDomainService
	userDomainService     *userdomainservice.UserDomainService
	passwordDomainService *userdomainservice.PasswordDomainService
	trxMgr                transaction.TransactionManager
}

var _ CreateUserCommandHandler = (*createUserHandler)(nil)

func NewCreateUserHandler(
	userRepository facade.UserRepository,
	userDomainService *userdomainservice.UserDomainService,
	passwordDomainService *userdomainservice.PasswordDomainService,
	rbacDomainService *userdomainservice.RBACDomainService,
	trxMgr transaction.TransactionManager,
) *createUserHandler {
	return &createUserHandler{
		userRepository:        userRepository,
		userDomainService:     userDomainService,
		passwordDomainService: passwordDomainService,
		rbacDomainService:     rbacDomainService,
		trxMgr:                trxMgr,
	}
}

// TODO: add transaction in this use case
func (h *createUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
	return h.trxMgr.Do(ctx, func(ctx context.Context) error {
		ok, err := h.rbacDomainService.ExistingRole(ctx, cmd.RoleID)
		if err != nil || !ok {
			return fmt.Errorf("role %d not exists", cmd.RoleID)
		}

		hashedPassword, err := h.passwordDomainService.EncryptPassword(cmd.Password)
		if err != nil {
			return fmt.Errorf("password encrypt error: %w", err)
		}

		existingUser, err := h.userRepository.GetUserByUsername(ctx, cmd.Username)
		if err != nil {
			return fmt.Errorf("user %s already exists", cmd.Username)
		}
		if existingUser != nil {
			return fmt.Errorf("user %s already exists", cmd.Username)
		}

		roleID, _ := valueobject.NewRoleID(cmd.RoleID)
		user, err := aggregate.NewUser(cmd.Username, hashedPassword, &roleID)
		if err != nil {
			return err
		}

		user.Create()

		err = h.userRepository.Save(ctx, user)
		if err != nil {
			return err
		}

		return nil

	}, transaction.PropagationRequired)
}
