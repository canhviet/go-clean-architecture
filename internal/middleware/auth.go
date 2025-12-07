package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/canhviet/go-clean-architecture/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   uint
	Issuer   string
	Audience string
}

func IssueTokens(userID uint) (*Tokens, error) {
	now := time.Now().UTC()
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "go-clean-app",
		Audience: "web-client",
	}

	// Access Token
	accessClaims := jwt.RegisteredClaims{
		Subject:   userIDStr,                             
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		Subject:   userIDStr,                            
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	t.Access = accessSigned
	t.Refresh = refreshSigned
	return t, nil
}

func Persist(ctx context.Context, r *repository.Redis, t *Tokens) error {
    userIDStr := strconv.FormatUint(uint64(t.UserID), 10)

    // Lưu access token JTI
    if err := r.SetJTI(
        ctx,
        "access:"+t.JTIAcc,
        userIDStr,           
        t.ExpAcc,           
    ); err != nil {
        return err
    }

    // Lưu refresh token JTI
    if err := r.SetJTI(
        ctx,
        "refresh:"+t.JTIRef,
        userIDStr,           
        t.ExpRef,            
    ); err != nil {
        return err
    }

    return nil
}

// SetAuthCookies - HttpOnly + Secure + SameSite
func SetAuthCookies(c *gin.Context, t *Tokens) {
	maxAgeAcc := int(time.Until(t.ExpAcc).Seconds())
	maxAgeRef := int(time.Until(t.ExpRef).Seconds())

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", t.Access, maxAgeAcc, "/", "", true, true)  // HttpOnly + Secure
	c.SetCookie("refresh_token", t.Refresh, maxAgeRef, "/", "", true, true)
}

// ClearAuthCookies
func ClearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}

// ParseAccess & ParseRefresh
func ParseAccess(tokenStr string) (*jwt.RegisteredClaims, error) {
	return parseToken(tokenStr, os.Getenv("ACCESS_SECRET"))
}

func ParseRefresh(tokenStr string) (*jwt.RegisteredClaims, error) {
	return parseToken(tokenStr, os.Getenv("REFRESH_SECRET"))
}

// Helper chung
func parseToken(tokenStr, secret string) (*jwt.RegisteredClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience("web-client"),
		jwt.WithIssuer("go-clean-app"),
	).ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}