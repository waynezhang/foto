package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileConfig(t *testing.T) {
	cfg := NewFileConfig("../../fs/static/foto.toml")

	assert.Equal(t, 640, cfg.GetExtractOption().ThumbnailWidth)
	assert.Equal(t, 2048, cfg.GetExtractOption().OriginalWidth)

	sections := cfg.GetSectionMetadata()
	assert.Equal(t, "Section 1", sections[0].Title)
	assert.Equal(t, "section-1", sections[0].Slug)
	assert.Equal(t, "~/photos/section-1", sections[0].Folder)
	assert.Equal(t, false, sections[0].Ascending)
	assert.Equal(t, "Section 2", sections[1].Title)
	assert.Equal(t, "section-2", sections[1].Slug)
	assert.Equal(t, "~/photos/section-2", sections[1].Folder)
	assert.Equal(t, false, sections[1].Ascending)

	assert.Equal(t, []string{"assets", "media"}, cfg.GetOtherFolders())

	// Test PhotoSwipe version
	assert.NotNil(t, cfg.AllSettings()["photoswipeversion"])
}
