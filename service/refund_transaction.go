package service

import (
	"errors"

	"github.com/asehra/bezier/storage"
)

func RefundTransaction(db storage.Storage, merchantID string, transactionID string, amount uint) error {
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return err
	}
	idx := findTransactionIndex(merchant.Transactions, transactionID)
	if idx == -1 {
		return errors.New("transaction not found")
	}
	if int(amount) > merchant.Transactions[idx].Captured {
		return errors.New("can not over-refund")
	}
	card, err := db.GetCard(merchant.Transactions[idx].CardNumber)
	if err != nil {
		return err
	}
	merchant.Transactions[idx].Captured = merchant.Transactions[idx].Captured - int(amount)
	merchant.Transactions[idx].Refunded = merchant.Transactions[idx].Refunded + int(amount)
	card.AvailableBalance = card.AvailableBalance + int(amount)
	{ // This should be done as a transaction on a real system
		err = db.StoreMerchant(merchant)
		if err != nil {
			return err
		}
		err = db.StoreCard(card)
		if err != nil {
			return err
		}
	}
	return nil
}
