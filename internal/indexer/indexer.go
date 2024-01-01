package indexer

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/log"
)

type Section struct {
	Title     string
	Text      string
	Slug      string
	Folder    string
	Ascending bool
	ImageSets []ImageSet
}

type ImageSize images.ImageSize

type ImageSet struct {
	FileName      string
	ThumbnailSize ImageSize
	OriginalSize  ImageSize
}

func Build(metadata []config.SectionMetadata, option config.ExtractOption) []Section {
	sections := []Section{}
	for _, val := range metadata {
		s := Section{
			Title:     val.Title,
			Text:      val.Text,
			Slug:      val.Slug,
			Folder:    val.Folder,
			Ascending: val.Ascending,
			ImageSets: buildImageSets(val.Folder, val.Ascending, option),
		}
		log.Debug("Extacting section [%s][/%s] %s", val.Title, val.Slug, val.Folder)

		sections = append(sections, s)
	}

	return sections
}

func buildImageSets(folder string, ascending bool, option config.ExtractOption) []ImageSet {
	imageSet := []ImageSet{}
	err := filepath.WalkDir(folder, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !images.IsPhotoSupported(path) {
			return nil
		}

		imgSet, err := buildImageSet(path, option)
		if err != nil {
			return err
		}
		imageSet = append(imageSet, *imgSet)

		return nil
	})
	if err != nil {
		log.Fatal("Failed to get photos from %s (%v)", folder, err)
	}

	sort.SliceStable(imageSet, func(i, j int) bool {
		if ascending {
			return imageSet[i].FileName < imageSet[j].FileName
		} else {
			return imageSet[i].FileName > imageSet[j].FileName
		}
	})

	return imageSet
}

func buildImageSet(path string, option config.ExtractOption) (*ImageSet, error) {
	imageSize, err := images.GetPhotoSize(path)
	if err != nil {
		return nil, err
	}

	thumbnailWidth := option.ThumbnailWidth
	thumbnailHeight := images.AspectedHeight(*imageSize, thumbnailWidth)

	originalWidth := option.OriginalWidth
	originalHeight := images.AspectedHeight(*imageSize, originalWidth)

	return &ImageSet{
		FileName: filepath.Base(path),
		ThumbnailSize: ImageSize{
			thumbnailWidth,
			thumbnailHeight,
		},
		OriginalSize: ImageSize{
			originalWidth,
			originalHeight,
		},
	}, nil
}
