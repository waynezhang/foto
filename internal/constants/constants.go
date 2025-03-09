package constants

import "path/filepath"

const (
	PhotoSwipeVersion  = "5.4.4"
	CacheDirectoryName = ".foto"
	CacheVersion       = "3"

	PhotosURLPath          string = "/photos/"
	DefaultCompressQuality        = 75
)

var (
	TemplateFilePath string = filepath.Join("templates", "template.html")
)
