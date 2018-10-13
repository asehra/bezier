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
}

func (s *Storage) StoreCard(card model.Card) error {
	s.StoreCardCall.Receives.Card = card
	return s.StoreCardCall.Returns.Error
}

func (s *Storage) GetCard(cardNumber int64) (model.Card, error) {
	return s.GetCardCall.Returns.Card, s.GetCardCall.Returns.Error
}
