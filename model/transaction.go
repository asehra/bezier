package model

type Transaction struct {
	ID         string `json:"id"`
	CardNumber int64  `json:"card_number"`
	Authorized int32  `json:"amount"`
	Captured   int32  `json:"captured"`
}
