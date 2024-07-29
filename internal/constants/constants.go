package constants

import "path/filepath"

const (
	PhotoSwipeVersion  = "5.4.2"
	CacheDirectoryName = ".foto"
	CacheVersion       = "2"

	PhotosURLPath string = "/photos/"
)

var (
	TemplateFilePath string = filepath.Join("templates", "template.html")
)
