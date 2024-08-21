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

func NewVerifyCommand(tokenValue string) *VerifyTokenCommand {
	return &VerifyTokenCommand{
		TokenValue: tokenValue,
	}
}

type VerifyTokenCommandResult struct {
	UserID   int
	Username string
	RoleID   int
}

type AuthCommands struct {
	Login       LoginCommandHandler
	Logout      LogoutCommandHandler
	VerifyToken VerifyTokenCommandHandler
}

func NewAuthCommands(
	login LoginCommandHandler,
	logout LogoutCommandHandler,
	verifyToken VerifyTokenCommandHandler,
) *AuthCommands {
	return &AuthCommands{
		Login:       login,
		Logout:      logout,
		VerifyToken: verifyToken,
	}
}
