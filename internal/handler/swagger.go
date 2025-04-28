package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/http-swagger"
)

func SwaggerHandler() gin.HandlerFunc {
	return gin.WrapH(httpSwagger.WrapHandler)
}
