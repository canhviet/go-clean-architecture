package main

import (
	"github.com/canhviet/go-clean-architecture/internal/config"
	"github.com/canhviet/go-clean-architecture/internal/controller"
	"github.com/canhviet/go-clean-architecture/internal/database"
	"github.com/canhviet/go-clean-architecture/internal/middleware"
	"github.com/canhviet/go-clean-architecture/internal/repository"
	"github.com/canhviet/go-clean-architecture/internal/service"
	"github.com/canhviet/go-clean-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	defer logger.Log.Sync()

	if err := database.Init(); err != nil {
		panic("Database connection failed: " + err.Error())
	}	

    userRepo := repository.NewUserRepository(database.DB)
    userService := service.NewUserService(userRepo)
    userCtrl := controller.NewUserController(userService)
	redisRepo := repository.NewRedis()

	r := gin.New()

	public := r.Group("/api/v1")
	{
		public.POST("/login", func(c *gin.Context) {
            controller.LoginHandler(c, database.DB, redisRepo)
        })  

		public.POST("/register", func(c *gin.Context) {
            controller.RegisterHandler(c, database.DB, redisRepo)
        }) 
	}

    auth := r.Group("/api/v1")
	auth.Use(middleware.AuthMiddleware(redisRepo))
	{
		auth.GET("/users", userCtrl.GetAll)
		auth.GET("/users/:id", userCtrl.GetByID)
		auth.PUT("/users/:id", userCtrl.Update)
		auth.DELETE("/users/:id", userCtrl.Delete)

		// auth.POST("/logout", controller.LogoutHandler)
		// auth.GET("/me", controller.GetMeHandler)
	}

	r.Run(":8080")
}