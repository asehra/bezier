package webserver

import (
	"net/http"

	"github.com/asehra/bezier/config"
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
}

func authorizeTransactionHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		var params AuthorizeTransactionRequest
		if err := c.ShouldBindJSON(&params); err != nil {
			c.String(http.StatusBadRequest, `{"error":"bad JSON format"}`)
			return
		}
		transactionID := service.AuthorizeTransaction(
			config.DB,
			params.CardNumber,
			params.MerchantID,
			params.Amount,
			config.TransactionIDGenerator,
		)
		c.JSON(http.StatusOK, AuthorizeTransactionResponse{transactionID})
	}
}
