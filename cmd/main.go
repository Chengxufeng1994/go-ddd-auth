package main

import (
	authapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth"

	roledao "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/dao"
	rolepersistence "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/persistence"
	rolepo "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/repository/po"
	roledomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/role/service"

	tokenpersistence "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/repository/persistence"
	tokendomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/token/service"

	userapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/user"
	userdao "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/dao"
	userpersistence "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/persistence"
	userpo "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/repository/po"
	userdomainservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/user/service"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructue/badger"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructue/token/jwt"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/middleware"

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

	badgerDB, err := badger.NewBadger("", badger.WithPath("/tmp/badger"))
	if err != nil {
		panic(err)
	}
	defer badgerDB.Close()

	roleDao := roledao.NewRoleDao(db)
	roleRepository := rolepersistence.NewRoleRepository(db, roleDao)
	roleDomainService := roledomainservice.NewRoleDomainService(roleRepository)

	userDao := userdao.NewUserDao(db)
	userRepository := userpersistence.NewUserRepository(db, userDao)
	userDomainService := userdomainservice.NewUserDomainService(userRepository)
	userApplicationService := userapplication.NewUserApplicationService(userDomainService, roleDomainService)

	secretKey := []byte("secret")
	tokenEnhancer := jwt.NewJwtTokenEnhancer(secretKey)
	tokenRepository := tokenpersistence.NewTokenRepository(badgerDB)

	authDomainService := tokendomainservice.NewTokenDomainService(tokenEnhancer, tokenRepository)
	authApplicationService := authapplication.NewAuthApplicationService(authDomainService, userDomainService)

	userApi := facade.NewUserApi(userApplicationService)
	authApi := facade.NewAuthApi(authApplicationService)

	authMiddleware := middleware.NewAuthenticateMiddleware(authApplicationService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	// public route
	auth := api.Group("/auth")
	auth.POST("/login", authApi.Login)
	auth.POST("/logout", authApi.Logout)
	auth.POST("/refresh_token")

	// private route
	user := api.Group("/users")
	user.Use(authMiddleware.Middleware())
	user.POST("/", userApi.CreateUser)
	user.GET("/:user_id", userApi.GetUserByID)
	user.PUT("/:user_id", userApi.UpdateUser)
	user.DELETE("/:user_id", userApi.DeleteUserByID)

	r.Run("0.0.0.0:8080")
}
