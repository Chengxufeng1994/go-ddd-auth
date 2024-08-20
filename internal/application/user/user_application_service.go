package application

import (
	"context"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/command"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/query"

	roledomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/service"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
)

type UserUseCase interface {
	GetByID(ctx context.Context, query *query.GetUserByIDQuery) (*query.GetUserByIDQueryResult, error)
	CreateUser(ctx context.Context, cmd *command.CreateUserCommand) error
	UpdateUser(ctx context.Context, cmd *command.UpdateUserCommand) error
	DeleteUser(ctx context.Context, cmd *command.DeleteUserCommand) error
}

type UserApplicationService struct {
	Commands *command.UserCommands
	Queries  *query.UserQueries
}

var _ UserUseCase = (*UserApplicationService)(nil)

func NewUserApplicationService(userDomainService *userdomainservice.UserDomainService, roleDomainService *roledomainservice.RoleDomainService) *UserApplicationService {
	createUserHandler := command.NewCreateUserHandler(roleDomainService, userDomainService)
	updateUserHandler := command.NewUpdateUserHandler(roleDomainService, userDomainService)
	deleteUserHandler := command.NewDeleteUserHandler(userDomainService)

	getUserByIDHandler := query.NewGetUserByIDHandler(userDomainService)

	userCommands := command.NewUserCommands(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
	)
	userQueries := query.NewUserQueries(getUserByIDHandler)

	return &UserApplicationService{
		Commands: userCommands,
		Queries:  userQueries,
	}
}

func (svc *UserApplicationService) CreateUser(ctx context.Context, cmd *command.CreateUserCommand) error {
	return svc.Commands.CreateUser.Handle(ctx, cmd)
}

func (svc *UserApplicationService) GetByID(ctx context.Context, q *query.GetUserByIDQuery) (*query.GetUserByIDQueryResult, error) {
	return svc.Queries.GetUserByID.Handle(ctx, q)
}

func (svc *UserApplicationService) UpdateUser(ctx context.Context, cmd *command.UpdateUserCommand) error {
	return svc.Commands.UpdateUser.Handle(ctx, cmd)
}

func (svc *UserApplicationService) DeleteUser(ctx context.Context, cmd *command.DeleteUserCommand) error {
	return svc.Commands.DeleteUser.Handle(ctx, cmd)
}
