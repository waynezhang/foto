package images

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"math"
	"path/filepath"
	"slices"
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
	return slices.Contains(
		[]string{".jpeg", ".jpg", ".webp", ".png"},
		strings.ToLower(filepath.Ext(path)))
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

func AspectedSize(size ImageSize, width int, minHeight int) ImageSize {
	ratio := float64(size.Height) / float64(size.Width)
	height := int(math.Round(float64(width) * ratio))
	if minHeight > height {
		height = minHeight
		width = int(math.Round(float64(height) / ratio))
	}

	return ImageSize{width, height}
}

func ResizeImage(src string, to string, width int, height int, compressQuality int) error {
	log.Debug().Msgf("Resizing %s to %dx%d", src, width, height)
	data, err := ResizeData(src, width, height, compressQuality)
	if err != nil {
		return err
	}

	if err := files.WriteDataToFile(data.Bytes(), to); err != nil {
		return err
	}
	return nil
}

func ResizeData(path string, width int, height int, compressQuality int) (*bytes.Buffer, error) {
	src, err := openImage(path)
	if err != nil {
		return nil, err
	}

	// If either width or height is 0, preserve aspect ratio
	// If both are specified, resize to exact dimensions
	resized := imaging.Resize(src, width, height, imaging.Lanczos)
	buf := new(bytes.Buffer)
	opt := jpeg.Options{Quality: compressQuality}
	if err := jpeg.Encode(buf, resized, &opt); err != nil {
		return nil, err
	}

	return buf, nil
}

func openImage(path string) (image.Image, error) {
	return imaging.Open(path, imaging.AutoOrientation(true))
}
