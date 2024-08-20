package main

import (
	roledao "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/dao"
	rolepersistence "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/persistence"
	rolepo "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/po"
	roledomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/service"

	userapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/user"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/dao"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/persistence"
	userpo "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/po"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/facade"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=P@ssw0rd dbname=postgres host=10.1.5.7 port=31970 sslmode=disable TimeZone=UTC search_path=go_ddd_auth",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&rolepo.Role{}, &userpo.User{})

	roleDao := roledao.NewRoleDao(db)
	roleRepository := rolepersistence.NewRoleRepository(db, roleDao)
	roleDomainService := roledomainservice.NewRoleDomainService(roleRepository)

	userDao := dao.NewUserDao(db)
	userRepository := persistence.NewUserRepository(db, userDao)
	userDomainService := userdomainservice.NewUserDomainService(userRepository)
	userApplicationService := userapplication.NewUserApplicationService(userDomainService, roleDomainService)
	userApi := facade.NewUserApi(userApplicationService)

	r := gin.Default()
	user := r.Group("/users")
	user.POST("/", userApi.CreateUser)
	user.GET("/:user_id", userApi.GetUserByID)
	user.PUT("/:user_id", userApi.UpdateUser)
	user.DELETE("/:user_id", userApi.DeleteUserByID)
	r.Run("0.0.0.0:8080")
}
