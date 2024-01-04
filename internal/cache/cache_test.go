package cache

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/testdata"
)

func TestFolderCache(t *testing.T) {
	dirName, err := os.MkdirTemp("", "foto-cache")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	cache := NewFolderCache(dirName).(folderCache)

	assert.Equal(t, dirName, cache.directoryName)

	img := cache.CachedImage(testdata.Testfile, 640)
	assert.Nil(t, img)

	cache.AddImage(testdata.Testfile, 640, testdata.ThumbnailFile)
	img = cache.CachedImage(testdata.Testfile, 640)
	expectedPath := fmt.Sprintf("%s/%s-640", dirName, testdata.ExpectedChecksum)
	assert.Equal(t, expectedPath, *img)

	resizedChecksum, _ := files.Checksum(expectedPath)
	assert.Equal(t, testdata.ExpectedThubmnailChecksum, *resizedChecksum)

	// no failure on invalid file
	cache.AddImage("nonexisting-file.jpg", 640, testdata.ThumbnailFile)
	img = cache.CachedImage("nonexisting-file.jpg", 640)
	assert.Nil(t, img)

	cache.Clear()
	assert.NoFileExists(t, dirName)
}

func TestShared(t *testing.T) {
	cache := Shared().(folderCache)
	assert.Equal(t, constants.CacheDirectoryName, cache.directoryName)
}

func TestImagePath(t *testing.T) {
	cache := NewFolderCache("some-path").(folderCache)
	path := cache.imagePath("some-checksum", 200)
	assert.Equal(t, "some-path/some-checksum-200", path)
}
