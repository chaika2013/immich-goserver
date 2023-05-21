package helper

import (
	"os"
	"path/filepath"

	"github.com/chaika2013/immich-goserver/config"
	"github.com/chaika2013/immich-goserver/model"
)

// MakeUserUploadDir creates and returns the upload path for the user
func MakeUserUploadDir(user *model.User) (string, error) {
	path := filepath.Join(*config.UploadPath, StringID(user.ID))
	return path, os.MkdirAll(path, os.ModePerm)
}
