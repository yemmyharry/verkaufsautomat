package resource

import (
	"github.com/gin-gonic/gin"
	ports "verkaufsautomat/internal/ports/resource"
)

type HTTPHandler struct {
	MachineService ports.MachineService
}

func (s *HTTPHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c)
		if err != nil {
			c.JSON(401, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func NewHTTPHandler(MachineService ports.MachineService) *HTTPHandler {
	handler := &HTTPHandler{
		MachineService: MachineService,
	}
	return handler
}
