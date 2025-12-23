package controller

import (
	"net/http"
	"time"

	"github.com/canhviet/go-clean-architecture/internal/dto"
	"github.com/canhviet/go-clean-architecture/internal/middleware"
	"github.com/canhviet/go-clean-architecture/internal/model"
	"github.com/canhviet/go-clean-architecture/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginHandler(c *gin.Context, db *gorm.DB, rds *repository.Redis) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    // Bind + validate JSON body
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": "Invalid JSON or missing fields",
        })
        return
    }

    var user model.User
    if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    toks, err := middleware.IssueTokens(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not issue tokens"})
        return
    }

    if err := middleware.Persist(c, rds, toks); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not persist tokens"})
        return
    }

    middleware.SetAuthResponse(c, toks)

    c.JSON(http.StatusOK, gin.H{
        "ok":      true,
        "user_id": user.ID,
        "access_token":  toks.Access, 
		"refresh_token": toks.Refresh,  
		"expires_in":    int(time.Until(toks.ExpAcc).Seconds()),
    })
}

func RegisterHandler(c *gin.Context, db *gorm.DB, rds *repository.Redis) {
	var input dto.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"details": err.Error(),
		})
		return
	}

	var existingUser model.User
	if err := db.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username or email already taken",
		})
		return
	}

	user := model.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		Name:     input.Name,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	toks, err := middleware.IssueTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not issue tokens"})
		return
	}

	if err := middleware.Persist(c.Request.Context(), rds, toks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not persist tokens"})
		return
	}

	middleware.SetAuthResponse(c, toks)

	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"ok":      true,
		"message": "Register successful",
		"user":    user,
		"access_token":  toks.Access,
		"refresh_token": toks.Refresh,
		"expires_in":    int(time.Until(toks.ExpAcc).Seconds()),
	})
}