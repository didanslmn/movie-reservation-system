package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/didanslmn/movie-reservation-system.git/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("panic recovered",
					zap.Any("error", r),
					zap.String("path", c.FullPath()),
					zap.String("method", c.Request.Method),
					zap.ByteString("stack", debug.Stack()))

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
