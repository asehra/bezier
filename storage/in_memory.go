package storage

import (
	"errors"

	"github.com/asehra/bezier/model"
)

type InMemory struct {
	CardStore     map[int64]model.Card
	MerchantStore map[string]model.Merchant
}

func NewInMemoryStore() *InMemory {
	return &InMemory{
		CardStore:     map[int64]model.Card{},
		MerchantStore: map[string]model.Merchant{},
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

func (i *InMemory) StoreMerchant(merchant model.Merchant) error {
	i.MerchantStore[merchant.ID] = merchant
	return nil
}

func (i *InMemory) GetMerchant(merchantID string) (model.Merchant, error) {
	merchant, ok := i.MerchantStore[merchantID]
	if !ok {
		return model.Merchant{}, errors.New("Merchant not found")
	}
	return merchant, nil
}
