package model

type Merchant struct {
	ID           string        `json:"id"`
	Transactions []Transaction `json:"authorized_transactions"`
}
