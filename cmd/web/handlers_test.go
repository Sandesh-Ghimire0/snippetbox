package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. This just
	// contains a couple of mock loggers (which discard anything written to them)
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	// create a new https server
	// pass the routes and handler, that server needs to test
	// server starts at randomly chosen port
	// defer to close the server after test finishes

	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	// make GET /ping request against the test server
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// check the response status and body
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
