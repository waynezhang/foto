package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteData(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmp)

	testdata := []byte("this is the content")
	file := filepath.Join(tmp, "sub-dir", "testfile")
	err = WriteDataToFile(testdata, file)
	assert.Nil(t, err)

	bytes, err := os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, testdata, bytes)
}
