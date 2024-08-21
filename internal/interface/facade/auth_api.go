package facade

import (
	"net/http"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth/command"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/dto"
	"github.com/gin-gonic/gin"
)

type AuthApi struct {
	authApplicationService auth.AuthUseCase
}

func NewAuthApi(authApplicationService auth.AuthUseCase) *AuthApi {
	return &AuthApi{
		authApplicationService: authApplicationService,
	}
}

func (api *AuthApi) Login(c *gin.Context) {
	var req dto.LoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.NewLoginCommand(req.Username, req.Password)
	cmdRes, err := api.authApplicationService.Login(c.Request.Context(), cmd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &dto.LoginResponseDto{
		AccessToken:  cmdRes.AccessToken,
		RefreshToken: cmdRes.RefreshToken,
	}

	c.JSON(http.StatusOK, resp)
}

func (api *AuthApi) Logout(c *gin.Context) {
}
