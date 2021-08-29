package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	rs := ts.Get(t, "/ping")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	if string(rs.Body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
