package files

import (
	"os"
	"path/filepath"
)

func OutputIndexFilePath(basePath string) string {
	return filepath.Join(basePath, "index.html")
}

func OutputPhotosFilePath(basePath string) string {
	return filepath.Join(basePath, "photos")
}

func OutputPhotoOriginalFilePath(basePath string, slug string, photoFilePath string) string {
	return filepath.Join(basePath, slug, "original", filepath.Base(photoFilePath))
}

func OutputPhotoThumbnailFilePath(basePath string, slug string, photoFilePath string) string {
	return filepath.Join(basePath, slug, "thumbnail", filepath.Base(photoFilePath))
}

func PruneDirectory(path string) error {
	return os.RemoveAll(path)
}

func EnsureDirectory(path string) error {
	if IsExisting(path) {
		return nil
	}

	return os.MkdirAll(path, 0755)
}

func EnsureParentDirectory(path string) error {
	return EnsureDirectory(filepath.Dir(path))
}

func IsExisting(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
