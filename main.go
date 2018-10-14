package main

import (
	"net/http"
	"os"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/storage"
	"github.com/asehra/bezier/webserver"
	"github.com/gin-gonic/gin"
)

var port string

func main() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	server := webserver.Create(
		config.Config{
			DB:                     storage.NewInMemoryStore(),
			CardIDGenerator:        &generator.CardIDIncrementor{LastID: 4921000000000000},
			MerchantIDGenerator:    &generator.StringIDIncrementor{Prefix: "M", LastID: 1000},
			TransactionIDGenerator: &generator.StringIDIncrementor{Prefix: "TX", LastID: 10000},
			StdOut:                 os.Stdout,
			StdErr:                 os.Stderr,
		},
	)
	server.Static("/assets", "./public/assets")
	server.LoadHTMLGlob("public/*.html")
	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	server.Run(":" + port)
}
