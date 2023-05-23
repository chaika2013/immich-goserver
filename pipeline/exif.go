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
	assetPath, err := model.GetAssetPathByID(p.ID)
	if err != nil {
		return err
	}

	et, err := exiftool.NewExiftool()
	if err != nil {
		return err
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(assetPath)
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
		// Orientation:
		CreateDate:   toDateTime(&fileInfo, "CreateDate"),
		ModifyDate:   toDateTime(&fileInfo, "ModifyDate"),
		LensModel:    toString(&fileInfo, "LensModel"),
		FNumber:      toFloat(&fileInfo, "FNumber"),
		FocalLength:  toFloat(&fileInfo, "FocalLength"),
		ISO:          toInt(&fileInfo, "ISO"),
		ExposureTime: toString(&fileInfo, "ExposureTime"),
	}

	fmt.Println(exif)
	return nil
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

func toFloat(fileInfo *exiftool.FileMetadata, k string) *float32 {
	if k == "FocalLength" {
		r := toString(fileInfo, k)
		if r == nil {
			return nil
		}
		r1 := strings.Trim(*r, " m")
		r2, err := strconv.ParseFloat(r1, 32)
		if err != nil {
			return nil
		}
		r3 := float32(r2)
		return &r3
	}
	r, err := fileInfo.GetFloat(k)
	if err != nil {
		return nil
	}
	r1 := float32(r)
	return &r1
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
