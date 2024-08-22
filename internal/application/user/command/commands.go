package command

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

type UpdateUserOpt struct {
	Username *string
	Password *string
	RoleID   *int
}

type UpdateUserCommand struct {
	userID int
	opt    UpdateUserOpt
}

func NewUpdateUserCommand(userID int, opt UpdateUserOpt) *UpdateUserCommand {
	return &UpdateUserCommand{
		userID: userID,
		opt:    opt,
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
