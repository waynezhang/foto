package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/waynezhang/foto/internal/log"
)

func CheckFatalError(err error, errMessage string) {
	if err == nil {
		return
	}

	log.Fatal("%s (%s)", errMessage, err)
	os.Exit(1)
}

func FileChecksum(path string) (*string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil, nil
	}

	value := hex.EncodeToString(hasher.Sum(nil))
	return &value, nil
}
