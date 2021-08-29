package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/terdia/snippetbox/pkg/logger"
)

type IntegrationTestResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func newTestApplication(t *testing.T) *application {
	return &application{
		logger: &logger.SnippetLogger{
			Info:  log.New(io.Discard, "", 0),
			Error: log.New(io.Discard, "", 0),
		},
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) Get(t *testing.T, urlPath string) IntegrationTestResponse {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return IntegrationTestResponse{
		StatusCode: rs.StatusCode,
		Header:     rs.Header,
		Body:       body,
	}
}
