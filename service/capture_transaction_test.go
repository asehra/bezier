package service_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCaptureTransaction(t *testing.T) {

	Convey("Given a merchant is attempting to capture a transaction", t, func() {
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
		db.GetCardCall.Returns.Card = model.Card{
			Number:           cardNumber,
			AvailableBalance: 200,
			BlockedBalance:   150,
		}

		Convey("With full amount", func() {
			service.CaptureTransaction(db, mockMerchantId, mockTransactionID, 100)
			Convey("Full amount is moved to captured", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Authorized: 0,
					Captured:   100,
				}})
			})
			Convey("Full amount is deducted from the Card's blocked balance", func() {
				So(db.StoreCardCall.Receives.Card.BlockedBalance, ShouldEqual, 50)
			})
		})

		Convey("With partial amount", func() {
			service.CaptureTransaction(db, mockMerchantId, mockTransactionID, 60)
			Convey("Partial amount is moved to captured", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Authorized: 40,
					Captured:   60,
				}})
			})
		})
		Convey("With more than authorized amount", func() {
			err := service.CaptureTransaction(db, mockMerchantId, mockTransactionID, 110)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("can not over-capture"))
			})

		})

		Convey("With invalid merchant ID", func() {
			db.GetMerchantCall.Returns.Error = errors.New("not found")
			err := service.CaptureTransaction(db, "invalid ID", mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("not found"))
			})
		})

		Convey("And the transaction does not exist", func() {
			err := service.CaptureTransaction(db, mockMerchantId, "bad transaction", 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("transaction not found"))
			})
		})

		Convey("And storing merchant raises error", func() {
			db.StoreMerchantCall.Returns.Error = errors.New("something went wrong")
			err := service.CaptureTransaction(db, "invalid ID", mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})

		Convey("And fetching card raises error", func() {
			db.GetCardCall.Returns.Error = errors.New("something went wrong")
			err := service.CaptureTransaction(db, mockMerchantId, mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})

		Convey("And storing card raises error", func() {
			db.StoreCardCall.Returns.Error = errors.New("something went wrong")
			err := service.CaptureTransaction(db, mockMerchantId, mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})
	})
}
