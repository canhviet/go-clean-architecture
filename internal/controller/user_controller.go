package controller

import (
	"net/http"
	"strconv"

	"github.com/canhviet/go-clean-architecture/internal/dto"
	"github.com/canhviet/go-clean-architecture/internal/service"
	"github.com/canhviet/go-clean-architecture/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
    service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
    return &UserController{service: service}
}

// GET /users
func (c *UserController) GetAll(ctx *gin.Context) {
    users, err := c.service.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"data": users})
}

// GET /users/:id
func (c *UserController) GetByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.ParseUint(idStr, 10, 32)

    user, err := c.service.GetByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// POST /users
func (c *UserController) Create(ctx *gin.Context) {
    var req dto.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := c.service.Create(req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("User created", zap.Uint("id", user.ID))
    ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

// PUT /users/:id
func (c *UserController) Update(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.ParseUint(idStr, 10, 32)

    var req dto.UpdateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := c.service.Update(uint(id), req)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// DELETE /users/:id
func (c *UserController) Delete(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.ParseUint(idStr, 10, 32)

    if err := c.service.Delete(uint(id)); err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}