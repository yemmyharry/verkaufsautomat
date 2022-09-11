package resource

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
	models "verkaufsautomat/internal/core/domain/resource"
	"verkaufsautomat/internal/core/logger"
)

type HealthCheckResponse struct {
	Status string
}

func (s *HTTPHandler) HealthCheck(c *gin.Context) {
	logger.Info("HealthCheck called")
	c.JSON(200, HealthCheckResponse{Status: "OK"})
}

func TokenValid(c *gin.Context) error {
	cookie, err := c.Cookie("myToken")
	if err != nil {
		return err
	}

	token, err := verifyToken(cookie)
	if err != nil {
		return err
	}

	if !isTokenValid(token) {
		return fmt.Errorf("token is not valid")
	}

	return nil
}

func (s *HTTPHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hashPassword)

	if err := s.MachineService.Register(&user); err != nil {
		logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user registered"})
}

func (s *HTTPHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("Error binding json: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := generateToken(&user)
	if err != nil {
		logger.Error("Error generating token: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("myToken", token, 3600, "/", "*", false, true)

	if err := s.MachineService.Login(&user); err != nil {
		logger.Error("Error logging in: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user logged in"})
}

func ComparePassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func generateToken(user *models.User) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := t.SignedString([]byte("secret"))
	if err != nil {
		logger.Error("Error signing token: " + err.Error())
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func isTokenValid(token *jwt.Token) bool {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	}

	return true
}
