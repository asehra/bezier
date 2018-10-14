package service

import "github.com/asehra/bezier/storage"

func TopUpCard(db storage.Storage, cardNumber int64, amount uint) error {
	card, err := db.GetCard(cardNumber)
	if err != nil {
		return err
	}
	card.AvailableBalance = card.AvailableBalance + int(amount)
	card.TotalLoaded = card.TotalLoaded + int(amount)
	err = db.StoreCard(card)
	if err != nil {
		return err
	}
	return nil
}
