package service

import "github.com/asehra/bezier/model"

type Storage interface {
	StoreCard(model.Card) error
}

type IDGenerator interface {
	Generate() int64
}

func CreateCard(storage Storage, idGenerator IDGenerator) (int64, error) {
	cardNumber := idGenerator.Generate()
	err := storage.StoreCard(model.Card{
		Number: cardNumber,
	})
	return cardNumber, err
}
