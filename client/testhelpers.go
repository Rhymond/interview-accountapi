package client

import (
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
)

func createTestServer() (*mux.Router, *httptest.Server, *Client) {
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	su, _ := url.Parse(server.URL)
	return router, server, NewClient(nil, su)
}