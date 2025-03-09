package cache

import (
	"fmt"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
)

// Implenmentation

type folderCache struct {
	directoryName string
}

func NewFolderCache(directoryName string) Cache {
	return folderCache{
		directoryName: directoryName,
	}
}

// Purge the cache if it's not compatible
func (cache folderCache) Migrate() {
	ver := cache.version()
	if ver == constants.CacheVersion {
		return
	}

	log.Debug().Msgf("Cache version is not compatible to new version(%s), purging", constants.CacheDirectoryName)

	cache.Clear()
	cache.writeVersion(constants.CacheVersion)
}

// `src` is used to compute checksum, `file` will be copied to the cache
func (cache folderCache) AddImage(src string, width int, height int, compressQuality int, file string) {
	checksum, err := files.Checksum(src)
	if err != nil {
		return
	}

	path := cache.imagePath(*checksum, width, height, compressQuality)
	log.Debug().Msgf("Add cache image %s for %s", path, src)
	err = files.EnsureParentDirectory(path)
	if err != nil {
		return
	}

	_ = cp.Copy(file, path)
}

func (cache folderCache) CachedImage(src string, width int, height int, compressQuality int) *string {
	checksum, err := files.Checksum(src)
	if err != nil {
		log.Warn().Msgf("Failed to generate file hash %s (%s).", src, err.Error())
		return nil
	}

	path := cache.imagePath(*checksum, width, height, compressQuality)
	if !files.IsExisting(path) {
		return nil
	}

	return &path
}

func (cache folderCache) Clear() {
	dir := cache.directoryName
	if !files.IsExisting(dir) {
		log.Warn().Msgf("Failed to find cache directory %s.", dir)
	}
	_ = files.PruneDirectory(dir)
}

func (cache folderCache) imagePath(checksum string, width int, height int, compressQuality int) string {
	return filepath.Join(cache.directoryName, fmt.Sprintf("%s-%d-%d-%d", checksum, width, height, compressQuality))
}

func (cache folderCache) version() string {
	path := filepath.Join(cache.directoryName, "version")
	ver, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(ver)
}

func (cache folderCache) writeVersion(ver string) {
	path := filepath.Join(cache.directoryName, "version")
	files.WriteDataToFile([]byte(ver), path)
}
