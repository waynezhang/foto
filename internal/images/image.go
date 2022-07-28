package images

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/log"
)

type ImageSize struct {
  Width int
  Height int
}

func IsPhotoSupported(path string) bool {
  lowerExt := strings.ToLower(filepath.Ext(path))
  return lowerExt == ".jpeg" || lowerExt == ".jpg"
}

func GetPhotoSize(path string) (*ImageSize, error) {
  f, err := os.Open(path)
  if err != nil {
    return nil, err
  }
	defer f.Close()

  img, _, err := image.DecodeConfig(f)
  if err != nil {
    return nil, err
  }

  return &ImageSize {
    img.Width,
    img.Height,
  }, nil
}

func ResizeImage(src string, to string, width int) error {
  log.Debug("Resizing %s to %d", src, width)
  data, err := ResizeData(src, width)
  if err != nil {
    return err
  }

  if err := files.WriteDataToFile(data.Bytes(), to); err != nil {
    return err
  }
  return nil
}

func ResizeData(path string, width int) (*bytes.Buffer, error) {
  src, err := imaging.Open(path)
  if err != nil {
    return nil, err
  }

	resized := imaging.Resize(src, width, 0, imaging.Lanczos)
  buf := new(bytes.Buffer)
  if err := jpeg.Encode(buf, resized, nil); err != nil {
    return nil, err
  }

  return buf, nil
}
