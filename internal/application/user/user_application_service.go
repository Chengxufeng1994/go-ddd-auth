package user

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/command"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/query"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/transaction"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/repository/facade"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, cmd *command.CreateUserCommand) error
	UpdateUser(ctx context.Context, cmd *command.UpdateUserCommand) error
	DeleteUser(ctx context.Context, cmd *command.DeleteUserCommand) error
	GetByID(ctx context.Context, query *query.GetUserByIDQuery) (*query.GetUserByIDQueryResult, error)
	SearchUsers(ctx context.Context, query *query.SearchUsersQuery) (*query.SearchUsersQueryResult, error)
}

type UserApplicationService struct {
	Commands *command.UserCommands
	Queries  *query.UserQueries
}

var _ UserUseCase = (*UserApplicationService)(nil)

func NewUserApplicationService(
	userRepository facade.UserRepository,
	userDomainService *userdomainservice.UserDomainService,
	passwordDomainService *userdomainservice.PasswordDomainService,
	rbacDomainServcie *userdomainservice.RBACDomainService,
	trxMgr transaction.TransactionManager,
) *UserApplicationService {

	createUserHandler := command.NewCreateUserHandler(userRepository, userDomainService, passwordDomainService, rbacDomainServcie, trxMgr)
	updateUserHandler := command.NewUpdateUserHandler(userRepository, userDomainService, passwordDomainService, rbacDomainServcie, trxMgr)
	deleteUserHandler := command.NewDeleteUserHandler(userDomainService)

	getUserByIDHandler := query.NewGetUserByIDHandler(userDomainService)
	searchUsersHandler := query.NewSearchUsersHandler(userDomainService)

	userCommands := command.NewUserCommands(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
	)
	userQueries := query.NewUserQueries(
		getUserByIDHandler,
		searchUsersHandler,
	)

	return &UserApplicationService{
		Commands: userCommands,
		Queries:  userQueries,
	}
}

func (svc *UserApplicationService) CreateUser(ctx context.Context, cmd *command.CreateUserCommand) error {
	return svc.Commands.CreateUser.Handle(ctx, cmd)
}

func (svc *UserApplicationService) UpdateUser(ctx context.Context, cmd *command.UpdateUserCommand) error {
	return svc.Commands.UpdateUser.Handle(ctx, cmd)
}

func (svc *UserApplicationService) DeleteUser(ctx context.Context, cmd *command.DeleteUserCommand) error {
	return svc.Commands.DeleteUser.Handle(ctx, cmd)
}

func (svc *UserApplicationService) GetByID(ctx context.Context, q *query.GetUserByIDQuery) (*query.GetUserByIDQueryResult, error) {
	return svc.Queries.GetUserByID.Handle(ctx, q)
}

func (svc *UserApplicationService) SearchUsers(ctx context.Context, query *query.SearchUsersQuery) (*query.SearchUsersQueryResult, error) {
	return svc.Queries.SearchUsers.Handle(ctx, query)
}
