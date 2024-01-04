package cache

import (
	"fmt"
	"path/filepath"
	"sync"

	cp "github.com/otiai10/copy"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/log"
)

// Implenmentation

type folderCache struct {
	directoryName string
}

var (
	once     sync.Once
	instance Cache
)

func NewFolderCache(directoryName string) Cache {
	instance = folderCache{
		directoryName: directoryName,
	}
	return instance
}

// `src` is used to compute checksum, `file` will be copied to the cache
func (cache folderCache) AddImage(src string, width int, file string) {
	checksum, err := files.Checksum(src)
	if err != nil {
		return
	}

	path := cache.imagePath(*checksum, width)
	log.Debug("Add cache image %s for %s", path, src)
	err = files.EnsureParentDirectory(path)
	if err != nil {
		return
	}

	_ = cp.Copy(file, path)
}

func (cache folderCache) CachedImage(src string, width int) *string {
	checksum, err := files.Checksum(src)
	if err != nil {
		log.Fatal("Failed to generate file hash %s (%s).", src, err.Error())
		return nil
	}

	path := cache.imagePath(*checksum, width)
	if !files.IsExisting(path) {
		return nil
	}

	return &path
}

func (cache folderCache) Clear() {
	dir := cache.directoryName
	if !files.IsExisting(dir) {
		log.Fatal("Failed to find cache directory %s.", dir)
	}
	files.PruneDirectory(dir)
}

func (cache folderCache) imagePath(checksum string, width int) string {
	return filepath.Join(cache.directoryName, fmt.Sprintf("%s-%d", checksum, width))
}
