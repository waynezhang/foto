package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := New("../../fs/static/foto.toml")

	site := cfg["site"].(map[string]interface{})
	assert.Equal(t, "A new site", site["title"])
	assert.Equal(t, "Author Here", site["author"])

	nav := site["nav"].([]interface{})
	assert.Equal(t, 3, len(nav))
	assert.Equal(t, "https://example.com", nav[0].(map[string]interface{})["link"])
	assert.Equal(t, "https://instagram.com/xxx", nav[1].(map[string]interface{})["link"])
	assert.Equal(t, "https://twitter.com/xxx", nav[2].(map[string]interface{})["link"])

	assert.Equal(t, 640, cfg.GetExtractOption().ThumbnailWidth)
	assert.Equal(t, 2048, cfg.GetExtractOption().OriginalWidth)

	layout := cfg["layout"].(map[string]interface{})
	assert.Equal(t, int64(1), layout["mincolumn"])
	assert.Equal(t, int64(4), layout["maxcolumn"])
	assert.Equal(t, int64(200), layout["minwidth"])

	sections := cfg["section"].([]interface{})
	assert.Equal(t, "Section 1", sections[0].(map[string]interface{})["title"])
	assert.Equal(t, "section-1", sections[0].(map[string]interface{})["slug"])
	assert.Equal(t, "~/photos/section-1", sections[0].(map[string]interface{})["folder"])
	assert.Equal(t, false, sections[0].(map[string]interface{})["ascending"])
	assert.Equal(t, "Section 2", sections[1].(map[string]interface{})["title"])
	assert.Equal(t, "section-2", sections[1].(map[string]interface{})["slug"])
	assert.Equal(t, "~/photos/section-2", sections[1].(map[string]interface{})["folder"])
	assert.Equal(t, false, sections[1].(map[string]interface{})["ascending"])

	others := cfg["others"].(map[string]interface{})
	assert.Equal(t, []string{"assets", "media"}, cfg.OtherFolders())
	assert.Equal(t, true, others["show_foto_footer"])

	// Test PhotoSwipe version
	assert.NotNil(t, cfg["PhotoSwipeVersion"])
}
