package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	cp "github.com/otiai10/copy"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/utils"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/log"
)

type Cache struct {
  directoryName string
}

var (
  once sync.Once
  instance Cache
)

func Shared() Cache {
  once.Do(func() {
    instance = New(constants.CacheDirectoryName)
  })
  return instance
}

func New(directoryName string) Cache {
    instance = Cache{}
    instance.directoryName = directoryName
    return instance
}

// `src` is used to compute checksum, `file` will be copied to the cache
func (cache Cache) AddImage(src string, width int, file string) {
	checksum, err := utils.FileChecksum(src)
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

func (cache Cache) CachedImage(src string, width int) *string {
	checksum, err := utils.FileChecksum(src)
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

func (cache Cache) Clear() {
  if _, err := os.Stat(cache.directoryName); err != nil {
    if !os.IsNotExist(err) {
      log.Fatal("Failed to find cache directory %s (%s).", cache.directoryName, err.Error()) 
    }
    return
  }
  err := os.RemoveAll(cache.directoryName) 
  if err != nil { 
    log.Fatal("Failed to remove cache directory %s (%s).", cache.directoryName, err.Error()) 
  }
}

func (cache Cache) imagePath(checksum string, width int) string {
	return filepath.Join(cache.directoryName, fmt.Sprintf("%s-%d", checksum, width))
}
