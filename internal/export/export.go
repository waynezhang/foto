package export

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	cp "github.com/otiai10/copy"
	"github.com/theckman/yacspin"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/indexer"
	"github.com/waynezhang/foto/internal/log"
	mm "github.com/waynezhang/foto/internal/minimize"
	"github.com/waynezhang/foto/internal/utils"
)

type ProgressFunc func(path string)

type Context interface {
	cleanDirectory(outputPath string) error
	buildIndex(cfg config.Config) []indexer.Section
	exportPhotos(sections []indexer.Section, outputPath string, cache cache.Cache, progressFunc ProgressFunc)
	generateIndexHtml(cfg config.Config, sections []indexer.Section, path string, minimize bool)
	processOtherFolders(folders []string, outputPath string, minimize bool, messageFunc func(src string, dst string))
}

func Export(outputPath string, minimize bool) error {
	return export(config.Shared(), outputPath, minimize, cache.Shared(), new(DefaultExportContext))
}

func export(cfg config.Config, outputPath string, minimize bool, cache cache.Cache, ctx Context) error {
	spinner, _ := yacspin.New(yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[14],
		Suffix:          fmt.Sprintf(" exporting to %s", outputPath),
		SuffixAutoColon: true,
		StopMessage:     " succeed",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	})
	_ = spinner.Start()

	spinnerMsg := func(format string, a ...interface{}) {
		spinner.Message(fmt.Sprintf(format, a...))
	}

	spinnerMsg("Removing directory %s", outputPath)
	err := ctx.cleanDirectory(outputPath)
	if err != nil {
		return err
	}

	spinnerMsg("Building index %s", outputPath)
	photosDirectory := files.OutputPhotosFilePath(outputPath)
	section := ctx.buildIndex(cfg)
	ctx.exportPhotos(section, photosDirectory, cache, func(path string) {
		spinnerMsg("%s", path)
	})

	indexPath := files.OutputIndexFilePath(outputPath)
	log.Debug("Exporting photos to %s", indexPath)
	ctx.generateIndexHtml(cfg, section, indexPath, minimize)

	msgFunc := func(src string, dst string) {
		spinnerMsg("copying folder %s to %s", src, dst)
	}
	ctx.processOtherFolders(cfg.GetOtherFolders(), outputPath, minimize, msgFunc)

	_ = spinner.Stop()

	return nil
}

type DefaultExportContext struct{}

func (ctx DefaultExportContext) cleanDirectory(outputPath string) error {
	return files.PruneDirectory(outputPath)
}

func (ctx DefaultExportContext) buildIndex(cfg config.Config) []indexer.Section {
	return indexer.Build(cfg.GetSectionMetadata(), cfg.GetExtractOption())
}

func (ctx DefaultExportContext) exportPhotos(sections []indexer.Section, outputPath string, cache cache.Cache, progressFunc ProgressFunc) {
	if err := files.EnsureDirectory(outputPath); err != nil {
		utils.CheckFatalError(err, "Failed to prepare output directory")
		return
	}

	for _, s := range sections {
		for _, set := range s.ImageSets {
			srcPath := filepath.Join(s.Folder, set.FileName)

			log.Debug("Processing image %s", srcPath)
			if progressFunc != nil {
				progressFunc(srcPath)
			}

			thumbnailPath := files.OutputPhotoThumbnailFilePath(outputPath, s.Slug, srcPath)
			err := resizeImageAndCache(srcPath, thumbnailPath, set.ThumbnailSize.Width, cache)
			utils.CheckFatalError(err, "Failed to generate thumbnail image")

			originalPath := files.OutputPhotoOriginalFilePath(outputPath, s.Slug, srcPath)
			err = resizeImageAndCache(srcPath, originalPath, set.OriginalSize.Width, cache)
			utils.CheckFatalError(err, "Failed to generate original image")
		}
	}
}

func (ctx DefaultExportContext) generateIndexHtml(cfg config.Config, sections []indexer.Section, path string, minimize bool) {
	f, err := os.Create(path)
	utils.CheckFatalError(err, "Failed to create index file.")
	defer f.Close()

	tmpl := template.Must(template.ParseFiles(constants.TemplateFilePath))
	err = tmpl.Execute(f, struct {
		Config   map[string]interface{}
		Sections []indexer.Section
	}{
		cfg.AllSettings(),
		sections,
	})
	utils.CheckFatalError(err, "Failed to generate index page.")

	if minimize {
		_ = mm.MinimizeFile(path, path)
	}
}

func (ctx DefaultExportContext) processOtherFolders(folders []string, outputPath string, minimize bool, messageFunc func(src string, dst string)) {
	for _, folder := range folders {
		targetFolder := filepath.Join(outputPath, folder)
		if messageFunc != nil {
			messageFunc(folder, targetFolder)
		}

		if err := cp.Copy(folder, targetFolder); err != nil {
			log.Fatal("Failed to copy folder %s to %s (%s).", folder, targetFolder, err)
		}
		if minimize {
			_ = filepath.WalkDir(targetFolder, func(path string, d fs.DirEntry, err error) error {
				if mm.Minimizable(path) {
					return mm.MinimizeFile(path, path)
				}
				return nil
			})
		}
	}
}

func resizeImageAndCache(src string, to string, width int, cache cache.Cache) error {
	cached := cache.CachedImage(src, width)
	if cached != nil {
		log.Debug("Found cached image for %s", src)
		err := cp.Copy(*cached, to)
		if err == nil {
			return nil
		}
	}

	err := images.ResizeImage(src, to, width)
	if err != nil {
		return err
	}

	cache.AddImage(src, width, to)

	return nil
}
