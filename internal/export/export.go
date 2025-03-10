package export

import (
	"fmt"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/rs/zerolog/log"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/indexer"
	mm "github.com/waynezhang/foto/internal/minimize"
	"github.com/waynezhang/foto/internal/utils"
)

func Export(outputPath string, minimize bool) {
	export(
		config.Shared(),
		outputPath,
		minimizer(minimize),
		cache.Shared(),
		new(defaultExportContext),
	)
}

type progressFunc func(path string)

type context interface {
	cleanDirectory(outputPath string) error
	buildIndex(cfg config.Config) ([]indexer.Section, error)
	exportPhotos(
		sections []indexer.Section,
		outputPath string,
		cache cache.Cache,
		postProgressFn progressFunc,
	)
	generateIndexHtml(
		cfg config.Config,
		templatePath string,
		sections []indexer.Section,
		path string,
		minimizer mm.Minimizer,
	)
	processOtherFolders(
		folders []string,
		outputPath string,
		minimizer mm.Minimizer,
		messageFunc func(src string, dst string),
	)
}

func export(
	cfg config.Config,
	outputPath string,
	minimizer mm.Minimizer,
	cache cache.Cache,
	ctx context,
) {
	sm := ysmrr.NewSpinnerManager(
		ysmrr.WithAnimation(animations.Dots),
	)
	prefixSpinnerMsg := fmt.Sprintf("exporting to %s: ", outputPath)
	spinner := sm.AddSpinner(prefixSpinnerMsg)
	sm.Start()

	spinnerMsg := func(format string, a ...any) {
		spinner.UpdateMessagef(prefixSpinnerMsg+format, a...)
	}

	spinnerMsg("removing directory %s", outputPath)
	err := ctx.cleanDirectory(outputPath)
	if err != nil {
		utils.CheckFatalError(err, "Failed to remove directory.")
	}

	spinnerMsg("building index")
	photosDirectory := files.OutputPhotosFilePath(outputPath)
	section, err := ctx.buildIndex(cfg)
	if err != nil {
		_ = ctx.cleanDirectory(outputPath)
		utils.CheckFatalError(err, "Failed to build index.")
	}

	ctx.exportPhotos(section, photosDirectory, cache, func(path string) {
		spinnerMsg("processed image %s", path)
	})

	indexPath := files.OutputIndexFilePath(outputPath)
	log.Debug().Msgf("Exporting photos to %s", indexPath)
	ctx.generateIndexHtml(cfg, constants.TemplateFilePath, section, indexPath, minimizer)

	ctx.processOtherFolders(cfg.GetOtherFolders(), outputPath, minimizer, func(src string, dst string) {
		spinnerMsg("copying folder %s to %s", src, dst)
	})

	spinner.UpdateMessage(prefixSpinnerMsg + "succeeded")

	spinner.Complete()
	sm.Stop()
}

func minimizer(minimize bool) mm.Minimizer {
	if minimize {
		return mm.MinifyMinimizer{}
	} else {
		return mm.NoneMinimizer{}
	}
}
