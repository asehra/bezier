package storage

import (
	"errors"

	"github.com/asehra/bezier/model"
)

type InMemory struct {
	CardStore map[int64]model.Card
}

func NewInMemoryStore() *InMemory {
	return &InMemory{
		CardStore: map[int64]model.Card{},
	}
}

func (i *InMemory) StoreCard(card model.Card) error {
	i.CardStore[card.Number] = card
	return nil
}

func (i *InMemory) GetCard(cardNumber int64) (model.Card, error) {
	card, ok := i.CardStore[cardNumber]
	if !ok {
		return model.Card{}, errors.New("Card not found")
	}
	return card, nil
}
