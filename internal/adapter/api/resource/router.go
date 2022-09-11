package resource

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *HTTPHandler) Routes(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apirouter := router.Group("api/v1")
	apirouter.GET("/healthcheck", s.HealthCheck)
	apirouter.POST("/register", s.Register)
	apirouter.POST("/login", s.Login)

	auth := router.Group("/auth")
	auth.Use(s.AuthMiddleware())
	auth.POST("/create_product", s.CreateProduct)
	auth.GET("/get_products", s.GetProducts)
	auth.GET("/get_product/:id", s.GetProduct)
	auth.PUT("/update_product/:id", s.UpdateProduct)
	auth.DELETE("/delete_product/:id", s.DeleteProduct)
	auth.PATCH("deposit_money", s.DepositMoney)
	auth.POST("/buy_product", s.BuyProduct)
	auth.PATCH("reset_deposit", s.ResetDeposit)
	router.NoRoute(func(c *gin.Context) { c.JSON(404, "no route") })
}
