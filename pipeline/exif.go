package pipeline

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/chaika2013/immich-goserver/model"
)

func (p *asset) extractExif() error {
	asset, err := model.GetAssetWithRealPathByID(p.ID)
	if err != nil {
		return err
	}

	et, err := exiftool.NewExiftool()
	if err != nil {
		return err
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(asset.AssetPath)
	if len(fileInfos) != 1 {
		return fmt.Errorf("EXIF not found")
	}

	fileInfo := fileInfos[0]
	if fileInfo.Err != nil {
		return fileInfo.Err
	}

	// put exif data into the database
	exif := model.Exif{
		AssetID:         p.ID,
		Make:            toString(&fileInfo, "Make"),
		CameraModel:     toString(&fileInfo, "Model"),
		ExifImageWidth:  toInt(&fileInfo, "ExifImageWidth"),
		ExifImageHeight: toInt(&fileInfo, "ExifImageHeight"),
		Orientation:     fromOrientation(&fileInfo),
		CreateDate:      toDateTime(&fileInfo, "CreateDate"),
		ModifyDate:      toDateTime(&fileInfo, "ModifyDate"),
		LensModel:       toString(&fileInfo, "LensModel"),
		FNumber:         toFloat(&fileInfo, "FNumber"),
		FocalLength:     fromFocalLength(&fileInfo),
		ISO:             toInt(&fileInfo, "ISO"),
		ExposureTime:    toString(&fileInfo, "ExposureTime"),
	}
	return model.UpsertExif(&exif, toDateTime(&fileInfo, "DateTimeOriginal"))
}

func toString(fileInfo *exiftool.FileMetadata, k string) *string {
	r, err := fileInfo.GetString(k)
	if err != nil {
		return nil
	}
	return &r
}

func toInt(fileInfo *exiftool.FileMetadata, k string) *int {
	r, err := fileInfo.GetInt(k)
	if err != nil {
		return nil
	}
	r1 := int(r)
	return &r1
}

func toFloat(fileInfo *exiftool.FileMetadata, k string) *float64 {
	r, err := fileInfo.GetFloat(k)
	if err != nil {
		return nil
	}
	return &r
}

func toDateTime(fileInfo *exiftool.FileMetadata, k string) *time.Time {
	const layout = "2006:01:02 15:04:05"
	r, err := fileInfo.GetString(k)
	if err != nil {
		return nil
	}
	r1, err := time.Parse(layout, r)
	if err != nil {
		return nil
	}
	return &r1
}

func fromFocalLength(fileInfo *exiftool.FileMetadata) *float64 {
	r := toString(fileInfo, "FocalLength")
	if r == nil {
		return nil
	}
	r1 := strings.Trim(*r, " m")
	r2, err := strconv.ParseFloat(r1, 64)
	if err != nil {
		return nil
	}
	return &r2
}

func fromOrientation(fileInfo *exiftool.FileMetadata) *int {
	r := toString(fileInfo, "Orientation")
	if r == nil {
		return nil
	}
	orientationMap := map[string]int{
		"Horizontal (normal)":                 1,
		"Mirror horizontal":                   2,
		"Rotate 180":                          3,
		"Mirror vertical":                     4,
		"Mirror horizontal and rotate 270 CW": 5,
		"Rotate 90 CW":                        6,
		"Mirror horizontal and rotate 90 CW":  7,
		"Rotate 270 CW":                       8,
	}
	if r1, found := orientationMap[*r]; found {
		return &r1
	}
	return nil
}
