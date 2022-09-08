package resource

import "github.com/gin-gonic/gin"

type HealthCheckResponse struct {
	Status string
}

func (s *HTTPHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, HealthCheckResponse{Status: "OK"})
}
