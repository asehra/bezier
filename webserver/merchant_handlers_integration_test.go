package webserver_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/mock"
	"github.com/asehra/bezier/storage"
	"github.com/asehra/bezier/webserver"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMerchantHandlersAPI(t *testing.T) {
	mockMerchantID := "M012345"
	// mockCardNumber := int64(9000000000000001)
	db := storage.NewInMemoryStore()
	testConfig := config.Config{
		DB: db,
		// CardIDGenerator:     &mock.CardIDGenerator{Generates: mockCardNumber},
		MerchantIDGenerator: &mock.MerchantIDGenerator{Generates: mockMerchantID},
		StdErr:              ioutil.Discard,
		StdOut:              ioutil.Discard,
	}

	Convey("GET /v1/merchant/create", t, func() {
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
			So(db.MerchantStore[response.MerchantID].ID, ShouldEqual, mockMerchantID)
		})
	})
}
