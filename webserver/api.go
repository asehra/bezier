package webserver

import (
	"github.com/asehra/bezier/config"
	"github.com/gin-gonic/gin"
)

func Create(config config.Config) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = config.StdOut
	gin.DefaultErrorWriter = config.StdErr

	v1 := r.Group("/v1")
	{
		v1.GET("/card/create", createCardHandler(config))
		v1.GET("/card/details", getCardHandler(config))
		v1.POST("/card/top-up", topUpCardHandler(config))
		v1.GET("/merchant/create", createMerchantHandler(config))
		v1.POST("/merchant/authorize-transaction", authorizeTransactionHandler(config))
		v1.GET("/merchant/transactions", merchantTransactionsHandler(config))
		v1.POST("/merchant/capture-transaction", captureTransactionHandler(config))
		v1.POST("/merchant/reverse-transaction", reverseTransactionHandler(config))
	}
	return r
}
