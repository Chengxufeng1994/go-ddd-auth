package facade

import (
	"fmt"
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

	cmd := command.NewUpdateUserCommand(
		int(userID),
		command.UpdateUserOpt{
			Username: &dto.Username,
			Password: &dto.Password,
			RoleID:   &dto.RoleID,
		})
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

func (api *UserApi) SearchUsers(c *gin.Context) {
	queryPage := c.DefaultQuery("page", "1")
	querySize := c.DefaultQuery("size", "10")
	queryOrderBy := c.DefaultQuery("orderBy", "id")
	querySortBy := c.DefaultQuery("sortBy", "asc")
	searchText := c.DefaultQuery("q", "")
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(querySize)
	if err != nil {
		size = 0
	}

	q := query.NewSearchUsersQuery(page, size, queryOrderBy, querySortBy, searchText)
	fmt.Printf("%#v", q)
	res, err := api.userApplicationService.SearchUsers(c.Request.Context(), q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []*dto.UserDto
	for _, user := range res.Users {
		users = append(users, api.userAssembler.ToDTO(user))
	}

	c.JSON(http.StatusOK, &dto.SearchUserResponseDto{
		Users:      users,
		Pagination: res.Pagination,
	})
}
