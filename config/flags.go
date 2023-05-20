package config

import "flag"

// DatabasePath where the sqlite file is stored
var DatabasePath = flag.String(
	"database-path",
	"/var/lib/immich/database/immich.sqlite",
	"sqlite database file path",
)

// UploadPath where the uploaded files are stored
var UploadPath = flag.String(
	"upload-path",
	"/var/lib/immich/upload",
	"path for uploaded files",
)

// LibraryPath where the original assets are stored
var LibraryPath = flag.String(
	"library-path",
	"/var/lib/immich/library",
	"path where original assets are sored",
)

// ThumbnailPath where the thumbnails are stored
var ThumbnailPath = flag.String(
	"thumbnail-path",
	"/var/lib/immich/thumbnail",
	"path where thumbnails are sored",
)

// EncodedPath where the encoded videos are stored
var EncodedVideoPath = flag.String(
	"encoded-video-path",
	"/var/lib/immich/encoded-video",
	"path where encoded videos are sored",
)
