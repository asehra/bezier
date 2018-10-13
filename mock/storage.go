package mock

import "github.com/asehra/bezier/model"

type Storage struct {
	StoredCard   model.Card
	ReturnsError error
}

func (s *Storage) StoreCard(card model.Card) error {
	s.StoredCard = card
	return s.ReturnsError
}
