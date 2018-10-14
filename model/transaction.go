package model

type Transaction struct {
	ID         string `json:"id"`
	CardNumber int64  `json:"card_number"`
	Amount     int32  `json:"amount"`
}
