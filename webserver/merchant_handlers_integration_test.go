package webserver_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		mockCardNumber := int64(9000000000000001)
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
			output := simulateGet(testConfig, "/v1/merchant/create")

			Convey("Returns 200 status code", func() {
				So(output.Code, ShouldEqual, 200)
			})

			Convey("Returns merchant ID in body", func() {
				bodyAsString := output.Body.String()
				So(bodyAsString, ShouldEqual, fmt.Sprintf(`{"merchant_id":"%s","error":null}`, mockMerchantID))
			})

			Convey("Storage should have a merchant with the returned merchant_id", func() {
				var response webserver.CreateMerchantResponse
				err := json.Unmarshal(output.Body.Bytes(), &response)
				So(err, ShouldBeNil)
				merchant, err := db.GetMerchant(response.MerchantID)
				So(err, ShouldBeNil)
				So(merchant.ID, ShouldEqual, mockMerchantID)
			})
		})

		Convey("GET /v1/merchant/transactions", func() {
			merchant := model.Merchant{
				ID: mockMerchantID,
				AuthorizedTransactions: []model.Transaction{
					{
						ID:         "TX101",
						CardNumber: mockCardNumber,
						Amount:     100,
					},
				},
			}
			err := db.StoreMerchant(merchant)
			So(err, ShouldBeNil)

			Convey("with a valid merchant ID", func() {
				path := fmt.Sprintf("/v1/merchant/transactions?merchant_id=%s", mockMerchantID)
				output := simulateGet(testConfig, path)

				Convey("Returns 200 status code", func() {
					So(output.Code, ShouldEqual, 200)
				})

				Convey("Returns the transaction activity by a merchant", func() {
					var response webserver.MerchantTransactionsResponse
					err := json.Unmarshal(output.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response.Merchant, ShouldResemble, merchant)
				})
			})

			Convey("with invalid merchant ID", func() {
				path := fmt.Sprintf(`/v1/merchant/transactions?merchant_id=%d`, int64(0))
				output := simulateGet(testConfig, path)

				Convey("Returns 400 status code", func() {
					So(output.Code, ShouldEqual, 400)
				})

				Convey("Returns error in body", func() {
					var response webserver.MerchantTransactionsResponse

					err := json.Unmarshal(output.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response.Error, ShouldEqual, "Merchant not found")
				})
			})

			Convey("with no merchant ID", func() {
				output := simulateGet(testConfig, "/v1/merchant/transactions?merchant_id=")

				Convey("Returns 400 status code", func() {
					So(output.Code, ShouldEqual, 400)
				})

				Convey("Returns error in body", func() {
					var response webserver.MerchantTransactionsResponse

					err := json.Unmarshal(output.Body.Bytes(), &response)
					So(err, ShouldBeNil)
					So(response.Error, ShouldEqual, "Bad merchant ID format")
				})
			})
		})

		Convey("POST /v1/merchant/transaction-authorization-request", func() {
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
				output := simulatePost(testConfig, "/v1/merchant/authorize-transaction", requestBody)

				Convey("Returns 200 status code", func() {
					So(output.Code, ShouldEqual, 200)
				})

				Convey("Adds transaction to Merchant's Authorized Transactions List", func() {
					merchant, err := db.GetMerchant(mockMerchantID)
					So(err, ShouldBeNil)
					So(merchant.AuthorizedTransactions, ShouldResemble, []model.Transaction{{mockTransactionID, mockCardNumber, transactionAmount}})
				})

				Convey("Returns transaction ID in body", func() {
					bodyAsString := output.Body.String()
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
				largeAmount := int32(5000)
				requestBody := authRequestBody(mockCardNumber, mockMerchantID, largeAmount)
				output := simulatePost(testConfig, "/v1/merchant/authorize-transaction", requestBody)

				Convey("Returns 400", func() {
					So(output.Code, ShouldEqual, 400)
				})

				Convey("Returns appropriate error in body", func() {
					bodyAsString := output.Body.String()
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
