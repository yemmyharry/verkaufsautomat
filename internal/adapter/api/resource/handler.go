package resource

import (
	"errors"
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

func insertCoin(coin int) (bool, error) {
	if coin == 5 || coin == 10 || coin == 20 || coin == 50 || coin == 100 {
		logger.Info("Coin inserted")
		return true, nil
	}
	logger.Error("Invalid coin")
	return false, errors.New("invalid coin")
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

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)

	if err := s.MachineService.Login(&user); err != nil {
		logger.Error("Error logging in: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user logged in", "token": token})
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

func (s *HTTPHandler) DepositMoney(c *gin.Context) {

	var deposit struct {
		Amount int `json:"amount"`
	}

	if err := c.ShouldBindJSON(&deposit); err != nil {
		logger.Error("Error binding json: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := insertCoin(deposit.Amount)
	if err != nil {
		logger.Error("Error inserting coin: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) != 2 {
		logger.Error("Error getting token from header")

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

	if roleID != 1 {
		logger.Error("User cannot deposit money")
		c.JSON(400, gin.H{"error": "user cannot deposit money"})
		return
	}

	if err := s.MachineService.DepositMoney(userID, deposit.Amount); err != nil {
		logger.Error("Error depositing money: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "money deposited"})

}

func (s *HTTPHandler) BuyProduct(c *gin.Context) {

	var buyProduct struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	var response struct {
		TotalPrice int   `json:"total_price"`
		Change     []int `json:"change"`
		Quantity   int   `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&buyProduct); err != nil {
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

	if roleID != 1 {
		logger.Error("User cannot buy product")
		c.JSON(400, gin.H{"error": "user cannot buy product"})
		return
	}

	product, err := s.MachineService.GetProductById(buyProduct.ProductID)
	if err != nil {
		logger.Error("Error getting product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := s.MachineService.GetUserById(userID)
	if err != nil {
		logger.Error("Error getting user: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Deposit < product.Cost {
		logger.Error("User does not have enough money")
		c.JSON(400, gin.H{"error": "user does not have enough money"})
		return
	}

	if product.ProductID == 0 {
		logger.Error("Product does not exist")
		c.JSON(400, gin.H{"error": "product does not exist"})
		return
	}

	if product.AmountAvailable < buyProduct.Quantity {
		logger.Error("Product quantity is not enough")
		c.JSON(400, gin.H{"error": "product quantity is not enough"})
		return
	}

	totalPrice := product.Cost * buyProduct.Quantity

	change := user.Deposit - totalPrice

	user.Deposit = change
	if err := s.MachineService.UpdateUser(user); err != nil {
		logger.Error("Error updating user: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	product.AmountAvailable = product.AmountAvailable - buyProduct.Quantity

	if err := s.MachineService.UpdateProductByID(int(product.ProductID), &product); err != nil {
		logger.Error("Error updating product: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response.TotalPrice = totalPrice
	response.Change = getChange(change)
	response.Quantity = buyProduct.Quantity

	c.JSON(200, response)

}

func getChange(amount int) []int {

	coins := []int{100, 50, 20, 10, 5}
	change := []int{}

	for _, coin := range coins {
		for amount >= coin {
			amount -= coin
			change = append(change, coin)
		}
	}

	return change
}

func (s *HTTPHandler) ResetDeposit(context *gin.Context) {

	tokenString := context.Request.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) != 2 {
		logger.Error("Error getting token from header")
		context.JSON(400, gin.H{"error": "Error getting token from header"})
		return
	}

	token, err := verifyToken(strings.Split(tokenString, " ")[1])
	if err != nil {
		logger.Error("Error verifying token: " + err.Error())
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	roleID := int(claims["role_id"].(float64))
	userID := int(claims["user_id"].(float64))

	if roleID != 1 {
		logger.Error("User cannot reset deposit")
		context.JSON(400, gin.H{"error": "user cannot reset deposit"})
		return
	}

	user, err := s.MachineService.GetUserById(userID)
	if err != nil {
		logger.Error("Error getting user: " + err.Error())
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user.Deposit = 0
	if err := s.MachineService.UpdateUser(user); err != nil {
		logger.Error("Error updating user: " + err.Error())
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "deposit reset"})
}
