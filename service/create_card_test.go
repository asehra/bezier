package service

import (
	"testing"

	"github.com/asehra/bezier/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateCard(t *testing.T) {

	Convey("Given A card creation system", t, func() {

		Convey("When a card is created", func() {
			expectedCardNumber := int64(9000000000000001)
			idGenerator := &mock.IDGenerator{Generates: expectedCardNumber}
			storage := &mock.Storage{}
			actualCardNumber := CreateCard(storage, idGenerator)

			Convey("Returns the card number of the new card as Generated by IDGenerator", func() {
				So(actualCardNumber, ShouldEqual, expectedCardNumber)
			})

			Convey("Stores the new card on the storage with a generated ID", func() {
				So(storage.StoredCard.Number, ShouldEqual, expectedCardNumber)
			})

			Convey("The new card has 0 available balance", func() {
				So(storage.StoredCard.AvailableBalance, ShouldEqual, 0)
			})

			// TODO: Create new card with unique numbers
			Convey("When another card is created", nil)
		})
	})
}