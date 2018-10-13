package service

import (
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func CreateCard(db storage.Storage, idGenerator generator.IDGenerator) (int64, error) {
	cardNumber := idGenerator.Generate()
	err := db.StoreCard(model.Card{
		Number: cardNumber,
	})
	return cardNumber, err
}
