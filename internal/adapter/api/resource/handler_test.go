package resource

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"verkaufsautomat/internal/core/domain/resource"
	services "verkaufsautomat/internal/core/services/mock"
)

func TestApplication_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedService := services.NewMockMachineService(ctrl)
	handler := NewHTTPHandler(mockedService)

	router := gin.Default()

	handler.Routes(router)

	t.Run("Deposit money", func(t *testing.T) {
		deposit := struct {
			Amount int `json:"amount"`
		}{
			Amount: 100,
		}
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjMwMjU0MjYsInJvbGVfaWQiOjIsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiaGFycnkifQ.fQpdaclPKqcGIKB5ng5UVwstZBbwwgze00mq3asZF-g"
		mockedService.EXPECT().DepositMoney(1, deposit.Amount).Return(nil)
		m, _ := json.Marshal(deposit)
		req, err := http.NewRequest("PATCH", "/auth/deposit_money", strings.NewReader(string(m)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		if strings.Contains(response.Body.String(), "money deposited") {
			t.Log("Deposit money test passed")
		} else {
			t.Error("Deposit money test failed")
		}

		if response.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
		}

	})

}

func TestApplication_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedService := services.NewMockMachineService(ctrl)
	handler := NewHTTPHandler(mockedService)

	router := gin.Default()

	handler.Routes(router)

	t.Run("Create product", func(t *testing.T) {

		var product resource.Product

		product.ProductID = 1
		product.AmountAvailable = 10
		product.Cost = 100
		product.ProductName = "cocacola"
		product.SellerID = 1

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjMwMjU0MjYsInJvbGVfaWQiOjIsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiaGFycnkifQ.fQpdaclPKqcGIKB5ng5UVwstZBbwwgze00mq3asZF-g"
		mockedService.EXPECT().CreateProduct(&product).Return(nil)
		m, err := json.Marshal(product)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/auth/create_product", strings.NewReader(string(m)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		if strings.Contains(response.Body.String(), "product created") {
			t.Log("Create product test passed")
		} else {
			t.Error("Create product test failed")
		}

		if response.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
		}

	})

}
