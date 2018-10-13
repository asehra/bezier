package service

import (
	"errors"

	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func AuthorizeTransaction(db storage.Storage, cardNumber int64, merchantID string, amount int32, idGenerator generator.StringIDGenerator) (string, error) {
	card, _ := db.GetCard(cardNumber) //TODO: handle error
	if amount > card.AvailableBalance {
		return "", errors.New("insufficient funds")
	}
	card.AvailableBalance = card.AvailableBalance - amount
	card.BlockedBalance = card.BlockedBalance + amount
	transaction := model.Transaction{
		ID:         idGenerator.Generate(),
		CardNumber: cardNumber,
		Amount:     amount,
	}
	merchant, _ := db.GetMerchant(merchantID) //TODO: handle error
	merchant.AuthorizedTransactions = append(merchant.AuthorizedTransactions, transaction)
	db.StoreCard(card)         //TODO: handle error
	db.StoreMerchant(merchant) //TODO: handle herror
	return transaction.ID, nil
}
