package command

import "context"

type LogoutCommandHandler interface {
	Handle(ctx context.Context, cmd *LogoutCommand) (*LogoutCommandResult, error)
}

type logoutHandler struct{}

var _ LogoutCommandHandler = (*logoutHandler)(nil)

func NewLogoutHandler() LogoutCommandHandler {
	return &logoutHandler{}
}

func (h *logoutHandler) Handle(ctx context.Context, cmd *LogoutCommand) (*LogoutCommandResult, error) {
	panic("unimplemented")
}
