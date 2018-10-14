package model

type Merchant struct {
	ID                     string        `json:"id"`
	AuthorizedTransactions []Transaction `json:"authorized_transactions"`
}
