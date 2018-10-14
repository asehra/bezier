package service

import (
	"errors"

	"github.com/asehra/bezier/storage"
)

func ReverseTransaction(db storage.Storage, merchantID string, transactionID string, amount uint) error {
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return err
	}
	idx := findTransactionIndex(merchant.Transactions, transactionID)
	if idx == -1 {
		return errors.New("transaction not found")
	}
	if int(amount) > merchant.Transactions[idx].Authorized {
		return errors.New("can not over-reverse")
	}
	merchant.Transactions[idx].Authorized = merchant.Transactions[idx].Authorized - int(amount) // TODO handle over-capture
	merchant.Transactions[idx].Reversed = merchant.Transactions[idx].Reversed + int(amount)
	return db.StoreMerchant(merchant)
}
