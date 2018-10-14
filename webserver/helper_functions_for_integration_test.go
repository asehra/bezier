package webserver_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/webserver"
)

func simulateGet(config config.Config, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()

	api := webserver.Create(config)
	api.ServeHTTP(w, req)
	return w
}

func simulatePost(config config.Config, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	w := httptest.NewRecorder()

	api := webserver.Create(config)
	api.ServeHTTP(w, req)
	return w
}
