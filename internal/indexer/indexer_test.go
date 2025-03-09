package indexer

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/testdata"
)

var (
	defaultOption = config.ExtractOption{
		ThumbnailWidth:  testdata.ThumbnailWidth,
		OriginalWidth:   testdata.OriginalWidth,
		CompressQuality: testdata.CompressQuality,
	}
)

func TestValidSlug(t *testing.T) {
	assert.True(t, validSlug("abcde-efg_9999"))
	assert.False(t, validSlug("abcde efg_9999"))
	assert.False(t, validSlug("abcde-efgかたかな"))
	assert.False(t, validSlug("abcde.efg_999"))
	assert.False(t, validSlug(""))
}

func TestBuild(t *testing.T) {
	var meta1 config.SectionMetadata
	var meta2 config.SectionMetadata
	_ = mapstructure.Decode(testdata.Collection1, &meta1)
	_ = mapstructure.Decode(testdata.Collection2, &meta2)

	data := []config.SectionMetadata{meta1, meta2}

	sections, _ := Build(data, defaultOption)
	assert.Equal(t, 2, len(sections))
	assert.Equal(t, testdata.Collection1["title"], sections[0].Title)

	assert.Equal(t, template.HTML("This is Section 1"), sections[0].Text)

	assert.Equal(t, 3, len(sections[0].ImageSets))
	assert.Equal(t, testdata.Collection1FileName1, sections[0].ImageSets[0].FileName)
	assert.Equal(t, 3, len(sections[1].ImageSets))
	assert.Equal(t, testdata.Collection2FileName1, sections[1].ImageSets[0].FileName)
}

func TestBuildSizeOverride(t *testing.T) {
	var meta1 config.SectionMetadata
	var meta2 config.SectionMetadata
	_ = mapstructure.Decode(testdata.Collection1, &meta1)
	_ = mapstructure.Decode(testdata.Collection2, &meta2)

	data := []config.SectionMetadata{meta1, meta2}

	sections, _ := Build(data, defaultOption)
	assert.Equal(t, 640, sections[0].ImageSets[0].ThumbnailSize.Width)
	assert.Equal(t, 480, sections[0].ImageSets[0].ThumbnailSize.Height)
	assert.Equal(t, 2048, sections[0].ImageSets[0].OriginalSize.Width)
	assert.Equal(t, 1536, sections[0].ImageSets[0].OriginalSize.Height)

	assert.Equal(t, 151, sections[1].ImageSets[0].ThumbnailSize.Width)
	assert.Equal(t, 100, sections[1].ImageSets[0].ThumbnailSize.Height)
	assert.Equal(t, 605, sections[1].ImageSets[0].OriginalSize.Width)
	assert.Equal(t, 400, sections[1].ImageSets[0].OriginalSize.Height)
}

func TestBuildDuplicatedSlugs(t *testing.T) {
	var meta1 config.SectionMetadata
	var meta2 config.SectionMetadata
	_ = mapstructure.Decode(testdata.Collection1, &meta1)
	_ = mapstructure.Decode(testdata.Collection2, &meta2)
	meta2.Slug = meta1.Slug

	data := []config.SectionMetadata{meta1, meta2}

	_, err := Build(data, defaultOption)
	assert.NotNil(t, err)
}

func TestBuildEmptySection(t *testing.T) {
	var meta config.SectionMetadata
	var emptyMeta config.SectionMetadata
	_ = mapstructure.Decode(testdata.Collection1, &meta)
	_ = mapstructure.Decode(testdata.EmptyCollection, &emptyMeta)

	data := []config.SectionMetadata{meta, emptyMeta}

	sections, _ := Build(data, defaultOption)
	assert.Equal(t, 1, len(sections))
	assert.Equal(t, testdata.Collection1["title"], sections[0].Title)
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

func TestInvalidBuildImageSets(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "foto-test")
	path := filepath.Join(tmp, "folder-not-exist")
	// no crash expected
	_ = buildImageSets(path, true, defaultOption)
}

