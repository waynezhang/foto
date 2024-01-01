package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/testdata"
)

func TestResizeImageCache(t *testing.T) {
	tmp, cache := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	assert.Nil(t, cache.CachedImage(testdata.Testfile, testdata.ThumbnailWidth))

	_ = resizeImageAndCache(testdata.Testfile, filepath.Join(tmp, "resized.jpg"), testdata.ThumbnailWidth, cache)

	image := cache.CachedImage(testdata.Testfile, testdata.ThumbnailWidth)
	checksum, _ := files.Checksum(*image)
	assert.Equal(t, testdata.ExpectedThubmnailChecksum, *checksum)
}

func prepareTempDirAndCache(t *testing.T) (string, cache.Cache) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.New(cachePath)
	assert.NotNil(t, cache)

	return tmp, cache
}
