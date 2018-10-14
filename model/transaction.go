package model

type Transaction struct {
	ID         string `json:"id"`
	CardNumber int64  `json:"card_number"`
	Authorized int    `json:"authorized"`
	Captured   int    `json:"captured"`
	Reversed   int    `json:"reversed"`
	Refunded   int    `json:"refunded"`
}
