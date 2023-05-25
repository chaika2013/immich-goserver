package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllJobsStatusNotLoggedIn(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	fmt.Println(w.Body.String())
}

func TestGetAllJobsStatusNotAdmin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	fmt.Println(w.Body.String())
}

func TestGetAllJobsStatus(t *testing.T) {
	// router := FromScratch(t)
	// AddTestUsers(t)
	// token := FakeLogin(t, router, 2)

	// w := httptest.NewRecorder()
	// req, _ := http.NewRequest("GET", "/jobs", nil)
	// req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	// router.e.ServeHTTP(w, req)

	// assert.Equal(t, 200, w.Code)
	// fmt.Println(w.Body.String())
	// // assert.JSONEq(t, `
	// // 	[
	// // 		{
	// // 			"id":"1",
	// // 			"type":"IMAGE",
	// // 			"deviceAssetId":"IMG_0_0.jpg-12345",
	// // 			"ownerId":"1",
	// // 			"deviceId":"CLI",
	// // 			"originalFileName":"IMG_0_0.jpg",
	// // 			"fileCreatedAt":"2000-01-01T00:00:00Z",
	// // 			"isFavorite":false,
	// // 			"isArchived":false,
	// // 			"duration":"0:00:00.000000"
	// // 		}
	// // 	]`,
	// // 	w.Body.String())
}
