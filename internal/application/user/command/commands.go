package command

import (
	"context"
)

type Commands interface {
	CreateUser(ctx context.Context, cmd *CreateUserCommand) error
}

type CreateUserCommand struct {
	Username string
	Password string
	RoleID   int
}

func NewCreateUserCommand(username string, password string, roleID int) *CreateUserCommand {
	return &CreateUserCommand{
		Username: username,
		Password: password,
		RoleID:   roleID,
	}
}

type UpdateUserCommand struct {
	UserID   int
	Username string
	Password string
	RoleID   int
}

func NewUpdateUserCommand(userID, roleID int, username, password string) *UpdateUserCommand {
	return &UpdateUserCommand{
		UserID:   userID,
		Username: username,
		Password: password,
		RoleID:   roleID,
	}
}

type DeleteUserCommand struct {
	UserID int
}

func NewDeleteUserCommand(userID int) *DeleteUserCommand {
	return &DeleteUserCommand{
		UserID: userID,
	}
}

type UserCommands struct {
	CreateUser CreateUserCommandHandler
	UpdateUser UpdateUserCommandHandler
	DeleteUser DeleteUserCommandHandler
}

func NewUserCommands(
	createUser CreateUserCommandHandler,
	updateUser UpdateUserCommandHandler,
	deleteUser DeleteUserCommandHandler,
) *UserCommands {
	return &UserCommands{
		CreateUser: createUser,
		UpdateUser: updateUser,
		DeleteUser: deleteUser,
	}
}
