package service_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/service"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthorizeTransaction(t *testing.T) {

	Convey("Given a merchant is attempting to authorize a transaction", t, func() {
		mockTransactionID := "transactionID"
		cardNumber := int64(123456789012345)
		merchantID := "M12345"
		amount := 50
		idGenerator := &mock.StringIDGenerator{MockID: mockTransactionID}
		db := &mock.Storage{}
		Convey("On a valid card", func() {
			db.GetCardCall.Returns.Card = model.Card{
				Number:           cardNumber,
				AvailableBalance: 1000,
			}

			Convey("With Sufficient funds", func() {
				transactionID, err := service.AuthorizeTransaction(db, cardNumber, merchantID, uint(amount), idGenerator)
				Convey("Transaction ID generted by the idGenerator is returned", func() {
					So(transactionID, ShouldEqual, mockTransactionID)
					So(err, ShouldBeNil)
				})
			})
			Convey("With Insufficient funds", func() {
				amount := 1001
				transactionID, err := service.AuthorizeTransaction(db, cardNumber, merchantID, uint(amount), idGenerator)
				Convey("Transaction ID generted by the idGenerator is returned", func() {
					So(transactionID, ShouldEqual, "")
					So(err, ShouldResemble, errors.New("insufficient funds"))
				})
			})
		})

		Convey("On an invalid card", func() {
			expectedError := errors.New("Bad card")
			db.GetCardCall.Returns.Error = expectedError
			transactionID, err := service.AuthorizeTransaction(db, cardNumber, merchantID, uint(amount), idGenerator)
			Convey("Transaction is declined with error", func() {
				So(transactionID, ShouldEqual, "")
				So(err, ShouldResemble, expectedError)
			})
		})

		Convey("With an invalid merchant ID", func() {
			db.GetCardCall.Returns.Card = model.Card{
				Number:           cardNumber,
				AvailableBalance: 1000,
			}
			expectedError := errors.New("Bad merchant")
			db.GetMerchantCall.Returns.Error = expectedError
			transactionID, err := service.AuthorizeTransaction(db, cardNumber, merchantID, uint(amount), idGenerator)
			Convey("Transaction is declined with error", func() {
				So(transactionID, ShouldEqual, "")
				So(err, ShouldResemble, expectedError)
			})
		})
	})
}
