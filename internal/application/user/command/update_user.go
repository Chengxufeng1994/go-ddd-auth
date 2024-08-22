package command

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	iamservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/transaction"
)

type UpdateUserCommandHandler interface {
	Handle(ctx context.Context, cmd *UpdateUserCommand) error
}

type updateUserHandler struct {
	userRepository        facade.UserRepository
	userDomainService     *iamservice.UserDomainService
	passwordDomainService *iamservice.PasswordDomainService
	rbacDomainService     *iamservice.RBACDomainService
	trxMgr                transaction.TransactionManager
}

var _ UpdateUserCommandHandler = (*updateUserHandler)(nil)

func NewUpdateUserHandler(
	userRepository facade.UserRepository,
	userDomainService *iamservice.UserDomainService,
	passwordDomainService *iamservice.PasswordDomainService,
	rbacDomainService *iamservice.RBACDomainService,
	trxMgr transaction.TransactionManager,
) *updateUserHandler {
	return &updateUserHandler{
		userRepository:        userRepository,
		userDomainService:     userDomainService,
		passwordDomainService: passwordDomainService,
		rbacDomainService:     rbacDomainService,
		trxMgr:                trxMgr,
	}
}

func (h *updateUserHandler) Handle(ctx context.Context, cmd *UpdateUserCommand) error {
	return h.trxMgr.Do(ctx, func(txctx context.Context) error {
		if cmd.opt.RoleID != nil && *cmd.opt.RoleID != 0 {
			isExists, err := h.rbacDomainService.ExistingRole(txctx, *cmd.opt.RoleID)
			if err != nil || isExists == false {
				return fmt.Errorf("role %d not exists", *cmd.opt.RoleID)
			}
		}

		user, err := h.userRepository.GetUserByID(txctx, cmd.userID)

		var hashedPassword string
		if cmd.opt.Password != nil && *cmd.opt.Password != "" {
			hashedPassword, err = h.passwordDomainService.ChangedPassword(user.Password, *cmd.opt.Password)
			if err != nil {
				return fmt.Errorf("password changed error: %w", err)
			}
		}

		user.Update(aggregate.UpdateUserOpt{
			Username: cmd.opt.Username,
			Password: &hashedPassword,
			RoleID:   cmd.opt.RoleID,
		})

		err = h.userRepository.Update(txctx, user)
		if err != nil {
			return fmt.Errorf("update error: %w", err)
		}

		return nil
	})
}
