package cmd

import (
	"html/template"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/log"
	"github.com/waynezhang/foto/internal/utils"
)

var outputPath string

var ExportCmd = func() *cobra.Command {
  cmd := &cobra.Command {
    Use: "export",
    Short: "Export sites",
    Run: export,
  }
  cmd.Flags().StringVarP(&outputPath, "output", "o", "dist", "Output directory")

  return cmd
}()

func export(cmd *cobra.Command, args []string) {
  log.Debug("Exprorting sites to %s...", outputPath)

  cfg := config.Shared()

  log.Debug("Removing directory %s", outputPath)
  err := files.PruneDirectory(outputPath)
  utils.CheckFatalError(err, "Failed to clean directory")

  photosDirectory := files.OutputPhotosFilePath(outputPath)
  log.Debug("Exporting photos to %s", photosDirectory)
  section := exportPhotos(cfg, photosDirectory)

  indexPath := files.OutputIndexFilePath(outputPath)
  log.Debug("Exporting photos to %s", indexPath)
  generateIndex(cfg, section, indexPath)

  assetsPath := files.OutputAssetsFilePath(outputPath)
  log.Debug("Copying assets to %s", assetsPath)
  err = copyAssets(assetsPath)
  utils.CheckFatalError(err, "Failed to copy assets files")

  for _, folder := range cfg.OtherFolders() {
    targetFolder := filepath.Join(outputPath, folder)
    log.Debug("Copying folder %s to %s", folder, targetFolder)
    
    if err := cp.Copy(folder, targetFolder); err != nil {
      log.Fatal("Failed to copy folder %s to %s (%s).", folder, targetFolder, err)
    }
  }
} 

func exportPhotos(cfg config.Config, outputPath string) []images.Section {
  if err := files.EnsureDirectory(outputPath); err != nil {
    return nil
  }

  return images.ExtractPhotos(cfg, &outputPath)
}

func generateIndex(cfg map[string]interface{}, sections []images.Section, path string) {
  f, err := os.Create(path)
  utils.CheckFatalError(err, "Failed to create index file.")
  defer f.Close()

  tmpl := template.Must(template.ParseFiles(files.TemplateFilePath))
  err = tmpl.Execute(f, struct {
    Config map[string]interface{}
    Sections []images.Section
  } {
    cfg,
    sections,
  })
  utils.CheckFatalError(err, "Failed to generate index page.")
} 

func copyAssets(to string) error {
  return cp.Copy(files.AssetsDir, to)
}
