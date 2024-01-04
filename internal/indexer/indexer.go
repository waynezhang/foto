package indexer

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/log"
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
	FileName      string
	ThumbnailSize ImageSize
	OriginalSize  ImageSize
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

		log.Debug("Extacting section [%s][/%s] %s", val.Title, val.Slug, val.Folder)
		s := Section{
			Title:     val.Title,
			Text:      val.Text,
			Slug:      slug,
			Folder:    val.Folder,
			Ascending: val.Ascending,
			ImageSets: buildImageSets(val.Folder, val.Ascending, option),
		}
		slugs[slug] = true

		sections = append(sections, s)
	}

	return sections, nil
}

func buildImageSets(folder string, ascending bool, option config.ExtractOption) []ImageSet {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal("Failed to get photos from %s (%v)", folder, err)
		return []ImageSet{}
	}

	ptrs := make([]*ImageSet, len(files))
	wg := &sync.WaitGroup{}
	for i, f := range files {
		path := filepath.Join(folder, f.Name())
		if f.IsDir() || !images.IsPhotoSupported(path) {
			continue
		}

		wg.Add(1)

		go func(idx int) {
			defer wg.Done()

			s, err := buildImageSet(path, option)
			if err != nil {
				log.Fatal("Failed to extract info from %s (%v)", path, err)
			}
			ptrs[idx] = s
		}(i)
	}
	wg.Wait()

	sets := []ImageSet{}
	for _, s := range ptrs {
		if s != nil {
			sets = append(sets, *s)
		}
	}

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

func validSlug(slug string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9-_]+$", slug)
	return matched
}
