package command

type LoginCommand struct {
	Username string
	Password string
}

func NewLoginCommand(username, password string) *LoginCommand {
	return &LoginCommand{
		Username: username,
		Password: password,
	}
}

type LoginCommandResult struct {
	AccessToken  string
	RefreshToken string
}

type LogoutCommand struct{}

type LogoutCommandResult struct{}

type VerifyTokenCommand struct {
	TokenValue string
}

func NewVerifyTokenCommand(tokenValue string) *VerifyTokenCommand {
	return &VerifyTokenCommand{
		TokenValue: tokenValue,
	}
}

type VerifyTokenCommandResult struct {
	UserID   int
	Username string
	RoleID   int
}
type VerifyPermissionCommand struct {
	RoleID   int
	Resource string
	Action   string
}

func NewVerifyPermissionCommand(roleID int, resource, action string) *VerifyPermissionCommand {
	return &VerifyPermissionCommand{
		RoleID:   roleID,
		Resource: resource,
		Action:   action,
	}
}

type VerifyPermissionCommandResult struct {
}

type AuthCommands struct {
	Login            LoginCommandHandler
	Logout           LogoutCommandHandler
	VerifyToken      VerifyTokenCommandHandler
	VerifyPermission VerifyPermissionCommandHandler
}

func NewAuthCommands(
	login LoginCommandHandler,
	logout LogoutCommandHandler,
	verifyToken VerifyTokenCommandHandler,
	verifyPermission VerifyPermissionCommandHandler,
) *AuthCommands {
	return &AuthCommands{
		Login:            login,
		Logout:           logout,
		VerifyToken:      verifyToken,
		VerifyPermission: verifyPermission,
	}
}
