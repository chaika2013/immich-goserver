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

	focalLength := fromFocalLength(&fileInfo)
	assert.NotNil(t, focalLength)
	assert.EqualValues(t, 6.3, *focalLength)
}

func TestOrientationToInt(t *testing.T) {
	{
		fileInfo := exiftool.FileMetadata{
			Fields: map[string]interface{}{
				"Orientation": "Horizontal (normal)",
			},
		}
		orientation := fromOrientation(&fileInfo)
		assert.NotNil(t, orientation)
		assert.Equal(t, 1, *orientation)
	}
	{
		fileInfo := exiftool.FileMetadata{
			Fields: map[string]interface{}{
				"Orientation": "Mirror horizontal and rotate 270 CW",
			},
		}
		orientation := fromOrientation(&fileInfo)
		assert.NotNil(t, orientation)
		assert.Equal(t, 5, *orientation)
	}
	{
		fileInfo := exiftool.FileMetadata{
			Fields: map[string]interface{}{
				"Orientation": "Bad format value",
			},
		}
		orientation := fromOrientation(&fileInfo)
		assert.Nil(t, orientation)
	}
}
