package bezier

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateCard(t *testing.T) {

	Convey("Given A card creation system", t, func() {

		Convey("When a card is created", func() {
			card := CreateCard()

			Convey("The card has a number", func() {
				So(card.Number, ShouldEqual, 9000000000000001)
			})

			Convey("The card has an available balance", func() {
				So(card.AvailableBalance, ShouldEqual, 0)
			})

			// TODO: Create new card with unique numbers
			Convey("When another card is created", nil)
		})
	})
}
