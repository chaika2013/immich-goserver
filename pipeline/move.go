package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/helper"
	"github.com/chaika2013/immich-goserver/model"
)

func (p *asset) moveToLibrary() error {
	var newFile string

	asset, err := model.GetAssetWithRealPathByID(p.ID)
	if err != nil {
		return err
	}

	// generate new asset path
	// TODO take format from configuration
	folderPath := generateAssetDir(asset.DateTimeOriginal)
	newAssetPath := filepath.Join(folderPath, asset.OriginalFileName)
	fileSuffix := filepath.Ext(newAssetPath)
	filePrefix := newAssetPath[:len(newAssetPath)-len(fileSuffix)]

	// create path
	if err := os.MkdirAll(
		filepath.Join(*config.LibraryPath, helper.StringID(asset.UserID), folderPath),
		os.ModePerm); err != nil {
		return err
	}

	// hard-link file
	index := 0
	for {
		newFile = filepath.Join(*config.LibraryPath, helper.StringID(asset.UserID), newAssetPath)
		err := os.Link(asset.AssetPath, newFile)
		if err == nil {
			break
		} else if os.IsExist(err) {
			if index++; index >= 10000 {
				return err
			}
			newAssetPath = fmt.Sprintf("%s-%d%s", filePrefix, index, fileSuffix)
		} else {
			return err
		}
	}

	// update db
	if err := model.MoveAssetToLibrary(asset, newAssetPath); err != nil {
		// if db update is failed, remove new file
		os.Remove(newFile)
		return err
	}

	// remove old file if all is ok
	if err := os.Remove(asset.AssetPath); err != nil {
		// TODO log that file remove failed
		_ = err
	}
	return nil
}

func generateAssetDir(dateTimeOriginal *time.Time) string {
	if dateTimeOriginal == nil {
		return "blank_date"
	}
	return dateTimeOriginal.Format("2006/2006-01-02")
}
