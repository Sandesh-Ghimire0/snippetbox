package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// records the response status code, headers and body instead of actually writing them to a HTTP connection.
	rr := httptest.NewRecorder()

	// dummy http.Requeset
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// call the ping handler with by passing response recorder and dummy http.request
	ping(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// check response body returned by the Ping handler
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
