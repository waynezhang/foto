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

func Minimizable(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".css" || ext == ".html" || ext == ".js"
}

func MinimizeFile(src string, dest string) error {
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

