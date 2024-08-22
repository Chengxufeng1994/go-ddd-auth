package middleware

import (
	"net/http"
	"strings"

	authapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth/command"
	"github.com/gin-gonic/gin"
)

type Authorize struct {
	authApplicationService authapplication.AuthUseCase
}

func NewAuthorizeMiddleware(
	authApplicationService authapplication.AuthUseCase,
) *Authorize {
	return &Authorize{
		authApplicationService: authApplicationService,
	}
}

func (m *Authorize) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userValidatedData, ok := c.Value(CtxUserData).(*UserValidateData)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, nil)
			return
		}
		path := c.Request.URL.Path
		method := c.Request.Method
		cmd := command.NewVerifyPermissionCommand(userValidatedData.RoleID, path, m.mappingAction(method))
		_, err := m.authApplicationService.VerifyPermission(c, cmd)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, nil)
			return
		}

		c.Next()
	}
}

func (m *Authorize) mappingAction(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PATCH", "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "none"
	}
}
