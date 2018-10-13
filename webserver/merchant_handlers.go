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
