package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/log"
)

var cacheDirectory = ".foto"

func AddImage(src string, width int, file string) {
	checksum, err := checksum(src)
	if err != nil {
		return
	}

	path := imagePath(*checksum, width)
	log.Debug("Add cache image %s for %s", path, src)
	err = files.EnsureParentDirectory(path)
	if err != nil {
		return
	}

	_ = cp.Copy(file, path)
}

func CachedImage(src string, width int) *string {
	checksum, err := checksum(src)
	if err != nil {
		log.Fatal("Failed to generate file hash %s (%s).", src, err.Error())
		return nil
	}

	path := imagePath(*checksum, width)
	if !files.IsExisting(path) {
		return nil
	}
	
	return &path
}

func Clear() {
  if _, err := os.Stat(cacheDirectory); err != nil {
    if !os.IsNotExist(err) {
      log.Fatal("Failed to find cache directory %s (%s).", cacheDirectory, err.Error()) 
    }
    return
  }
  err := os.RemoveAll(cacheDirectory) 
  if err != nil { 
    log.Fatal("Failed to remove cache directory %s (%s).", cacheDirectory, err.Error()) 
  }
}

func checksum(path string) (*string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, nil
	}

	value := hex.EncodeToString(hasher.Sum(nil))
	return &value, nil
}

func imagePath(checksum string, width int) string {
	return filepath.Join(cacheDirectory, fmt.Sprintf("%s-%d", checksum, width))
}
