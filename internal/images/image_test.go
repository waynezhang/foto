package images

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/testdata"
)

func TestPhotoSupport(t *testing.T) {
	assert.True(t, IsPhotoSupported("photo.jpg"))
	assert.True(t, IsPhotoSupported("photo.jpeg"))
	assert.True(t, IsPhotoSupported("photo.webp"))
	assert.True(t, IsPhotoSupported("photo.png"))
	assert.False(t, IsPhotoSupported("photo.xxx"))
}

func TestGetPhotoSize(t *testing.T) {
	size, err := GetPhotoSize(testdata.Testfile)
	assert.Nil(t, err)
	assert.Equal(t, testdata.TestfileWidth, size.Width)
	assert.Equal(t, testdata.TestfileHeight, size.Height)

	// test against image with orientation data
	size, _ = GetPhotoSize(testdata.RotatedImageFile)
	assert.Equal(t, testdata.RotatedImageWidth, size.Width)
	assert.Equal(t, testdata.RotatedImageHeight, size.Height)

	_, err = GetPhotoSize("nonexisting-file.jpg")
	assert.True(t, os.IsNotExist(err))
}

func TestAspectedHeight(t *testing.T) {
	assert.Equal(t, 200, AspectedHeight(ImageSize{200, 200}, 200))
	assert.Equal(t, 100, AspectedHeight(ImageSize{2048, 1024}, 200))
}

func TestAspectedWidth(t *testing.T) {
	assert.Equal(t, 200, AspectedWidth(ImageSize{200, 200}, 200))
	assert.Equal(t, 400, AspectedWidth(ImageSize{2048, 1024}, 200))
}

func TestResizeImage(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized.jpg")

	err = ResizeImage("nonexisting-file.jpg", path, testdata.ThumbnailWidth, 0, testdata.CompressQuality)
	assert.True(t, os.IsNotExist(err))
	assert.False(t, files.IsExisting(path))

	err = ResizeImage(testdata.Testfile, path, testdata.ThumbnailWidth, 0, testdata.CompressQuality)
	assert.Nil(t, err)

	checksum, _ := files.Checksum(path)
	assert.Equal(t, testdata.ExpectedThubmnailChecksum, *checksum)
}

func TestResizeWithRotation(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized.jpg")

	err = ResizeImage(testdata.RotatedImageFile, path, testdata.ThumbnailWidth, 0, testdata.CompressQuality)
	assert.Nil(t, err)

	size, _ := GetPhotoSize(path)
	assert.Equal(t, testdata.RotatedImageThumbnailWidth, size.Width)
	assert.Equal(t, testdata.RotatedImageThumbnailHeight, size.Height)
}

func TestCompressQuality(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized.jpg")

	err = ResizeImage(testdata.Testfile, path, testdata.ThumbnailWidth, 0, testdata.CompressQualityHQ)
	assert.Nil(t, err)

	checksum, _ := files.Checksum(path)
	assert.Equal(t, testdata.ExpectedThubmnailHQChecksum, *checksum)
}

func TestWebpSupport(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	size, err := GetPhotoSize(testdata.WebpTestFile)
	assert.Equal(t, testdata.WebpTestfileWidth, size.Width)
	assert.Equal(t, testdata.WebpTestfileHeight, size.Height)

	path := filepath.Join(tmp, "resized.webp")

	err = ResizeImage(testdata.WebpTestFile, path, testdata.WebpThumbnailWidth, 0, testdata.CompressQuality)
	assert.Nil(t, err)

	size, err = GetPhotoSize(path)
	assert.Equal(t, testdata.WebpThumbnailWidth, size.Width)
	assert.Equal(t, testdata.WebpThumbnailHeight, size.Height)

	checksum, _ := files.Checksum(path)
	assert.Equal(t, testdata.WebpExpectedThubmnailChecksum, *checksum)
}

func TestPngSupport(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	size, err := GetPhotoSize(testdata.PngTestFile)
	assert.Equal(t, testdata.PngTestfileWidth, size.Width)
	assert.Equal(t, testdata.PngTestfileHeight, size.Height)

	path := filepath.Join(tmp, "resized.png")

	err = ResizeImage(testdata.PngTestFile, path, testdata.PngThumbnailWidth, 0, testdata.CompressQuality)
	assert.Nil(t, err)

	size, err = GetPhotoSize(path)
	assert.Equal(t, testdata.PngThumbnailWidth, size.Width)
	assert.Equal(t, testdata.PngThumbnailHeight, size.Height)

	checksum, _ := files.Checksum(path)
	assert.Equal(t, testdata.PngExpectedThubmnailChecksum, *checksum)
}

func TestResizeWithHeightOnly(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized-height.jpg")

	// Test resize with only height specified
	err = ResizeImage(testdata.Testfile, path, 0, testdata.ThumbnailHeight, testdata.CompressQuality)
	assert.Nil(t, err)

	size, _ := GetPhotoSize(path)
	assert.Equal(t, testdata.ThumbnailHeight, size.Height)
	// Width should be calculated to maintain aspect ratio
	expectedWidth := AspectedWidth(ImageSize{testdata.TestfileWidth, testdata.TestfileHeight}, testdata.ThumbnailHeight)
	assert.Equal(t, expectedWidth, size.Width)
}

func TestResizeWithBothDimensions(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	path := filepath.Join(tmp, "resized-both.jpg")

	// Test resize with both width and height specified (aspect ratio not preserved)
	customWidth := 400
	customHeight := 400
	err = ResizeImage(testdata.Testfile, path, customWidth, customHeight, testdata.CompressQuality)
	assert.Nil(t, err)

	size, _ := GetPhotoSize(path)
	assert.Equal(t, customWidth, size.Width)
	assert.Equal(t, customHeight, size.Height)
}
