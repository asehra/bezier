package service

import (
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

type IDGenerator interface {
	Generate() int64
}

func CreateCard(db storage.Storage, idGenerator IDGenerator) (int64, error) {
	cardNumber := idGenerator.Generate()
	err := db.StoreCard(model.Card{
		Number: cardNumber,
	})
	return cardNumber, err
}
