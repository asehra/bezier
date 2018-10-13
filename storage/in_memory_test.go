package storage_test

import (
	"errors"
	"testing"

	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemoryStorage(t *testing.T) {
	db := storage.NewInMemoryStore()
	cardNumber := int64(9)
	card := model.Card{Number: cardNumber}

	Convey("StoreCard", t, func() {
		Convey("Updates in memory Card store", func() {
			err := db.StoreCard(card)
			So(db.CardStore[cardNumber], ShouldResemble, card)
			So(err, ShouldBeNil)
		})
	})

	Convey("GetCard", t, func() {
		Convey("Fetches card when exists", func() {
			db.CardStore[cardNumber] = card
			card, err := db.GetCard(cardNumber)
			So(err, ShouldBeNil)
			So(card, ShouldResemble, card)
		})
		Convey("Returns error when it doesn't exist", func() {
			_, err := db.GetCard(int64(10))
			So(err, ShouldResemble, errors.New("Card not found"))
		})
	})
}
