package resource

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
	"verkaufsautomat/internal/adapter/repositories/mysql/resource"
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
	tokenString := c.Request.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := verifyToken(tokenString)
	if err != nil {
		logger.Error("Error verifying token: " + err.Error())
		return err
	}

	if !isTokenValid(token) {
		logger.Error("Token is not valid")
		return err
	}

	return nil
}

func (s *HTTPHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("Error binding json: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hashPassword)

	if err := s.MachineService.Register(&user); err != nil {
		logger.Error("Error registering user: " + err.Error())
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

	userid, roleid, err := resource.NewMachineRepositoryDB().GetUserIdAndRoleId(user.Username)
	if err != nil {
		return
	}

	user.UserID = userid
	user.RoleID = roleid

	token, err := generateToken(&user)
	if err != nil {
		logger.Error("Error generating token: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+token)

	if err := s.MachineService.Login(&user); err != nil {
		logger.Error("Error logging in: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user logged in"})
}

func generateToken(user *models.User) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.UserID,
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

func (s *HTTPHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Error("Error binding json: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) != 2 {
		logger.Error("Error getting token from header")
		c.JSON(400, gin.H{"error": "Error getting token from header"})
		return
	}

	token, err := verifyToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		logger.Error("Error verifying token: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	roleID := int(claims["role_id"].(float64))
	userID := int(claims["user_id"].(float64))

	if roleID != 2 {
		logger.Error("User cannot create product")
		c.JSON(400, gin.H{"error": "user cannot create product"})
		return
	}

	product.SellerID = uint(userID)

	if err := s.MachineService.CreateProduct(&product); err != nil {
		logger.Error("Error creating product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "product created"})
}

func (s *HTTPHandler) GetProducts(c *gin.Context) {
	products, err := s.MachineService.GetProducts()
	if err != nil {
		logger.Error("Error getting products: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}

func (s *HTTPHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error converting id to int: " + err.Error())
		return
	}

	product, err := s.MachineService.GetProductById(atoi)
	if err != nil {
		logger.Error("Error getting product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (s *HTTPHandler) UpdateProduct(c *gin.Context) {

	id := c.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error converting id to int: " + err.Error())
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Error("Error binding json: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) != 2 {
		logger.Error("Error getting token from header")
		c.JSON(400, gin.H{"error": "Error getting token from header"})
		return
	}

	token, err := verifyToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		logger.Error("Error verifying token: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	roleID := int(claims["role_id"].(float64))
	userID := int(claims["user_id"].(float64))

	if roleID != 2 {
		logger.Error("User cannot update product")
		c.JSON(400, gin.H{"error": "user cannot update product"})
		return
	}

	product.SellerID = uint(userID)

	if err := s.MachineService.UpdateProductByID(atoi, &product); err != nil {
		logger.Error("Error updating product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "product updated"})
}

func (s *HTTPHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error converting id to int: " + err.Error())
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) != 2 {
		logger.Error("Error getting token from header")
		c.JSON(400, gin.H{"error": "Error getting token from header"})
		return
	}

	token, err := verifyToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		logger.Error("Error verifying token: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	roleID := int(claims["role_id"].(float64))

	if roleID != 2 {
		logger.Error("User cannot delete product")
		c.JSON(400, gin.H{"error": "user cannot delete product"})
		return
	}

	if err := s.MachineService.DeleteProductByID(atoi); err != nil {
		logger.Error("Error deleting product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "product deleted"})
}
