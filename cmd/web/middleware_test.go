package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippit/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock http handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass mock to middleware
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()

	// check headers
	want := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), want)

	want = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), want)

	want = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), want)

	want = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), want)

	want = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), want)

	// check status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// check body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
