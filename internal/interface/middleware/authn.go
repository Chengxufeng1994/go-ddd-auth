package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth/command"
)

type Authenticate struct {
	authApplicationService authapplication.AuthUseCase
}

func NewAuthenticateMiddleware(
	authenticateService authapplication.AuthUseCase,
) *Authenticate {
	return &Authenticate{
		authApplicationService: authenticateService,
	}
}

func (m *Authenticate) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if !strings.Contains(authorization, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		tokenParts := strings.Split(authorization, " ")
		if len(tokenParts) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		tokenValue := tokenParts[1]
		cmd := command.NewVerifyTokenCommand(tokenValue)
		cmdResult, err := m.authApplicationService.VerifyToken(c, cmd)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		m.updateValidatedDataToContext(c, cmdResult)

		c.Next()
	}
}

type UserValidateData struct {
	UserID   int
	RoleID   int
	Username string
}

func (m *Authenticate) updateValidatedDataToContext(c *gin.Context, cmd *command.VerifyTokenCommandResult) {
	userValidateData, _ := c.Value(CtxUserData).(*UserValidateData)
	if userValidateData == nil {
		userValidateData = &UserValidateData{}
	}
	userValidateData.UserID = cmd.UserID
	userValidateData.RoleID = cmd.RoleID
	userValidateData.Username = cmd.Username
	c.Set(CtxUserData, userValidateData)
}
