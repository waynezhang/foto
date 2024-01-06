package cmd

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	staticFs "github.com/waynezhang/foto/fs"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/utils"
)

var CreateCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [directory]",
		Short: "Create a new site",
		Run:   create,
	}

	return cmd
}()

func create(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal().Msg("Directory argument not found")
	}

	targetPath := args[0]
	if files.IsExisting(targetPath) {
		log.Fatal().Msg("Directory " + targetPath + " already exists.")
	}

	log.Debug().Msgf("Creating directory %s...", targetPath)
	err := files.EnsureDirectory(targetPath)
	utils.CheckFatalError(err, "Failed to create directory")

	err = extractFiles(targetPath)
	utils.CheckFatalError(err, "Failed to extract files")
}

func extractFiles(to string) error {
	return fs.WalkDir(staticFs.FS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		data, err := staticFs.FS.ReadFile(path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(to, strings.TrimPrefix(path, "static/"))
		return files.WriteDataToFile(data, targetPath)
	})
}
