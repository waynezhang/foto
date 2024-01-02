package export

import (
	"fmt"
	"time"

	"github.com/theckman/yacspin"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/indexer"
	"github.com/waynezhang/foto/internal/log"
	mm "github.com/waynezhang/foto/internal/minimize"
	"github.com/waynezhang/foto/internal/utils"
)

type ProgressFunc func(path string)

type Context interface {
	cleanDirectory(outputPath string) error
	buildIndex(cfg config.Config) ([]indexer.Section, error)
	exportPhotos(sections []indexer.Section, outputPath string, cache cache.Cache, progressFunc ProgressFunc)
	generateIndexHtml(cfg config.Config, templatePath string, sections []indexer.Section, path string, minimizer mm.Minimizer)
	processOtherFolders(folders []string, outputPath string, minimizer mm.Minimizer, messageFunc func(src string, dst string))
}

func Export(outputPath string, minimize bool) error {
	return export(
		config.Shared(),
		outputPath,
		minimizer(minimize),
		cache.Shared(),
		new(DefaultExportContext),
	)
}

func export(cfg config.Config, outputPath string, minimizer mm.Minimizer, cache cache.Cache, ctx Context) error {
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
	section, err := ctx.buildIndex(cfg)
	if err != nil {
		ctx.cleanDirectory(outputPath)
		utils.CheckFatalError(err, "Failed t build index.")
	}

	ctx.exportPhotos(section, photosDirectory, cache, func(path string) {
		spinnerMsg("%s", path)
	})

	indexPath := files.OutputIndexFilePath(outputPath)
	log.Debug("Exporting photos to %s", indexPath)
	ctx.generateIndexHtml(cfg, constants.TemplateFilePath, section, indexPath, minimizer)

	msgFunc := func(src string, dst string) {
		spinnerMsg("copying folder %s to %s", src, dst)
	}
	ctx.processOtherFolders(cfg.GetOtherFolders(), outputPath, minimizer, msgFunc)

	_ = spinner.Stop()

	return nil
}

func minimizer(minimize bool) mm.Minimizer {
	if minimize {
		return mm.MinifyMinimizer{}
	} else {
		return mm.NoneMinimizer{}
	}
}
