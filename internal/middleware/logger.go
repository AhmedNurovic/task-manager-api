package middleware

import (
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        logger.Info("Request handled",
            zap.Int("status", c.Writer.Status()),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
        )
    }
}