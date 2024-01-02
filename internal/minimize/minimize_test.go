package minimize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/testdata"
)

func TestMinimizable(t *testing.T) {
	assert.True(t, Minimizable("/a/b/c/d.css"))
	assert.True(t, Minimizable("/a/b/c/d.js"))
	assert.True(t, Minimizable("/a/b/c/d.html"))

	assert.False(t, Minimizable("/a/b/c/d.ext"))
}

func TestMinimizeFile(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "foto-test")

	var err error

	err = MinimizeFile(testdata.TestHtmlFile, filepath.Join(tmp, "test.html"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.html")))

	err = MinimizeFile(testdata.TestCssFile, filepath.Join(tmp, "test.css"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.css")))

	err = MinimizeFile(testdata.TestJavascriptFile, filepath.Join(tmp, "test.js"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.js")))

	err = MinimizeFile(testdata.TestTxtFile, filepath.Join(tmp, "test.txt"))
	assert.NotNil(t, err)
	assert.False(t, files.IsExisting(filepath.Join(tmp, "test.txt")))

	err = MinimizeFile("nonexisting.html", filepath.Join(tmp, "nonexisting.html"))
	assert.NotNil(t, err)
	assert.False(t, files.IsExisting(filepath.Join(tmp, "nonexisting.html")))
}
