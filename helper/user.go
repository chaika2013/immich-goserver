package helper

import (
	"os"
	"path/filepath"

	"github.com/chaika2013/immich-goserver/config"
)

// MakeUserUploadDir creates and returns the upload path for the user
func MakeUserUploadDir(userID uint) (string, error) {
	path := filepath.Join(*config.UploadPath, StringID(userID))
	return path, os.MkdirAll(path, os.ModePerm)
}
