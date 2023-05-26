package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginWrongEmail(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(
		[]byte(`{"email":"email","password":"password"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestLoginCorrectEmailWithWrongPassword(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(
		[]byte(`{"email":"test1.user1@gmail.com","password":"password"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	token, w := DoLogin(t, router, "test1.user1@gmail.com", "123456")

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `
	{
		"accessToken":"`+token+`",
		"userId":"1",
		"isAdmin":true,
		"userEmail":"test1.user1@gmail.com",
		"firstName":"Test1",
		"lastName":"User1",
		"shouldChangePassword":false,
		"profileImagePath":""
	}`, w.Body.String())
}

func TestLogout(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `{"redirectUri":"","successful":true}`,
		w.Body.String())
}

func TestLogoutAfterLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	token, _ := DoLogin(t, router, "test1.user1@gmail.com", "123456")

	// do logout with the cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "immich_access_token",
		Value: token,
	})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `{"redirectUri":"","successful":true}`,
		w.Body.String())
}
