package service

import (
	"errors"

	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func AuthorizeTransaction(db storage.Storage, cardNumber int64, merchantID string, amount uint, idGenerator generator.StringIDGenerator) (string, error) {
	card, err := db.GetCard(cardNumber)
	if err != nil {
		return "", err
	}
	if int(amount) > card.AvailableBalance {
		return "", errors.New("insufficient funds")
	}
	card.AvailableBalance = card.AvailableBalance - int(amount)
	card.BlockedBalance = card.BlockedBalance + int(amount)
	transaction := model.Transaction{
		ID:         idGenerator.Generate(),
		CardNumber: cardNumber,
		Authorized: int(amount),
	}
	merchant, err := db.GetMerchant(merchantID)
	if err != nil {
		return "", err
	}
	merchant.Transactions = append(merchant.Transactions, transaction)
	{ // NOTE This should be transactional on a real system
		db.StoreCard(card)         //TODO: handle error
		db.StoreMerchant(merchant) //TODO: handle herror
	}
	return transaction.ID, nil
}
