package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerVersion(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server-info/version", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"major":1,"minor":55,"patch":1}`,
		w.Body.String())
}

func TestPingServer(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server-info/ping", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"res":"pong"}`,
		w.Body.String())
}
