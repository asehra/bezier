package webserver

import (
	"net/http"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/service"
	"github.com/gin-gonic/gin"
)

type CreateMerchantResponse struct {
	MerchantID string `json:"merchant_id"`
	Error      error  `json:"error"`
}

func createMerchantHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		merchantID, err := service.CreateMerchant(config.DB, config.MerchantIDGenerator)
		c.JSON(http.StatusOK, CreateMerchantResponse{merchantID, err})
	}
}

type AuthorizeTransactionRequest struct {
	CardNumber int64  `json:"card_number"`
	MerchantID string `json:"merchant_id"`
	Amount     int32  `json:"amount"`
}
type AuthorizeTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	Error         string `json:"error"`
}

func authorizeTransactionHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		var params AuthorizeTransactionRequest
		if err := c.ShouldBindJSON(&params); err != nil {
			c.String(http.StatusBadRequest, `{"error":"bad JSON format"}`)
			return
		}
		transactionID, err := service.AuthorizeTransaction(
			config.DB,
			params.CardNumber,
			params.MerchantID,
			params.Amount,
			config.TransactionIDGenerator,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, AuthorizeTransactionResponse{"", err.Error()})
			return
		}
		c.JSON(http.StatusOK, AuthorizeTransactionResponse{transactionID, ""})
	}
}

type MerchantTransactionsResponse struct {
	Merchant model.Merchant `json:"merchant_activity"`
	Error    string         `json:"error"`
}

func merchantTransactionsHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		merchantID := c.Query("merchant_id")
		if merchantID == "" {
			c.JSON(http.StatusBadRequest, MerchantTransactionsResponse{model.Merchant{}, "Bad merchant ID format"})
			return
		}
		merchant, err := service.GetMerchant(config.DB, merchantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, MerchantTransactionsResponse{merchant, err.Error()})
			return
		}
		c.JSON(http.StatusOK, MerchantTransactionsResponse{merchant, ""})
	}
}
