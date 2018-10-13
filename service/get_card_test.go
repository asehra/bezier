package service_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/service"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCardStatus(t *testing.T) {

	Convey("Get Card Status", t, func() {
		cardNumber := int64(9000000000000001)
		db := &mock.Storage{}

		Convey("When card is found in storage", func() {
			db.GetCardCall.Returns.Card = model.Card{
				Number:           cardNumber,
				AvailableBalance: 1000,
			}

			Convey("Retrieves the card from storage", func() {
				card, _ := service.GetCard(db, cardNumber)
				So(card.Number, ShouldEqual, cardNumber)
			})
		})

		Convey("When storage returns error", func() {
			Convey("GetCard service returns error", func() {
				db.GetCardCall.Returns.Error = errors.New("Unable to find card")
				_, err := service.GetCard(db, cardNumber)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