func TestBuildImageSet(t *testing.T) {
	set, _ := buildImageSet(testdata.Testfile, defaultOption)
	assert.Equal(t, filepath.Base(testdata.Testfile), set.FileName)
	assert.Equal(t, testdata.ThumbnailWidth, set.ThumbnailSize.Width)
	assert.Equal(t, testdata.ThumbnailHeight, set.ThumbnailSize.Height)
	assert.Equal(t, testdata.OriginalWidth, set.OriginalSize.Width)
	assert.Equal(t, testdata.OriginalHeight, set.OriginalSize.Height)
	assert.Equal(t, testdata.CompressQuality, set.CompressQuality)
}

func TestSectionExtractOption(t *testing.T) {
	testCases := []struct {
		name           string
		global         config.ExtractOption
		sectionMeta    config.SectionMetadata
		expectedOption config.ExtractOption
	}{
		{
			name: "no section overrides",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalWidth:  2048,
			},
			sectionMeta: config.SectionMetadata{},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  640,
				ThumbnailHeight: 0,
				OriginalWidth:   2048,
				OriginalHeight:  0,
			},
		},
		{
			name: "section override thumbnail",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalWidth:  2048,
			},
			sectionMeta: config.SectionMetadata{
				ThumbnailWidth:  800,
				ThumbnailHeight: 0,
				OriginalWidth:   2048,
				OriginalHeight:  0,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth: 800,
				OriginalWidth:  2048,
			},
		},
		{
			name: "section override thumbnail height",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalWidth:  2048,
			},
			sectionMeta: config.SectionMetadata{
				ThumbnailHeight: 600,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  0,
				ThumbnailHeight: 600,
				OriginalWidth:   2048,
				OriginalHeight:  0,
			},
		},
		{
			name: "section override original",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalWidth:  2048,
			},
			sectionMeta: config.SectionMetadata{
				OriginalWidth: 1920,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  640,
				ThumbnailHeight: 0,
				OriginalWidth:   1920,
				OriginalHeight:  0,
			},
		},
		{
			name: "section override original height",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalHeight: 1536,
			},
			sectionMeta: config.SectionMetadata{
				OriginalHeight: 1080,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  640,
				ThumbnailHeight: 0,
				OriginalWidth:   0,
				OriginalHeight:  1080,
			},
		},
		{
			name: "section override all dimensions",
			global: config.ExtractOption{
				ThumbnailWidth: 640,
				OriginalWidth:  2048,
			},
			sectionMeta: config.SectionMetadata{
				ThumbnailWidth:  800,
				ThumbnailHeight: 600,
				OriginalWidth:   1920,
				OriginalHeight:  1080,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  800,
				ThumbnailHeight: 600,
				OriginalWidth:   1920,
				OriginalHeight:  1080,
			},
		},
		{
			name: "zero values for section should be ignored",
			global: config.ExtractOption{
				ThumbnailWidth:  640,
				ThumbnailHeight: 480,
				OriginalWidth:   2048,
				OriginalHeight:  1536,
			},
			sectionMeta: config.SectionMetadata{
				ThumbnailWidth:  0,
				ThumbnailHeight: 0,
				OriginalWidth:   0,
				OriginalHeight:  0,
			},
			expectedOption: config.ExtractOption{
				ThumbnailWidth:  640,
				ThumbnailHeight: 480,
				OriginalWidth:   2048,
				OriginalHeight:  1536,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sectionExtractOption(tc.global, tc.sectionMeta)

			assert.Equal(t, tc.expectedOption.ThumbnailWidth, result.ThumbnailWidth, "ThumbnailWidth should match expected value")
			assert.Equal(t, tc.expectedOption.ThumbnailHeight, result.ThumbnailHeight, "ThumbnailHeight should match expected value")
			assert.Equal(t, tc.expectedOption.OriginalWidth, result.OriginalWidth, "OriginalWidth should match expected value")
			assert.Equal(t, tc.expectedOption.OriginalHeight, result.OriginalHeight, "OriginalHeight should match expected value")
			assert.Equal(t, tc.expectedOption.CompressQuality, result.CompressQuality, "CompressQuality should match expected value")
		})
	}
}
