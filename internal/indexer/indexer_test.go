package indexer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/test"
)

var (
	defaultOption = config.ExtractOption{
		ThumbnailWidth: test.ThumbnailWidth,
		OriginalWidth:  test.OriginalWidth,
	}
)

func TestBuild(t *testing.T) {
	data := []config.SectionMetadata{
		test.Collection1,
		test.Collection2,
	}
	sections := Build(data, defaultOption)
	assert.Equal(t, 2, len(sections))
	assert.Equal(t, test.Collection1["title"], sections[0].Title)
	assert.Equal(t, 3, len(sections[0].ImageSets))
	assert.Equal(t, test.Collection1FileName3, sections[0].ImageSets[0].FileName)
	assert.Equal(t, 3, len(sections[1].ImageSets))
	assert.Equal(t, test.Collection2FileName3, sections[1].ImageSets[0].FileName)
}

func TestParseSection(t *testing.T) {
	metadata := test.Collection1

	section := parseSection(metadata)
	assert.Equal(t, test.Collection1["title"], section.Title)
	assert.Equal(t, test.Collection1["text"], string(section.Text))
	assert.Equal(t, test.Collection1["slug"], section.Slug)
	assert.Equal(t, test.Collection1["folder"], section.Folder)

	// ascending is false by default
	assert.Nil(t, metadata["ascending"])
	assert.Equal(t, false, section.IsAscending)

	metadata["ascending"] = true
	section = parseSection(metadata)
	assert.Equal(t, true, section.IsAscending)
}

func TestBuildImageSets(t *testing.T) {
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

	folder := test.Collection1["folder"].(string)

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

	set, _ := buildImageSet(test.Testfile, defaultOption)
	assert.Equal(t, filepath.Base(test.Testfile), set.FileName)
	assert.Equal(t, test.ThumbnailWidth, set.ThumbnailSize.Width)
	assert.Equal(t, test.ThumbnailHeight, set.ThumbnailSize.Height)
	assert.Equal(t, test.OriginalWidth, set.OriginalSize.Width)
	assert.Equal(t, test.OriginalHeight, set.OriginalSize.Height)
}
