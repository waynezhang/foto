package files

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func WriteDataToFile(data []byte, to string) error {
	if err := EnsureParentDirectory(to); err != nil {
		return err
	}

	return os.WriteFile(to, data, 0644)
}

func Checksum(path string) (*string, error) {
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
