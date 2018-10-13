package webserver_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/model"
	"github.com/asehra/bezier/storage"

	"github.com/asehra/bezier/webserver"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCardHandlersAPI(t *testing.T) {
	Convey("User Actions", t, func() {

		mockCardNumber := int64(9000000000000001)
		db := storage.NewInMemoryStore()
		testConfig := config.Config{
			DB:              db,
			CardIDGenerator: &mock.CardIDGenerator{Generates: mockCardNumber},
			StdErr:          ioutil.Discard,
			StdOut:          ioutil.Discard,
		}

		Convey("GET /v1/card/create", func() {
			req, _ := http.NewRequest("GET", "/v1/card/create", nil)
			w := httptest.NewRecorder()

			api := webserver.Create(testConfig)
			api.ServeHTTP(w, req)

			Convey("Returns 200 status code", func() {
				So(w.Code, ShouldEqual, 200)
			})

			Convey("Returns card number in body", func() {
				bodyAsString := w.Body.String()
				So(bodyAsString, ShouldEqual, fmt.Sprintf(`{"card_number":%d,"error":null}`, mockCardNumber))
			})

			Convey("Storage should have a card with the returned card_number", func() {
				var response webserver.CreateCardResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)
				So(db.CardStore[response.CardNumber].Number, ShouldEqual, mockCardNumber)
			})
		})

		Convey("GET /v1/card/details", func() {
			Convey("With valid card number", func() {
				card := model.Card{Number: mockCardNumber, AvailableBalance: 1000, BlockedBalance: 500, TotalLoaded: 2000}
				path := fmt.Sprintf(`/v1/card/details?card_number=%d`, card.Number)
				req, _ := http.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()
				testConfig.DB.StoreCard(card)
				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 200 status code", func() {
					So(w.Code, ShouldEqual, 200)
				})

				Convey("Returns card details in body", func() {
					bodyAsString := w.Body.String()
					expectedBody := fmt.Sprintf(
						`"card_number":%d,"available_balance":%d,"blocked_balance":%d,"total_loaded":%d`,
						card.Number,
						card.AvailableBalance,
						card.BlockedBalance,
						card.TotalLoaded,
					)
					So(bodyAsString, ShouldContainSubstring, expectedBody)
				})
			})

			Convey("with invalid card number", func() {
				path := fmt.Sprintf(`/v1/card/details?card_number=%d`, int64(0))
				req, _ := http.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()
				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 400 status code", func() {
					So(w.Code, ShouldEqual, 400)
				})

				Convey("Returns error in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, `"error":"Card not found"`)
				})
			})

			Convey("with no card number", func() {
				req, _ := http.NewRequest("GET", "/v1/card/details?card_number=", nil)
				w := httptest.NewRecorder()
				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 400 status code", func() {
					So(w.Code, ShouldEqual, 400)
				})

				Convey("Returns error in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, `"error":"Bad card number format"`)
				})
			})
		})

		Convey("POST /v1/card/top-up", func() {
			Convey("With valid json request body", func() {
				card := model.Card{Number: mockCardNumber, AvailableBalance: 500}
				db.StoreCard(card)
				data := fmt.Sprintf(`{"card_number":%d,"amount":%d}`, card.Number, 1000)
				req, _ := http.NewRequest("POST", "/v1/card/top-up", strings.NewReader(data))
				w := httptest.NewRecorder()

				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 200 status code", func() {
					So(w.Code, ShouldEqual, 200)
				})

				Convey("Updates balance for the card in db", func() {
					updatedCard, _ := db.GetCard(card.Number)
					So(updatedCard.AvailableBalance, ShouldEqual, 1500)
				})
			})

			Convey("With invalid json request body", func() {
				req, _ := http.NewRequest("POST", "/v1/card/top-up", strings.NewReader("bad json"))
				w := httptest.NewRecorder()

				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 400 status code", func() {
					So(w.Code, ShouldEqual, 400)
				})
				Convey("Returns error in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, `"error":"bad JSON format"`)
				})
			})

			Convey("With bad card details in body", func() {
				data := fmt.Sprintf(`{"card_number":%d,"amount":%d}`, int64(10), 1000)
				req, _ := http.NewRequest("POST", "/v1/card/top-up", strings.NewReader(data))
				w := httptest.NewRecorder()

				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 400 status code", func() {
					So(w.Code, ShouldEqual, 400)
				})

				Convey("Returns error in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, `"error":"Card not found"`)
				})
			})
		})
	})
}