package webserver_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asehra/bezier/model"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/storage"
	"github.com/asehra/bezier/webserver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMerchantHandlersAPI(t *testing.T) {
	Convey("Merchant Actions", t, func() {
		mockMerchantID := "M012345"
		db := storage.NewInMemoryStore()
		testConfig := config.Config{
			DB: db,
			// CardIDGenerator:     &mock.CardIDGenerator{Generates: mockCardNumber},
			MerchantIDGenerator: &mock.StringIDGenerator{MockID: mockMerchantID},
			StdErr:              ioutil.Discard,
			StdOut:              ioutil.Discard,
		}

		Convey("GET /v1/merchant/create", func() {
			req, _ := http.NewRequest("GET", "/v1/merchant/create", nil)
			w := httptest.NewRecorder()

			api := webserver.Create(testConfig)
			api.ServeHTTP(w, req)

			Convey("Returns 200 status code", func() {
				So(w.Code, ShouldEqual, 200)
			})

			Convey("Returns merchant ID in body", func() {
				bodyAsString := w.Body.String()
				So(bodyAsString, ShouldEqual, fmt.Sprintf(`{"merchant_id":"%s","error":null}`, mockMerchantID))
			})

			Convey("Storage should have a merchant with the returned merchant_id", func() {
				var response webserver.CreateMerchantResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)
				merchant, err := db.GetMerchant(response.MerchantID)
				So(err, ShouldBeNil)
				So(merchant.ID, ShouldEqual, mockMerchantID)
			})
		})

		Convey("POST /v1/merchant/transaction-authorization-request", func() {
			mockCardNumber := int64(9000000000000001)
			db.StoreMerchant(model.Merchant{mockMerchantID, []model.Transaction{}})
			db.StoreCard(model.Card{
				Number:           mockCardNumber,
				AvailableBalance: 1000,
				BlockedBalance:   100,
			})
			mockTransactionID := "TX88888"
			testConfig.TransactionIDGenerator = &mock.StringIDGenerator{MockID: mockTransactionID}

			Convey("When card has sufficient funds", func() {
				transactionAmount := int32(50)

				requestBody := authRequestBody(mockCardNumber, mockMerchantID, transactionAmount)
				req, _ := http.NewRequest("POST", "/v1/merchant/authorize-transaction", requestBody)
				w := httptest.NewRecorder()

				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 200 status code", func() {
					So(w.Code, ShouldEqual, 200)
				})

				Convey("Adds transaction to Merchant's Authorized Transactions List", func() {
					merchant, err := db.GetMerchant(mockMerchantID)
					So(err, ShouldBeNil)
					So(merchant.AuthorizedTransactions, ShouldResemble, []model.Transaction{{mockTransactionID, mockCardNumber, transactionAmount}})
				})

				Convey("Returns transaction ID in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, fmt.Sprintf(`"transaction_id":"%s"`, mockTransactionID))
				})

				Convey("Moves funds from AvailableBalance to BlockedBalance on card", func() {
					card, err := db.GetCard(mockCardNumber)
					So(err, ShouldBeNil)
					So(card.AvailableBalance, ShouldEqual, 950)
					So(card.BlockedBalance, ShouldEqual, 150)
				})
			})

			Convey("When card has insufficient funds", func() {
				transactionAmount := int32(5000)
				requestBody := authRequestBody(mockCardNumber, mockMerchantID, transactionAmount)
				req, _ := http.NewRequest("POST", "/v1/merchant/authorize-transaction", requestBody)
				w := httptest.NewRecorder()

				api := webserver.Create(testConfig)
				api.ServeHTTP(w, req)

				Convey("Returns 400", func() {
					So(w.Code, ShouldEqual, 400)
				})

				Convey("Returns appropriate error in body", func() {
					bodyAsString := w.Body.String()
					So(bodyAsString, ShouldContainSubstring, `"error":"insufficient funds"`)
				})
			})
		})
	})
}

func authRequestBody(cardNumber int64, merchantID string, transactionAmount int32) io.Reader {
	return strings.NewReader(fmt.Sprintf(
		`{
			"card_number": %d,
			"merchant_id": "%s",
			"amount": %d
		}`,
		cardNumber,
		merchantID,
		transactionAmount,
	))
}
