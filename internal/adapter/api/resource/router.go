package resource

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPHandler) Routes(router *gin.Engine) {
	apirouter := router.Group("api/v1")
	apirouter.GET("/healthcheck", s.HealthCheck)
	apirouter.POST("/register", s.Register)
	apirouter.POST("/login", s.Login)

	auth := apirouter.Group("/auth")
	auth.Use(s.AuthMiddleware())

	router.NoRoute(func(c *gin.Context) { c.JSON(404, "no route") })
}
