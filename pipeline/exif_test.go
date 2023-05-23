package pipeline

import (
	"testing"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/stretchr/testify/assert"
)

func TestCreateDateToDateTime(t *testing.T) {
	fileInfo := exiftool.FileMetadata{
		Fields: map[string]interface{}{
			"CreateDate": "2003:01:08 18:12:09",
		},
	}

	createDate := toDateTime(&fileInfo, "CreateDate")
	assert.NotNil(t, createDate)
	assert.Equal(t, time.Date(2003, time.Month(1), 8, 18, 12, 9, 0, time.UTC), *createDate)
}

func TestFocalLengthToFloat(t *testing.T) {
	fileInfo := exiftool.FileMetadata{
		Fields: map[string]interface{}{
			"FocalLength": "6.3 mm",
		},
	}

	focalLength := toFloat(&fileInfo, "FocalLength")
	assert.NotNil(t, focalLength)
	assert.Equal(t, float32(6.3), *focalLength)
}
