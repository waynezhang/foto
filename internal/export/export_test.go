package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/test"
)

func TestResizeImageCache(t *testing.T) {
	tmp, cache := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	assert.Nil(t, cache.CachedImage(test.Testfile, test.ThumbnailWidth))

	_ = resizeImageAndCache(test.Testfile, filepath.Join(tmp, "resized.jpg"), test.ThumbnailWidth, cache)

	image := cache.CachedImage(test.Testfile, test.ThumbnailWidth)
	checksum, _ := files.Checksum(*image)
	assert.Equal(t, test.ExpectedThubmnailChecksum, *checksum)
}

func prepareTempDirAndCache(t *testing.T) (string, cache.Cache) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.New(cachePath)
	assert.NotNil(t, cache)

	return tmp, cache
}
