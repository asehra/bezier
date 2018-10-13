package service

import (
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func AuthorizeTransaction(db storage.Storage, cardNumber int64, merchantID string, amount int32, idGenerator generator.StringIDGenerator) string {
	card, _ := db.GetCard(cardNumber) //TODO: handle error
	card.AvailableBalance = card.AvailableBalance - 50
	card.BlockedBalance = card.BlockedBalance + 50
	transaction := model.Transaction{
		ID:         idGenerator.Generate(),
		CardNumber: cardNumber,
		Amount:     amount,
	}
	merchant, _ := db.GetMerchant(merchantID) //TODO: handle error
	merchant.AuthorizedTransactions = append(merchant.AuthorizedTransactions, transaction)
	db.StoreCard(card)         //TODO: handle error
	db.StoreMerchant(merchant) //TODO: handle herror
	return idGenerator.Generate()
}
