package service

import (
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"
)

func GetMerchant(db storage.Storage, id string) (model.Merchant, error) {
	return db.GetMerchant(id)
}
