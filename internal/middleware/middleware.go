package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/canhviet/go-clean-architecture/internal/repository"
	"github.com/gin-gonic/gin"
)

func bearerFromHeader(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

func AuthMiddleware(r *repository.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, _ := c.Cookie("access_token")
		if tokenStr == "" {
			tokenStr = bearerFromHeader(c)
		}
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := ParseAccess(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		userIDFromRedis, err := r.GetUserByJTI(c.Request.Context(), "access:"+claims.ID)
		if err != nil || userIDFromRedis == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			return
		}

		c.Set("userID", claims.Subject) 
		c.Next()
	}
}

func MustCookie(c *gin.Context, name string) (string, error) {
	val, err := c.Cookie(name)
	if err != nil || val == "" {
		return "", errors.New("missing cookie: " + name)
	}
	return val, nil
}

func SetAuthResponse(c *gin.Context, toks *Tokens) {
	SetAuthCookies(c, toks)

	c.Header("Authorization", "Bearer "+toks.Access)
	
	c.Header("X-Access-Token", toks.Access)
	c.Header("X-Refresh-Token", toks.Refresh)
}