package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserCount(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/count?admin=true", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"userCount":1}`,
		w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/count", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"userCount":2}`,
		w.Body.String())
}

func TestGetMyUserInfoNotLoggedIn(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/me", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}
