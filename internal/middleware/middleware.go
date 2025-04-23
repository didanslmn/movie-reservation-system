package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"slices"

	"github.com/didanslmn/movie-reservation-system.git/internal/pkg/logger"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type contextKey string

const userContextKey contextKey = "user"

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logAndAbort(c, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			logAndAbort(c, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logAndAbort(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logAndAbort(c, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			logAndAbort(c, http.StatusUnauthorized, "Invalid sub claim")
			return
		}

		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		roleStr, _ := claims["role"].(string)

		user := &model.User{
			Model: gorm.Model{ID: uint(sub)},
			Name:  name,
			Email: email,
			Role:  model.Role(roleStr),
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userContextKey, user))
		c.Set("userID", user.ID)

		logger.Log.Info("Authenticated user", zap.Uint("userID", user.ID), zap.String("email", user.Email))

		c.Next()
	}
}

func RoleBasedAccess(allowedRoles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := GetUserFromContext(c.Request.Context())
		if !ok {
			logAndAbort(c, http.StatusUnauthorized, "User not found in context")
			return
		}

		if !slices.Contains(allowedRoles, user.Role) {
			logger.Log.Warn("Access denied", zap.Uint("userID", user.ID), zap.String("role", string(user.Role)))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.Next()
	}
}

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}

func logAndAbort(c *gin.Context, code int, msg string) {
	logger.Log.Warn("Auth error", zap.Int("status", code), zap.String("path", c.Request.URL.Path), zap.String("message", msg))
	c.AbortWithStatusJSON(code, gin.H{"error": msg})
}
