package images

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/test"
	"github.com/waynezhang/foto/internal/utils"
)

func TestExtractSection(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmp)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.New(cachePath)
	assert.NotNil(t, cache)

	metadata := test.Collection1
	option := config.ExtractOption{
		ThumbnailWidth: test.ThumbnailWidth,
		OriginalWidth:  test.OriginalWidth,
	}
	section := extractSection(metadata, option, nil, cache, nil)
	assert.Equal(t, test.Collection1["title"], section.Title)
	assert.Equal(t, test.Collection1["text"], string(section.Text))
	assert.Equal(t, test.Collection1["slug"], section.Slug)
	assert.Equal(t, test.Collection1["folder"], section.Folder)

	assert.Equal(t, option.ThumbnailWidth, section.ImageSets[0].ThumbnailSize.Width)
	assert.Equal(t, option.OriginalWidth, section.ImageSets[0].OriginalSize.Width)

	// test orders

	expectedAscendingFileNames := []string{
		test.Collection1FileName1,
		test.Collection1FileName2,
		test.Collection1FileName3,
	}
	expectedDesendingFileNames := []string{
		test.Collection1FileName3,
		test.Collection1FileName2,
		test.Collection1FileName1,
	}

	// ascending is false by default
	section = extractSection(metadata, option, &tmp, cache, nil)
	assert.Equal(t, expectedDesendingFileNames, []string{
		section.ImageSets[0].FileName,
		section.ImageSets[1].FileName,
		section.ImageSets[2].FileName,
	})

	metadata["ascending"] = true
	section = extractSection(metadata, option, &tmp, cache, nil)
	assert.Equal(t, expectedAscendingFileNames, []string{
		section.ImageSets[0].FileName,
		section.ImageSets[1].FileName,
		section.ImageSets[2].FileName,
	})

	metadata["ascending"] = false
	section = extractSection(metadata, option, &tmp, cache, nil)
	assert.Equal(t, expectedDesendingFileNames, []string{
		section.ImageSets[0].FileName,
		section.ImageSets[1].FileName,
		section.ImageSets[2].FileName,
	})

	// check files
	originalFileFolder := filepath.Join(tmp, section.Slug, "original")
	checksum1, _ := utils.FileChecksum(filepath.Join(originalFileFolder, test.Collection1FileName1))
	assert.Equal(t, test.Collection1OriginalChecksum1, *checksum1)
	checksum2, _ := utils.FileChecksum(filepath.Join(originalFileFolder, test.Collection1FileName2))
	assert.Equal(t, test.Collection1OriginalChecksum2, *checksum2)
	checksum3, _ := utils.FileChecksum(filepath.Join(originalFileFolder, test.Collection1FileName3))
	assert.Equal(t, test.Collection1OriginalChecksum3, *checksum3)
}

func TestExtractImage(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmp)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.New(cachePath)
	assert.NotNil(t, cache)

	option := config.ExtractOption{
		ThumbnailWidth: test.ThumbnailWidth,
		OriginalWidth:  test.OriginalWidth,
	}

	set, err := extractImage("nonexisting.jpg", option, "a-slug", &tmp, cache)
	assert.Nil(t, set)
	assert.NotNil(t, err)

	set, err = extractImage(test.Testfile, option, "a-slug", &tmp, cache)
	assert.Nil(t, err)

	// testfile.jpg {640 480} {2048 1536}
	assert.Equal(t, filepath.Base(test.Testfile), *&set.FileName)
	assert.Equal(t, test.OriginalWidth, *&set.OriginalSize.Width)
	assert.Equal(t, test.OriginalHeight, *&set.OriginalSize.Height)
	assert.Equal(t, test.ThumbnailWidth, *&set.ThumbnailSize.Width)
	assert.Equal(t, test.ThumbnailHeight, *&set.ThumbnailSize.Height)
}

func TestResizeImageCache(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmp)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.New(cachePath)
	assert.NotNil(t, cache)

	assert.Nil(t, cache.CachedImage(test.Testfile, test.ThumbnailWidth))

	err = resizeImageAndCache(test.Testfile, filepath.Join(tmp, "resized.jpg"), test.ThumbnailWidth, cache)
	assert.Nil(t, err)

	image := cache.CachedImage(test.Testfile, test.ThumbnailWidth)
	assert.NotNil(t, image)
	checksum, _ := utils.FileChecksum(*image)
	assert.Equal(t, test.ExpectedThubmnailChecksum, *checksum)
}
