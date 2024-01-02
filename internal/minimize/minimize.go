package minimize

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

type Minimizer interface {
	Minimizable(path string) bool
	MinimizeFile(src string, dest string) error
}

type NoneMinimizer struct{}

func (m NoneMinimizer) Minimizable(path string) bool {
	return true
}

func (m NoneMinimizer) MinimizeFile(src string, dest string) error {
	return nil
}

type MinifyMinimizer struct{}

func (m MinifyMinimizer) Minimizable(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".css" || ext == ".html" || ext == ".js"
}

func (m MinifyMinimizer) MinimizeFile(src string, dest string) error {
	var (
		buf bytes.Buffer
		err error
	)

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	switch m := minify.New(); filepath.Ext(src) {
	case ".css":
		err = css.Minify(m, &buf, f, nil)
	case ".html":
		err = html.Minify(m, &buf, f, nil)
	case ".js":
		err = js.Minify(m, &buf, f, nil)
	default:
		err = errors.New("Unsupported file")
	}
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, buf.Bytes(), 0644)
	return err
}
