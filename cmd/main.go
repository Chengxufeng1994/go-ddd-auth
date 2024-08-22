package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	authapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/auth"
	userapplication "github.com/Chengxufeng1994/go-ddd-auth/internal/application/user"

	iamservice "github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/service"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/cache/badger"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/cachelayer"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/dao"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/po"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/timerlayer"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/rbac"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/token/jwt"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/transaction"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/facade"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/interface/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initPO(db *gorm.DB) {
	if err := db.AutoMigrate(&po.Role{}, &po.User{}); err != nil {
		panic(err)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("data/license.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initPO(db)

	badgerDB, err := badger.NewBadger("", badger.WithPath("data/badger"))
	if err != nil {
		panic(err)
	}
	defer badgerDB.Close()

	casbinEnforcer, err := rbac.NewCasbinEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		panic(err)
	}
	defer casbinEnforcer.StopAutoLoadPolicy()

	casbinAdapter := rbac.NewCasbinAdapter(casbinEnforcer)

	dbFactory := transaction.NewDefaultDBFactory(db)
	trxMgr := transaction.NewGormTransactionManager(dbFactory)

	passwordDomainService := iamservice.NewPasswordDomainService()

	roleDao := dao.NewRoleDao(db)
	roleRepository := persistence.NewRbacRepository(db, roleDao)
	roleDomainService := iamservice.NewRbacDomainService(roleRepository, casbinAdapter)

	userDao := dao.NewUserDao(db)
	userRepository := timerlayer.NewTimerLayer(
		cachelayer.NewUserCacheLayer(
			persistence.NewUserRepository(dbFactory, userDao), badgerDB),
	)
	userDomainService := iamservice.NewUserDomainService(userRepository)
	userApplicationService := userapplication.NewUserApplicationService(userRepository, userDomainService, passwordDomainService, roleDomainService, trxMgr)

	secretKey := []byte("secret")
	tokenEnhancer := jwt.NewJwtTokenEnhancer(secretKey)
	tokenRepository := persistence.NewTokenRepository(badgerDB)

	authDomainService := iamservice.NewTokenDomainService(tokenEnhancer, tokenRepository)
	authApplicationService := authapplication.NewAuthApplicationService(authDomainService, userDomainService, passwordDomainService, roleDomainService, casbinAdapter)

	userApi := facade.NewUserApi(userApplicationService)
	authApi := facade.NewAuthApi(authApplicationService)

	authnMiddleware := middleware.NewAuthenticateMiddleware(authApplicationService)
	authzMiddleware := middleware.NewAuthorizeMiddleware(authApplicationService)

	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CORSMiddleware(),
	)

	api := r.Group("/api")
	// public route
	auth := api.Group("/auth")
	auth.POST("/login", authApi.Login)
	auth.POST("/logout", authApi.Logout)
	auth.POST("/refresh_token")

	// private route
	privateRoute := api.Group("/")
	privateRoute.Use(authnMiddleware.Middleware(), authzMiddleware.Middleware())
	user := privateRoute.Group("/users")
	user.POST("/", userApi.CreateUser)
	user.GET("/", userApi.SearchUsers)
	user.PUT("/:user_id", userApi.UpdateUser)
	user.GET("/:user_id", userApi.GetUserByID)
	user.DELETE("/:user_id", userApi.DeleteUserByID)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")
}
