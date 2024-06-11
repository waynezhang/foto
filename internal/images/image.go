package images

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog/log"
	"github.com/waynezhang/foto/internal/files"
	_ "golang.org/x/image/webp"
)

type ImageSize struct {
	Width  int
	Height int
}

func IsPhotoSupported(path string) bool {
	lowerExt := strings.ToLower(filepath.Ext(path))
	return lowerExt == ".jpeg" || lowerExt == ".jpg" || lowerExt == ".webp"
}

func GetPhotoSize(path string) (*ImageSize, error) {
	src, err := openImage(path)
	if err != nil {
		return nil, err
	}

	return &ImageSize{
		src.Bounds().Size().X,
		src.Bounds().Size().Y,
	}, nil
}

func AspectedHeight(size ImageSize, width int) int {
	ratio := float32(size.Height) / float32(size.Width)
	return int(float32(width) * ratio)
}

func ResizeImage(src string, to string, width int) error {
	log.Debug().Msgf("Resizing %s to %d", src, width)
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
	src, err := openImage(path)
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

func openImage(path string) (image.Image, error) {
	return imaging.Open(path, imaging.AutoOrientation(true))
}
