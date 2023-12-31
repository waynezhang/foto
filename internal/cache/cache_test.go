package cache

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/utils"
)

func TestCache(t *testing.T) {
	testfile := "../../testdata/testfile.jpg"
	resizedFile := "../../testdata/testfile-640.jpg"
	expectedChecksum := "2786728c2c9eb5334df492e1853e24c72f976e063ebd513b45bc47476178cc23"
	expectedResizedChecksum := "1c8a6195eefb53be554d86df9de1ae7c5559fa71938be1db595c3bef6c063796"

	dirName, err := os.MkdirTemp("", "foto-cache")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	instance = New(dirName)

	assert.Equal(t, dirName, instance.directoryName)

	img := instance.CachedImage(testfile, 640)
	assert.Nil(t, img)

	instance.AddImage(testfile, 640, resizedFile)
	img = instance.CachedImage(testfile, 640)
	expectedPath := fmt.Sprintf("%s/%s-640", dirName, expectedChecksum)
	assert.Equal(t, expectedPath, *img)

	resizedChecksum, _ := utils.FileChecksum(expectedPath)
	assert.Equal(t, expectedResizedChecksum, *resizedChecksum)

	instance.Clear()
	assert.NoFileExists(t, dirName)
}

func TestShared(t *testing.T) {
	instance := Shared()
	assert.Equal(t, constants.CacheDirectoryName, instance.directoryName)
}

func TestImagePath(t *testing.T) {
	instance := New("some-path")
	path := instance.imagePath("some-checksum", 200)
	assert.Equal(t, "some-path/some-checksum-200", path)
}
