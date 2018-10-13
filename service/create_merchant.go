package service

import (
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func CreateMerchant(db storage.Storage, idGenerator generator.MerchantIDGenerator) (string, error) {
	merchantID := idGenerator.Generate()
	err := db.StoreMerchant(model.Merchant{
		ID: merchantID,
	})
	return merchantID, err
}
