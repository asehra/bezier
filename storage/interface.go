package storage

import "github.com/asehra/bezier/model"

type Storage interface {
	StoreCard(card model.Card) error
	GetCard(cardNumber int64) (model.Card, error)
	StoreMerchant(merchant model.Merchant) error
	GetMerchant(merchantID string) (model.Merchant, error)
}
