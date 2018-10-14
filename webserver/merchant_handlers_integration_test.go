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
				Transactions: []model.Transaction{
					{
						ID:         "TX101",
						CardNumber: mockCardNumber,
						Authorized: 100,
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

		Convey("POST /v1/merchant/authorize-transaction", func() {
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
					So(merchant.Transactions, ShouldResemble, []model.Transaction{{
						ID:         mockTransactionID,
						CardNumber: mockCardNumber,
						Authorized: transactionAmount,
						Captured:   0,
					}})
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

			Convey("When request is badly formed", func() {
				output := simulatePost(testConfig, "/v1/merchant/authorize-transaction", strings.NewReader("bad data"))
				Convey("Returns 400", func() {
					So(output.Code, ShouldEqual, 400)
				})

				Convey("Returns bad JSON response", func() {
					bodyAsString := output.Body.String()
					So(bodyAsString, ShouldEqual, `{"error":"bad JSON format"}`)
				})
			})
		})

		Convey("POST /v1/merchant/capture-transaction", func() {
			mockTransactionID := "TX88888"
			authorizedTransaction := model.Transaction{
				ID:         mockTransactionID,
				CardNumber: mockCardNumber,
				Authorized: 100,
			}
			db.StoreMerchant(model.Merchant{
				ID:           mockMerchantID,
				Transactions: []model.Transaction{authorizedTransaction},
			})

			Convey("When capture is possible", func() {
				requestBody := captureRequestBody(mockMerchantID, mockTransactionID, 170)
				output := simulatePost(testConfig, "/v1/merchant/capture-transaction", requestBody)
				Convey("Returns 400", func() {
					So(output.Code, ShouldEqual, 400)
				})
				
				Convey("Returns error message", func() {
					bodyAsString := output.Body.String()
					So(bodyAsString, ShouldEqual, `{"error":"can not over-capture"}`)
				})

				Convey("Leaves DB unaffected", func() {
					merchant, _ := db.GetMerchant(mockMerchantID)
					So(merchant.Transactions, ShouldResemble, []model.Transaction{
						model.Transaction{
							ID:         mockTransactionID,
							CardNumber: mockCardNumber,
							Authorized: 100,
							Captured:   0,
						}})
				})
			})

			Convey("When capture is not possible", func() {
				requestBody := captureRequestBody(mockMerchantID, mockTransactionID, 70)
				output := simulatePost(testConfig, "/v1/merchant/capture-transaction", requestBody)
				Convey("Returns 200", func() {
					So(output.Code, ShouldEqual, 200)
				})

				Convey("Moves funds from Authorized Transactions to captured transactions", func() {
					merchant, _ := db.GetMerchant(mockMerchantID)
					So(merchant.Transactions, ShouldResemble, []model.Transaction{
						model.Transaction{
							ID:         mockTransactionID,
							CardNumber: mockCardNumber,
							Authorized: 30,
							Captured:   70,
						}})
				})
			})

			Convey("When request is badly formed", func() {
				output := simulatePost(testConfig, "/v1/merchant/capture-transaction", strings.NewReader("bad data"))
				Convey("Returns 400", func() {
					So(output.Code, ShouldEqual, 400)
				})

				Convey("Returns bad JSON response", func() {
					bodyAsString := output.Body.String()
					So(bodyAsString, ShouldEqual, `{"error":"bad JSON format"}`)
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
func captureRequestBody(merchantID string, tranasactionID string, transactionAmount int32) io.Reader {
	return strings.NewReader(fmt.Sprintf(
		`{
			"merchant_id": "%s",
			"transaction_id": "%s",
			"amount": %d
		}`,
		merchantID,
		tranasactionID,
		transactionAmount,
	))
}
