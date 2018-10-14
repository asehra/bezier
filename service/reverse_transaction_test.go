package service_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReverseTransaction(t *testing.T) {

	Convey("Given a merchant is attempting to reverse a transaction", t, func() {
		mockMerchantId := "M100"
		db := &mock.Storage{}
		mockTransactionID := "transactionID"
		cardNumber := int64(123456789012345)
		db.GetMerchantCall.Returns.Merchant = model.Merchant{
			ID: mockMerchantId,
			Transactions: []model.Transaction{{
				ID:         mockTransactionID,
				CardNumber: cardNumber,
				Authorized: 100,
			}},
		}

		Convey("With full amount", func() {

			service.ReverseTransaction(db, mockMerchantId, mockTransactionID, 100)
			Convey("Full amount is moved to Reversed", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Authorized: 0,
					Reversed:   100,
				}})
			})
		})

		Convey("With partial amount", func() {
			service.ReverseTransaction(db, mockMerchantId, mockTransactionID, 60)
			Convey("Partial amount is moved to Reversed", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Authorized: 40,
					Reversed:   60,
				}})
			})
		})
		Convey("With more than authorized amount", func() {
			err := service.ReverseTransaction(db, mockMerchantId, mockTransactionID, 110)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("can not over-reverse"))
			})

		})

		Convey("With invalid merchant ID", func() {
			db.GetMerchantCall.Returns.Error = errors.New("not found")
			err := service.ReverseTransaction(db, "invalid ID", mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("not found"))
			})
		})

		Convey("And the transaction does not exist", func() {
			err := service.ReverseTransaction(db, mockMerchantId, "bad transaction", 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("transaction not found"))
			})
		})

		Convey("And storing merchant raises error", func() {
			db.StoreMerchantCall.Returns.Error = errors.New("something went wrong")
			err := service.ReverseTransaction(db, "invalid ID", mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})
	})
}
