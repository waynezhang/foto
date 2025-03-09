package indexer

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/images"
)

type Section struct {
	Title     string
	Text      template.HTML
	Slug      string
	Folder    string
	Ascending bool
	ImageSets []ImageSet
}

type ImageSize images.ImageSize

type ImageSet struct {
	FileName        string
	ThumbnailSize   ImageSize
	OriginalSize    ImageSize
	CompressQuality int
}

func Build(metadata []config.SectionMetadata, option config.ExtractOption) ([]Section, error) {
	sections := []Section{}
	slugs := map[string]bool{}

	for _, val := range metadata {
		slug := val.Slug
		if !validSlug(slug) {
			return nil, fmt.Errorf("Slug \"%s\" is invalid. Only letters([a-zA-Z]), numbers([09-]), underscore(_) and hyphen(-) can be used.", slug)
		}
		if slugs[slug] {
			return nil, fmt.Errorf("Slug \"%s\" already exists. Slug needs to be unique.", slug)
		}

		log.Debug().Msgf("Extacting section [%s][/%s] %s", val.Title, val.Slug, val.Folder)

		sectionOption := sectionExtractOption(option, val)
		s := Section{
			Title:     val.Title,
			Text:      val.Text,
			Slug:      slug,
			Folder:    val.Folder,
			Ascending: val.Ascending,
			ImageSets: buildImageSets(val.Folder, val.Ascending, sectionOption),
		}
		slugs[slug] = true

		if len(s.ImageSets) > 0 {
			sections = append(sections, s)
		}
	}

	return sections, nil
}

func buildImageSets(folder string, ascending bool, option config.ExtractOption) []ImageSet {
	sets := []ImageSet{}

	wg := &sync.WaitGroup{}
	mutext := &sync.Mutex{}

	_ = filepath.WalkDir(folder, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			log.Warn().Msgf("Failed to extract info from %s (%v)", path, err)
			return nil
		}
		if info.IsDir() || !images.IsPhotoSupported(path) {
			return nil
		}

		wg.Add(1)

		go func(src string) {
			defer wg.Done()

			s, err := buildImageSet(src, option)
			if s != nil {
				mutext.Lock()
				sets = append(sets, *s)
				mutext.Unlock()
			} else {
				log.Warn().Msgf("Failed to extract info from %s (%v)", src, err)
			}
		}(path)

		return nil
	})
	wg.Wait()

	sort.SliceStable(sets, func(i, j int) bool {
		if ascending {
			return sets[i].FileName < sets[j].FileName
		} else {
			return sets[i].FileName > sets[j].FileName
		}
	})

	return sets
}

func buildImageSet(path string, option config.ExtractOption) (*ImageSet, error) {
	imageSize, err := images.GetPhotoSize(path)
	if err != nil {
		return nil, err
	}

	thumbnailWidth := option.ThumbnailWidth
	thumbnailHeight := option.ThumbnailHeight

	// If both width and height are specified, use them as is
	// If only one is specified, calculate the other to maintain aspect ratio
	if thumbnailWidth > 0 && thumbnailHeight == 0 {
		thumbnailHeight = images.AspectedHeight(*imageSize, thumbnailWidth)
	} else if thumbnailWidth == 0 && thumbnailHeight > 0 {
		thumbnailWidth = images.AspectedWidth(*imageSize, thumbnailHeight)
	} else if thumbnailWidth == 0 && thumbnailHeight == 0 {
		return nil, errors.New("Either thumbnailWidth or thumbnailHeight should be set")
	}

	// Handle original size
	originalWidth := option.OriginalWidth
	originalHeight := option.OriginalHeight

	// If both width and height are specified, use them as is
	// If only one is specified, calculate the other to maintain aspect ratio
	if originalWidth > 0 && originalHeight == 0 {
		originalHeight = images.AspectedHeight(*imageSize, originalWidth)
	} else if originalWidth == 0 && originalHeight > 0 {
		originalWidth = images.AspectedWidth(*imageSize, originalHeight)
	} else if originalWidth == 0 && originalHeight == 0 {
		return nil, errors.New("Either originalWidth or originalHeight should be set")
	}

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
		CompressQuality: option.CompressQuality,
	}, nil
}

func validSlug(slug string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9-_]+$", slug)
	return matched
}

func sectionExtractOption(global config.ExtractOption, metadata config.SectionMetadata) config.ExtractOption {
	sectionOption := global
	if metadata.ThumbnailWidth > 0 || metadata.ThumbnailHeight > 0 {
		sectionOption.ThumbnailWidth = metadata.ThumbnailWidth
		sectionOption.ThumbnailHeight = metadata.ThumbnailHeight
	}
	if metadata.OriginalWidth > 0 || metadata.OriginalHeight > 0 {
		sectionOption.OriginalWidth = metadata.OriginalWidth
		sectionOption.OriginalHeight = metadata.OriginalHeight
	}

	return sectionOption
}
