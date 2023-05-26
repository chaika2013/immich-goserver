package test

import (
	"encoding/json"
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

func TestGetMyUserInfo(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/me", nil)
	req.Header.Add("X-Api-Key", token)
	router.e.ServeHTTP(w, req)

	var msg map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msg)
	createdAt := msg["createdAt"].(string)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `
	{
		"id": 1,
		"email": "test1.user1@gmail.com",
		"firstName": "Test1",
		"lastName": "User1",
		"createdAt": "`+createdAt+`",
		"shouldChangePassword": false,
		"isAdmin": true
	} `, w.Body.String())
}

func TestGetAllUsersNoLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?isAll=false", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestGetAllUsersNotAdmin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?isAll=false", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestGetAllUsers(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?isAll=false", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	var msg []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msg)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `
	[
		{
			"id": 2,
			"email": "test2.user2@gmail.com",
			"firstName": "Test2",
			"lastName": "User2",
			"createdAt": "`+msg[0]["createdAt"].(string)+`",
			"shouldChangePassword": false,
			"isAdmin": false
		}
	]`, w.Body.String())
}

func TestGetAllUsersIsAllTrue(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?isAll=true", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	var msg []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msg)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `
	[
		{
			"id": 1,
			"email": "test1.user1@gmail.com",
			"firstName": "Test1",
			"lastName": "User1",
			"createdAt": "`+msg[0]["createdAt"].(string)+`",
			"shouldChangePassword": false,
			"isAdmin": true
		},
		{
			"id": 2,
			"email": "test2.user2@gmail.com",
			"firstName": "Test2",
			"lastName": "User2",
			"createdAt": "`+msg[1]["createdAt"].(string)+`",
			"shouldChangePassword": false,
			"isAdmin": false
		}
	]`, w.Body.String())
}
