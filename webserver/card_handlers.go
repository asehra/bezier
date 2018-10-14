package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/service"
	"github.com/gin-gonic/gin"
)

type CreateCardResponse struct {
	CardNumber int64 `json:"card_number"`
	Error      error `json:"error"`
}

func createCardHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		cardNumber, _ := service.CreateCard(config.DB, config.CardIDGenerator)
		c.JSON(http.StatusOK, CreateCardResponse{cardNumber, nil})
	}
}

type GetCardResponse struct {
	Card  model.Card `json:"card_details"`
	Error string     `json:"error"`
}

func getCardHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		cardNumber, err := strconv.ParseInt(c.Query("card_number"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, GetCardResponse{model.Card{}, "Bad card number format"})
			return
		}
		card, err := service.GetCard(config.DB, cardNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, GetCardResponse{card, err.Error()})
			return
		}
		c.JSON(http.StatusOK, GetCardResponse{card, ""})
	}
}

type TopUpCardRequest struct {
	CardNumber int64 `json:"card_number"`
	Amount     uint  `json:"amount"`
}

func topUpCardHandler(config config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		var params TopUpCardRequest
		if err := c.ShouldBindJSON(&params); err != nil {
			c.String(http.StatusBadRequest, `{"error":"bad JSON format"}`)
			return
		}
		err := service.TopUpCard(config.DB, params.CardNumber, params.Amount)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
			return
		}
		c.String(http.StatusOK, "")
	}
}
