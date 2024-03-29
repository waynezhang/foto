package minimize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/testdata"
)

func TestMinifyMinimizable(t *testing.T) {
	minimizer := MinifyMinimizer{}
	assert.True(t, minimizer.Minimizable("/a/b/c/d.css"))
	assert.True(t, minimizer.Minimizable("/a/b/c/d.js"))
	assert.True(t, minimizer.Minimizable("/a/b/c/d.html"))

	assert.False(t, minimizer.Minimizable("/a/b/c/d.ext"))
}

func TestMinifyMinimizeFile(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "foto-test")

	minimizer := MinifyMinimizer{}

	var err error

	err = minimizer.MinimizeFile(testdata.TestHtmlFile, filepath.Join(tmp, "test.html"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.html")))

	err = minimizer.MinimizeFile(testdata.TestCssFile, filepath.Join(tmp, "test.css"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.css")))

	err = minimizer.MinimizeFile(testdata.TestJavascriptFile, filepath.Join(tmp, "test.js"))
	assert.Nil(t, err)
	assert.True(t, files.IsExisting(filepath.Join(tmp, "test.js")))

	err = minimizer.MinimizeFile(testdata.TestTxtFile, filepath.Join(tmp, "test.txt"))
	assert.NotNil(t, err)
	assert.False(t, files.IsExisting(filepath.Join(tmp, "test.txt")))

	err = minimizer.MinimizeFile("nonexisting.html", filepath.Join(tmp, "nonexisting.html"))
	assert.NotNil(t, err)
	assert.False(t, files.IsExisting(filepath.Join(tmp, "nonexisting.html")))
}
