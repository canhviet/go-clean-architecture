package main

import (
	"github.com/canhviet/go-clean-architecture/internal/config"
	"github.com/canhviet/go-clean-architecture/internal/controller"
	"github.com/canhviet/go-clean-architecture/internal/database"
	"github.com/canhviet/go-clean-architecture/internal/repository"
	"github.com/canhviet/go-clean-architecture/internal/service"
	"github.com/canhviet/go-clean-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	logger.Init()
	defer logger.Log.Sync()
	r := gin.New()

	if err := database.Init(); err != nil {
		panic("Database connection failed: " + err.Error())
	}	

    userRepo := repository.NewUserRepository(database.DB)
    userService := service.NewUserService(userRepo)
    userCtrl := controller.NewUserController(userService)

    // Routes
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", userCtrl.GetAll)
        v1.GET("/users/:id", userCtrl.GetByID)
        v1.POST("/users", userCtrl.Create)
        v1.PUT("/users/:id", userCtrl.Update)
        v1.DELETE("/users/:id", userCtrl.Delete)
    }

    r.Run(":8080")
}