package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAssetCountByTimeBucketNoLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/count-by-time-bucket", bytes.NewBuffer(
		[]byte(`{"timeGroup":"month"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestGetAssetCountByTimeBucket(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/count-by-time-bucket", bytes.NewBuffer(
		[]byte(`{"timeGroup":"month"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `
		{
			"totalCount":15,
			"buckets":[
				{"count":5,"timeBucket":"2000-03-01T00:00:00.000Z"},
				{"count":5,"timeBucket":"2000-02-01T00:00:00.000Z"},
				{"count":5,"timeBucket":"2000-01-01T00:00:00.000Z"}
			]
		}`,
		w.Body.String())
}

func TestGetAssetByTimeBucketNoLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/time-bucket", bytes.NewBuffer(
		[]byte(`{"timeBucket":"2000-03-01T00:00:00.000Z"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestGetAssetByTimeBucket(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 1, 1)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/time-bucket", bytes.NewBuffer(
		[]byte(`{"timeBucket":["2000-01-01T00:00:00.000Z"]}`)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `
		[
			{
				"id":"1",
				"type":"IMAGE",
				"deviceAssetId":"IMG_0_0.jpg-12345",
				"ownerId":"1",
				"deviceId":"CLI",
				"originalFileName":"IMG_0_0.jpg",
				"fileCreatedAt":"2000-01-01T00:00:00Z",
				"isFavorite":false,
				"isArchived":false,
				"duration":"0:00:00.000000"
			}
		]`,
		w.Body.String())
}
