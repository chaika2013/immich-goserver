package pipeline

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/helper"
	"github.com/chaika2013/immich-goserver/model"
)

func (p *asset) generateThumbnail() error {
	asset, err := model.GetAssetWithRealPathByID(p.ID)
	if err != nil {
		return err
	}

	err = nil
	if asset.AssetType == "IMAGE" {
		err = generateImageThumbnail(asset)
	} else {
		err = fmt.Errorf("unknown method to create thumbnail for asset type %s",
			asset.AssetType)
	}

	return err
}

func generateImageThumbnail(asset *model.Asset) error {
	userAssetPath := filepath.Join(*config.ThumbnailPath, helper.StringID(asset.UserID))

	if err := os.MkdirAll(userAssetPath, os.ModePerm); err != nil {
		return err
	}

	if err := imageToWebpThumbnail(userAssetPath, asset); err != nil {
		return err
	}

	return nil
}

func imageToWebpThumbnail(userAssetPath string, asset *model.Asset) error {
	var stderr bytes.Buffer
	webpFilename := fmt.Sprintf("%d.webp", asset.ID)
	webpPath := filepath.Join(userAssetPath, webpFilename)
	cmd := exec.Command("convert", asset.AssetPath, "-auto-orient", "-thumbnail", "250x250", "-unsharp", "0x.5", webpPath)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	return nil
}
