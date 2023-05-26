package test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/model"
	"github.com/chaika2013/immich-goserver/router"
	"github.com/chaika2013/immich-goserver/session"
	"github.com/chaika2013/immich-goserver/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func FromScratch(t *testing.T) *RouterCtx {
	var err error

	// remove test folder
	os.RemoveAll(".test")

	// setup folder structure
	os.MkdirAll(".test/upload", os.ModePerm)
	os.MkdirAll(".test/library", os.ModePerm)
	os.MkdirAll(".test/thumbnail", os.ModePerm)
	os.MkdirAll(".test/encoded", os.ModePerm)

	// configuration
	*config.DatabasePath = ".test/immich.sqlite"
	*config.UploadPath = ".test/upload"
	*config.LibraryPath = ".test/library"
	*config.ThumbnailPath = ".test/thumbnail"
	*config.EncodedVideoPath = ".test/encoded"

	// setup database
	model.DB, err = gorm.Open(sqlite.Open(*config.DatabasePath), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	assert.NoError(t, err)

	err = model.DB.AutoMigrate(&model.Asset{}, &model.User{}, &model.Exif{})
	assert.NoError(t, err)

	return NewRouter(t)
}

type RouterCtx struct {
	e *gin.Engine
	s session.Store
}

func NewRouter(t *testing.T) *RouterCtx {
	r := RouterCtx{
		e: gin.Default(),
		s: session.NewStore(),
	}
	r.e.Use(sessions.Sessions("immich_access_token", r.s))
	router.Setup(r.e)
	return &r
}

func AddTestUsers(t *testing.T) {
	// append new users
	user1 := model.User{
		Email:                "test1.user1@gmail.com",
		Password:             "$2a$14$rRKBPSc.syVWf3AqoIvdXOEvb5Dq83WlxaO.La1/30Gt5ysB.TFzS",
		FirstName:            "Test1",
		LastName:             "User1",
		ShouldChangePassword: false,
		IsAdmin:              true,
	}
	err := model.DB.Create(&user1).Error
	assert.NoError(t, err)

	user2 := model.User{
		Email:                "test2.user2@gmail.com",
		Password:             "$2a$14$rRKBPSc.syVWf3AqoIvdXOEvb5Dq83WlxaO.La1/30Gt5ysB.TFzS",
		FirstName:            "Test2",
		LastName:             "User2",
		ShouldChangePassword: false,
		IsAdmin:              false,
	}
	err = model.DB.Create(&user2).Error
	assert.NoError(t, err)
}

func DoLogin(t *testing.T, router *RouterCtx, user string, password string) (string, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(
		[]byte(`{"email":"`+user+`","password":"`+password+`"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.e.ServeHTTP(w, req)

	// get token
	cookie := w.Result().Cookies()[0]
	assert.Equal(t, "immich_access_token", cookie.Name)

	return cookie.Value, w
}

func AddTestAssetsForUser(t *testing.T, userID uint, monthBuckets int, countPerBucket int, emptyBucket bool) {
	ts := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < monthBuckets+1; i++ {
		if i >= monthBuckets && !emptyBucket {
			break
		}
		for j := 0; j < countPerBucket; j++ {
			imageName := fmt.Sprintf("IMG_%d_%d.jpg", i, j)
			fileSize := i*monthBuckets + j
			asset := model.Asset{
				UserID:           userID,
				AssetType:        "IMAGE",
				DeviceID:         "CLI",
				DeviceAssetID:    fmt.Sprintf("%s-%d", imageName, fileSize),
				FileCreatedAt:    time.Now(),
				FileModifiedAt:   time.Now(),
				OriginalFileName: imageName,
				FileSize:         int64(fileSize),
				CRC32:            uint32(i*monthBuckets + j),
				DateTimeOriginal: &ts,
				Duration:         "0:00:00.000000",
			}
			if i == monthBuckets {
				asset.DateTimeOriginal = nil
			}
			err := model.DB.Create(&asset).Error
			assert.NoError(t, err)
		}
		ts = ts.AddDate(0, 1, 0)
	}
}

func FakeLogin(t *testing.T, router *RouterCtx, userID uint) string {
	req, _ := http.NewRequest("POST", "/auth/login", nil)
	session, err := router.s.New(req, "immich_access_token")
	assert.NoError(t, err)

	// set userid and save session
	session.Values["user_id"] = userID
	router.s.Save(req, nil, session)

	return session.ID
}

func UploadFile(t *testing.T, fileName string, assetInfo *view.UploadFile) (*bytes.Buffer, string) {
	var body bytes.Buffer

	file, err := os.Open(filepath.Join("assets", fileName))
	assert.NoError(t, err)
	defer file.Close()

	w := multipart.NewWriter(&body)
	defer w.Close()

	part, _ := w.CreateFormFile("assetData", fileName)
	io.Copy(part, file)
	part, _ = w.CreateFormField("assetType")
	io.Copy(part, strings.NewReader(assetInfo.AssetType))
	part, _ = w.CreateFormField("deviceAssetId")
	io.Copy(part, strings.NewReader(assetInfo.DeviceAssetID))
	part, _ = w.CreateFormField("deviceId")
	io.Copy(part, strings.NewReader(assetInfo.DeviceID))
	part, _ = w.CreateFormField("fileCreatedAt")
	io.Copy(part, strings.NewReader(assetInfo.FileCreatedAt.Format(time.RFC3339Nano)))
	part, _ = w.CreateFormField("fileModifiedAt")
	io.Copy(part, strings.NewReader(assetInfo.FileModifiedAt.Format(time.RFC3339Nano)))
	part, _ = w.CreateFormField("isFavorite")
	io.Copy(part, strings.NewReader(strconv.FormatBool(assetInfo.IsFavorite)))
	part, _ = w.CreateFormField("isArchived")
	io.Copy(part, strings.NewReader(strconv.FormatBool(assetInfo.IsArchived)))
	part, _ = w.CreateFormField("isVisible")
	io.Copy(part, strings.NewReader(strconv.FormatBool(assetInfo.IsVisible)))
	part, _ = w.CreateFormField("fileExtension")
	io.Copy(part, strings.NewReader(assetInfo.FileExtension))
	part, _ = w.CreateFormField("duration")
	io.Copy(part, strings.NewReader(assetInfo.Duration))

	return &body, w.FormDataContentType()
}
