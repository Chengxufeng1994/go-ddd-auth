package facade

import (
	"net/http"
	"strconv"

	application "github.com/Chengxufeng1994/go-ddd-auth/internal/application/user"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/command"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/application/user/query"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/assembler"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/dto"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	userAssembler          *assembler.UserAssembler
	userApplicationService *application.UserApplicationService
}

func NewUserApi(userApplicationService *application.UserApplicationService) *UserApi {
	return &UserApi{
		userApplicationService: userApplicationService,
	}
}

func (api *UserApi) CreateUser(c *gin.Context) {
	var createUserDto dto.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.NewCreateUserCommand(createUserDto.Username, createUserDto.Password, createUserDto.RoleID)
	if err := api.userApplicationService.CreateUser(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (api *UserApi) GetUserByID(c *gin.Context) {
	paramUserID := c.Param("user_id")
	userID, _ := strconv.ParseInt(paramUserID, 10, 64)

	q := query.NewGetUserByIDQuery(int(userID))
	res, err := api.userApplicationService.GetByID(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto := api.userAssembler.ToDTO(res.User)
	c.JSON(http.StatusOK, dto)
}

func (api *UserApi) UpdateUser(c *gin.Context) {
	paramUserID := c.Param("user_id")
	userID, _ := strconv.ParseInt(paramUserID, 10, 64)
	var dto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.NewUpdateUserCommand(int(userID), dto.RoleID, dto.Username, dto.Password)
	if err := api.userApplicationService.UpdateUser(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (api *UserApi) DeleteUserByID(c *gin.Context) {
	paramUserID := c.Param("user_id")
	userID, _ := strconv.ParseInt(paramUserID, 10, 64)

	var dto dto.DeleteUser
	dto.ID = int(userID)

	cmd := command.NewDeleteUserCommand(dto.ID)
	if err := api.userApplicationService.DeleteUser(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
