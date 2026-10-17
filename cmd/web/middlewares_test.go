package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock the next handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our secureHeaders middleware. Because
	// secureHeaders *returns* a http.Handler we can call its serveHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to execute it.
	secureHeaders(next).ServeHTTP(rr, r)

	// check the middleware has correctly set the X-Frame-Options in response header or not
	frameOptions := rr.Header().Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// check X-XSS-Protection
	xssProtection := rr.Header().Get("X-XSS-Protection")
	if xssProtection != "1; model=block" {
		t.Errorf("want %q; got %q", "1; model=block", xssProtection)
	}

	// response of the next handler
	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// read body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want %q; got %q", "OK", string(body))
	}

}
