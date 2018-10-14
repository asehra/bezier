package service_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRefundTransaction(t *testing.T) {

	Convey("Given a merchant is attempting to refund a transaction", t, func() {
		mockMerchantId := "M100"
		db := &mock.Storage{}
		mockTransactionID := "transactionID"
		cardNumber := int64(123456789012345)
		db.GetCardCall.Returns.Card = model.Card{
			Number:           cardNumber,
			AvailableBalance: 200,
		}
		db.GetMerchantCall.Returns.Merchant = model.Merchant{
			ID: mockMerchantId,
			Transactions: []model.Transaction{{
				ID:         mockTransactionID,
				CardNumber: cardNumber,
				Captured:   100,
			}},
		}

		Convey("With full amount", func() {
			service.RefundTransaction(db, mockMerchantId, mockTransactionID, 100)
			Convey("Full amount is moved to Refunded", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Captured:   0,
					Refunded:   100,
				}})
				So(db.StoreCardCall.Receives.Card.AvailableBalance, ShouldEqual, 300)
			})
		})

		Convey("With partial amount", func() {
			service.RefundTransaction(db, mockMerchantId, mockTransactionID, 60)
			Convey("Partial amount is moved to Refunded", func() {
				So(db.StoreMerchantCall.Receives.Merchant.Transactions, ShouldResemble, []model.Transaction{{
					ID:         mockTransactionID,
					CardNumber: cardNumber,
					Captured:   40,
					Refunded:   60,
				}})

				So(db.StoreCardCall.Receives.Card.AvailableBalance, ShouldEqual, 260)
			})
		})

		Convey("With more than Captured amount", func() {
			err := service.RefundTransaction(db, mockMerchantId, mockTransactionID, 110)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("can not over-refund"))
			})

		})

		Convey("With invalid merchant ID", func() {
			db.GetMerchantCall.Returns.Error = errors.New("not found")
			err := service.RefundTransaction(db, "invalid ID", mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("not found"))
			})
		})

		Convey("And card details fail to retrieve from database", func() {
			db.GetCardCall.Returns.Error = errors.New("can't fetch card")
			err := service.RefundTransaction(db, mockMerchantId, mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("can't fetch card"))
			})
		})

		Convey("And the transaction does not exist", func() {
			err := service.RefundTransaction(db, mockMerchantId, "bad transaction", 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("transaction not found"))
			})
		})

		Convey("And storing merchant raises error", func() {
			db.StoreMerchantCall.Returns.Error = errors.New("something went wrong")
			err := service.RefundTransaction(db, mockMerchantId, mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})
		Convey("And storing card raises error", func() {
			db.StoreCardCall.Returns.Error = errors.New("something went wrong")
			err := service.RefundTransaction(db, mockMerchantId, mockTransactionID, 50)
			Convey("Error is raised", func() {
				So(err, ShouldResemble, errors.New("something went wrong"))
			})
		})
	})
}
