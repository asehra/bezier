package service

import (
	"errors"

	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func AuthorizeTransaction(db storage.Storage, cardNumber int64, merchantID string, amount int32, idGenerator generator.StringIDGenerator) (string, error) {
	card, err := db.GetCard(cardNumber)
	if err != nil {
		return "", err
	}
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
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return "", err
	}
	merchant.AuthorizedTransactions = append(merchant.AuthorizedTransactions, transaction)
	{ // NOTE This should be transactional on a real system
		db.StoreCard(card)         //TODO: handle error
		db.StoreMerchant(merchant) //TODO: handle herror
	}
	return transaction.ID, nil
}
