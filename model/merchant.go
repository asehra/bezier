package model

type Transaction struct {
	ID         string
	CardNumber int64
	Amount     int32
}

type Merchant struct {
	ID                     string
	AuthorizedTransactions []Transaction
}
