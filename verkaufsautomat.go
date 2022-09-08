package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	adapter "verkaufsautomat/internal/adapter/api/resource"
	"verkaufsautomat/internal/adapter/repositories/mysql/resource"
	"verkaufsautomat/internal/core/logger"
	services "verkaufsautomat/internal/core/services/resource"
)

func main() {
	err := godotenv.Load("verkaufsautomat.env")
	if err != nil {
		logger.Error("Error loading .env file")
	}
	router := gin.Default()
	database := resource.NewMachineRepositoryDB()
	service := services.New(database)
	handler := adapter.NewHTTPHandler(service)
	handler.Routes(router)
	logger.Info("Starting server on port 8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
