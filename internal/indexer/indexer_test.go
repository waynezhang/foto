package indexer

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/testdata"
)

var (
	defaultOption = config.ExtractOption{
		ThumbnailWidth: testdata.ThumbnailWidth,
		OriginalWidth:  testdata.OriginalWidth,
	}
)

func TestBuild(t *testing.T) {
	var meta1 config.SectionMetadata
	var meta2 config.SectionMetadata
	mapstructure.Decode(testdata.Collection1, &meta1)
	mapstructure.Decode(testdata.Collection2, &meta2)
	data := []config.SectionMetadata{meta1, meta2}
	sections := Build(data, defaultOption)
	assert.Equal(t, 2, len(sections))
	assert.Equal(t, testdata.Collection1["title"], sections[0].Title)

	assert.Equal(t, 3, len(sections[0].ImageSets))
	assert.Equal(t, testdata.Collection1FileName3, sections[0].ImageSets[0].FileName)
	assert.Equal(t, 3, len(sections[1].ImageSets))
	assert.Equal(t, testdata.Collection2FileName3, sections[1].ImageSets[0].FileName)
}

func TestBuildImageSets(t *testing.T) {
	expectedAscendingFileNames := []string{
		testdata.Collection1FileName1,
		testdata.Collection1FileName2,
		testdata.Collection1FileName3,
	}
	expectedDesendingFileNames := []string{
		testdata.Collection1FileName3,
		testdata.Collection1FileName2,
		testdata.Collection1FileName1,
	}

	folder := testdata.Collection1["folder"].(string)

	sets := buildImageSets(folder, true, defaultOption)
	assert.Equal(t, expectedAscendingFileNames, []string{
		sets[0].FileName,
		sets[1].FileName,
		sets[2].FileName,
	})

	sets = buildImageSets(folder, false, defaultOption)
	assert.Equal(t, expectedDesendingFileNames, []string{
		sets[0].FileName,
		sets[1].FileName,
		sets[2].FileName,
	})
}

func TestBuildImageSet(t *testing.T) {
	set, _ := buildImageSet(testdata.Testfile, defaultOption)
	assert.Equal(t, filepath.Base(testdata.Testfile), set.FileName)
	assert.Equal(t, testdata.ThumbnailWidth, set.ThumbnailSize.Width)
	assert.Equal(t, testdata.ThumbnailHeight, set.ThumbnailSize.Height)
	assert.Equal(t, testdata.OriginalWidth, set.OriginalSize.Width)
	assert.Equal(t, testdata.OriginalHeight, set.OriginalSize.Height)
}
