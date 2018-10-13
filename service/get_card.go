package service

import (
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func GetCard(db storage.Storage, cardNumber int64) (model.Card, error) {
	return db.GetCard(cardNumber)
}
