package model

type Transaction struct {
	ID         string `json:"id"`
	CardNumber int64  `json:"card_number"`
	Authorized int    `json:"amount"`
	Captured   int    `json:"captured"`
}
