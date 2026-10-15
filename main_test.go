package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestHealthz(t *testing.T) {
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rr.Code != http.StatusOK || rr.Body.String() != "ok\n" {
		t.Fatalf("got %d %q", rr.Code, rr.Body.String())
	}
}

func TestWhoami(t *testing.T) {
	signingKey = []byte("drill")
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "drill"}).SignedString(signingKey)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	newHandler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "hello drill\n" {
		t.Fatalf("got %d %q", rr.Code, rr.Body.String())
	}
	rr = httptest.NewRecorder()
	newHandler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/whoami", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("missing token: got %d", rr.Code)
	}
}
