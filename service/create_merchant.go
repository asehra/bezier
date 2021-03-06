package service

import (
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func CreateMerchant(db storage.Storage, idGenerator generator.StringIDGenerator) (string, error) {
	merchantID := idGenerator.Generate()
	err := db.StoreMerchant(model.Merchant{
		ID:           merchantID,
		Transactions: []model.Transaction{},
	})
	return merchantID, err
}
