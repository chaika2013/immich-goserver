package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/view"
	"github.com/stretchr/testify/assert"
)

func TestGetAssetCountByTimeBucketNoLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5, false)

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
	AddTestAssetsForUser(t, 1, 3, 5, false)
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
	}`, w.Body.String())
}

func TestGetAssetByTimeBucketNoLogin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5, false)

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
	AddTestAssetsForUser(t, 1, 1, 1, false)
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
			"deviceAssetId":"IMG_0_0.jpg-0",
			"ownerId":"1",
			"deviceId":"CLI",
			"type":"IMAGE",
			"originalPath":"-",
			"originalFileName":"IMG_0_0.jpg",
			"fileCreatedAt":"2000-01-01T00:00:00Z",
			"isFavorite":false,
			"isArchived":false,
			"duration":"0:00:00.000000"
		}
	]`, w.Body.String())
}

func TestCheckDuplicateAssetFoundDuplicate(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 10, 1, false)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/check", bytes.NewBuffer(
		[]byte(`{"deviceAssetId":"IMG_0_0.jpg-0","deviceId":"CLI"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"id":"1", "isExist":true}`, w.Body.String())
}

func TestCheckDuplicateAssetNoDuplicate(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 10, 1, false)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/check", bytes.NewBuffer(
		[]byte(`{"deviceAssetId":"IMG_0_0.jpg-1","deviceId":"CLI"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"isExist":false}`, w.Body.String())
}

func TestGetAssetCountByTimeBucketWithEmptyBucket(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 3, 5, true)
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
		"totalCount":20,
		"buckets":[
			{"count":5,"timeBucket":"2000-03-01T00:00:00.000Z"},
			{"count":5,"timeBucket":"2000-02-01T00:00:00.000Z"},
			{"count":5,"timeBucket":"2000-01-01T00:00:00.000Z"},
			{"count":5,"timeBucket":""}
		]
	}`, w.Body.String())
}

func TestGetAssetByTimeBucketWithEmptyBucket(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 0, 1, true)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/time-bucket", bytes.NewBuffer(
		[]byte(`{"timeBucket":[""]}`)))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `
	[
		{
			"id":"1",
			"deviceAssetId":"IMG_0_0.jpg-0",
			"ownerId":"1",
			"deviceId":"CLI",
			"type":"IMAGE",
			"originalPath":"-",
			"originalFileName":"IMG_0_0.jpg",
			"isFavorite":false,
			"isArchived":false,
			"duration":"0:00:00.000000"
		}
	]`, w.Body.String())
}

func TestUploadFile(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 2)

	body, contentType := UploadFile(t, "file_example_JPG_100kB.jpg", &view.UploadFile{
		AssetType:      "IMAGE",
		DeviceAssetID:  "file_example_JPG_100kB.jpg-102117",
		DeviceID:       "CLI",
		FileCreatedAt:  time.Date(2009, 1, 2, 12, 34, 56, 0, time.UTC),
		FileModifiedAt: time.Date(2103, 2, 1, 21, 43, 12, 0, time.UTC),
		IsFavorite:     false,
		IsArchived:     false,
		IsVisible:      false,
		FileExtension:  "jpg",
		Duration:       "0:0.0",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/asset/upload", body)
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `{"duplicate":false,"id":"1"}`, w.Body.String())

	// verify asset in the database
	var asset model.Asset
	model.DB.First(&asset)
	assert.EqualValues(t, 1, asset.ID)
	assert.EqualValues(t, 2, asset.UserID)
	assert.EqualValues(t, "IMAGE", asset.AssetType)
	assert.EqualValues(t, "CLI", asset.DeviceID)
	assert.EqualValues(t, "file_example_JPG_100kB.jpg-102117", asset.DeviceAssetID)
	assert.EqualValues(t, time.Date(2009, 1, 2, 12, 34, 56, 0, time.UTC), asset.FileCreatedAt)
	assert.EqualValues(t, time.Date(2103, 2, 1, 21, 43, 12, 0, time.UTC), asset.FileModifiedAt)
	assert.False(t, asset.IsFavorite)
	assert.False(t, asset.IsArchived)
	assert.False(t, asset.IsVisible)
	assert.EqualValues(t, "0:0.0", asset.Duration)
	assert.EqualValues(t, "file_example_JPG_100kB.jpg", asset.OriginalFileName)
	assert.EqualValues(t, 102117, asset.FileSize)
	assert.EqualValues(t, uint32(0x993d34f3), asset.CRC32)
	assert.Nil(t, asset.DateTimeOriginal)
	assert.False(t, asset.InLibrary)
	assert.True(t, strings.HasPrefix(asset.AssetPath, "file_example_JPG_100kB.jpg-"))
}

func TestGetAssetByID(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	AddTestAssetsForUser(t, 1, 1, 1, false)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/asset/assetById/1", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `
	{
		"id":"1",
		"deviceAssetId":"IMG_0_0.jpg-0",
		"ownerId":"1",
		"deviceId":"CLI",
		"type":"IMAGE",
		"originalPath":"-",
		"originalFileName":"IMG_0_0.jpg",
		"fileCreatedAt":"2000-01-01T00:00:00Z",
		"isFavorite":false,
		"isArchived":false,
		"duration":"0:00:00.000000"
	}`, w.Body.String())
}
