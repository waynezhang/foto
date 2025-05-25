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

func TestAspectedSize(t *testing.T) {
	assert.Equal(t, ImageSize{640, 480}, AspectedSize(ImageSize{2048, 1536}, 640, 0))
	assert.Equal(t, ImageSize{640, 480}, AspectedSize(ImageSize{2048, 1536}, 640, 100))
	assert.Equal(t, ImageSize{1024, 768}, AspectedSize(ImageSize{2048, 1536}, 640, 768))
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

func TestGetImageEXIF(t *testing.T) {
	exif, err := GetEXIFValues(testdata.MetadataTestFile)
	assert.Nil(t, err)

	assert.Equal(t, testdata.ExpectedImageDescription, exif["ImageDescription"])
	assert.Equal(t, testdata.ExpectedMake, exif["Make"])
	assert.Equal(t, testdata.ExpectedModel, exif["Model"])
	assert.Equal(t, testdata.ExpectedExposureTime, exif["ExposureTime"])
	assert.Equal(t, testdata.ExpectedISO, exif["ISO"])
	assert.Equal(t, testdata.ExpectedApertureValue, exif["ApertureValue"])
}

func TestGetEmptyImageDescription(t *testing.T) {
	exif, err := GetEXIFValues(testdata.RotatedImageFile)
	assert.Nil(t, err)
	assert.Equal(t, "", exif["ImageDescription"])
}

func TestGetPngImageDescription(t *testing.T) {
	exif, err := GetEXIFValues(testdata.PngMetadataTestFile)
	assert.Nil(t, err)
	assert.Equal(t, testdata.PngExpectedImageDescription, exif["ImageDescription"])
}

func TestGetEmptyPngImageDescription(t *testing.T) {
	exif, err := GetEXIFValues(testdata.PngTestFile)
	assert.Nil(t, err)
	assert.Equal(t, "", exif["ImageDescription"])
}
