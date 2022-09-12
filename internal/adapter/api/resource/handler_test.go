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
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjMwMTE3MDcsInJvbGVfaWQiOjEsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiam9obiJ9.jbgrs1gySB42QzaoCTCtptzk4Vll_de_LMLzjleSd3E"
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

// test for buy product
func TestApplication_BuyProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedService := services.NewMockMachineService(ctrl)
	handler := NewHTTPHandler(mockedService)

	router := gin.Default()

	handler.Routes(router)

	t.Run("Buy product", func(t *testing.T) {
		BuyProduct := struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}{
			ProductID: 1,
			Quantity:  1,
		}

		mockedService.EXPECT().GetProductById(BuyProduct.ProductID).Return(resource.Product{}, nil)

		mockedService.EXPECT().GetUserById(1).Return(resource.User{}, nil)

	})

}
