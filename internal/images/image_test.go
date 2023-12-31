package images

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/test"
	"github.com/waynezhang/foto/internal/utils"
)

func TestPhotoSupport(t *testing.T) {
	assert.True(t, IsPhotoSupported("photo.jpg"))
	assert.True(t, IsPhotoSupported("photo.jpeg"))
	assert.False(t, IsPhotoSupported("photo.png"))
}

func TestGetPhotoSize(t *testing.T) {
	size, err := GetPhotoSize(test.Testfile)
	assert.Nil(t, err)
	assert.Equal(t, test.TestfileWidth, size.Width)
	assert.Equal(t, test.TestfileHeight, size.Height)

	_, err = GetPhotoSize("nonexisting-file.jpg")
	assert.True(t, os.IsNotExist(err))
}

func TestResizeImage(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized.jpg")

	err = ResizeImage("nonexisting-file.jpg", path, test.ResizedWidth)
	assert.True(t, os.IsNotExist(err))
	assert.False(t, files.IsExisting(path))

	err = ResizeImage(test.Testfile, path, test.ResizedWidth)
	assert.Nil(t, err)

	checksum, err := utils.FileChecksum(path)
	assert.Equal(t, test.ExpectedResizedChecksum, *checksum)
}
