package service

import (
	"errors"

	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func CaptureTransaction(db storage.Storage, merchantID string, transactionID string, amount uint) error {
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return err
	}
	idx := findTransactionIndex(merchant.Transactions, transactionID)
	if idx == -1 {
		return errors.New("transaction not found")
	}
	if int(amount) > merchant.Transactions[idx].Authorized {
		return errors.New("can not over-capture")
	}
	card, err := db.GetCard(merchant.Transactions[idx].CardNumber)
	if err != nil {
		return err
	}

	merchant.Transactions[idx].Authorized = merchant.Transactions[idx].Authorized - int(amount)
	merchant.Transactions[idx].Captured = merchant.Transactions[idx].Captured + int(amount)
	card.BlockedBalance = card.BlockedBalance - int(amount)
	err = db.StoreCard(card)
	if err != nil {
		return err
	}
	return db.StoreMerchant(merchant)
}

func findTransactionIndex(transactions []model.Transaction, transactionID string) int {
	for idx, transaction := range transactions {
		if transaction.ID == transactionID {
			return idx
		}
	}
	return -1
}
