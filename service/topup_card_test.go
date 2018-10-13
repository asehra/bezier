package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTopUpCard(t *testing.T) {

	Convey("Given A card number", t, func() {

		Convey("When it is topped up", func() {
			cardNumber := int64(9000000000000001)
			db := &mock.Storage{}

			Convey("It's available balance is incremented by the amount", func() {
				testCases := []struct {
					Description     string
					Amount          int32
					InitialBalance  int32
					ExpectedBalance int32
				}{
					{"zero balance", 5000, 0, 5000},
					{"non-zero", 5000, 100, 5100},
				}

				for _, testCase := range testCases {
					Convey(fmt.Sprintf("Works for: %s", testCase.Description), func() {
						db.GetCardCall.Returns.Card = model.Card{
							Number:           cardNumber,
							AvailableBalance: testCase.InitialBalance,
						}

						err := service.TopUpCard(db, cardNumber, testCase.Amount)
						So(err, ShouldBeNil)
						So(db.StoreCardCall.Receives.Card.AvailableBalance, ShouldEqual, testCase.ExpectedBalance)
					})
				}
			})

			Convey("When card lookup returns in error", func() {
				Convey("Balance is not updated", func() {
					db.GetCardCall.Returns.Error = errors.New("Unable to find card")

					err := service.TopUpCard(db, cardNumber, 100)
					So(err, ShouldNotBeNil)
				})
			})
			Convey("When card update returns in error", func() {
				Convey("Error is propogated", func() {
					db.StoreCardCall.Returns.Error = errors.New("Unable to find card")

					err := service.TopUpCard(db, cardNumber, 100)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}