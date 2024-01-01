package indexer

import (
	"path/filepath"
	"testing"

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
	data := []config.SectionMetadata{
		testdata.Collection1,
		testdata.Collection2,
	}
	sections := Build(data, defaultOption)
	assert.Equal(t, 2, len(sections))
	assert.Equal(t, testdata.Collection1["title"], sections[0].Title)

	assert.Equal(t, 3, len(sections[0].ImageSets))
	assert.Equal(t, testdata.Collection1FileName3, sections[0].ImageSets[0].FileName)
	assert.Equal(t, 3, len(sections[1].ImageSets))
	assert.Equal(t, testdata.Collection2FileName3, sections[1].ImageSets[0].FileName)
}

func TestParseSection(t *testing.T) {
	metadata := testdata.Collection1

	section := parseSection(metadata)
	assert.Equal(t, testdata.Collection1["title"], section.Title)
	assert.Equal(t, testdata.Collection1["text"], string(section.Text))
	assert.Equal(t, testdata.Collection1["slug"], section.Slug)
	assert.Equal(t, testdata.Collection1["folder"], section.Folder)

	// ascending is false by default
	assert.Nil(t, metadata["ascending"])
	assert.Equal(t, false, section.IsAscending)

	metadata["ascending"] = true
	section = parseSection(metadata)
	assert.Equal(t, true, section.IsAscending)
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
