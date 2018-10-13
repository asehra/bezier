package mock

import "github.com/asehra/bezier/model"

type Storage struct {
	StoreCardCall struct {
		Receives struct {
			Card model.Card
		}
		Returns struct {
			Error error
		}
	}
	GetCardCall struct {
		Receives struct {
			CardNumber int64
		}
		Returns struct {
			Card  model.Card
			Error error
		}
	}
	StoreMerchantCall struct {
		Receives struct {
			Merchant model.Merchant
		}
		Returns struct {
			Error error
		}
	}
	GetMerchantCall struct {
		Receives struct {
			MerchantID string
		}
		Returns struct {
			Merchant model.Merchant
			Error    error
		}
	}
}

func (s *Storage) StoreCard(card model.Card) error {
	s.StoreCardCall.Receives.Card = card
	return s.StoreCardCall.Returns.Error
}

func (s *Storage) GetCard(cardNumber int64) (model.Card, error) {
	return s.GetCardCall.Returns.Card, s.GetCardCall.Returns.Error
}

func (s *Storage) StoreMerchant(merchant model.Merchant) error {
	s.StoreMerchantCall.Receives.Merchant = merchant
	return s.StoreMerchantCall.Returns.Error
}

func (s *Storage) GetMerchant(MerchantID string) (model.Merchant, error) {
	return s.GetMerchantCall.Returns.Merchant, s.GetMerchantCall.Returns.Error
}
