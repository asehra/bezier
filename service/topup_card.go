package service

import "github.com/asehra/bezier/storage"

func TopUpCard(db storage.Storage, cardNumber int64, amount int32) error {
	card, err := db.GetCard(cardNumber)
	if err != nil {
		return err
	}
	card.AvailableBalance = card.AvailableBalance + amount
	card.TotalLoaded = card.TotalLoaded + amount
	err = db.StoreCard(card)
	if err != nil {
		return err
	}
	return nil
}
