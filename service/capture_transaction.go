package service

import (
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func CaptureTransaction(db storage.Storage, merchantID string, transactionID string, amount int32) error {
	merchant, _ := db.GetMerchant(merchantID) // TODO Handle err
	idx := findTransactionIndex(merchant.Transactions, transactionID)
	if idx == -1 {
		panic("unhandled missing transaction case")
	}
	merchant.Transactions[idx].Authorized = merchant.Transactions[idx].Authorized - amount // TODO handle over-capture
	merchant.Transactions[idx].Captured = merchant.Transactions[idx].Captured + amount
	db.StoreMerchant(merchant) //TODO Handle err

	return nil
}

func findTransactionIndex(transactions []model.Transaction, transactionID string) int {
	for idx, transaction := range transactions {
		if transaction.ID == transactionID {
			return idx
		}
	}
	return -1
}
