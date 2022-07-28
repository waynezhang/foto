package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
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
  spinner, _ := yacspin.New(yacspin.Config {
    Frequency:       100 * time.Millisecond,
    CharSet:         yacspin.CharSets[14],
    Suffix:          fmt.Sprintf(" exporting to %s", outputPath),
    SuffixAutoColon: true,
    StopMessage:     " succeed",
    StopCharacter:   "✓",
    StopColors:      []string{"fgGreen"},
  })
  _ = spinner.Start()

  spinnerMsg := func (format string, a ...interface{}) {
    spinner.Message(fmt.Sprintf(format, a...))
  }

  cfg := config.Shared()

  spinnerMsg("Removing directory %s", outputPath)
  err := files.PruneDirectory(outputPath)
  utils.CheckFatalError(err, "Failed to clean directory")

  photosDirectory := files.OutputPhotosFilePath(outputPath)
  section := exportPhotos(cfg, photosDirectory, func (path string) {
    spinnerMsg("%s", path)
  })

  indexPath := files.OutputIndexFilePath(outputPath)
  log.Debug("Exporting photos to %s", indexPath)
  generateIndex(cfg, section, indexPath)

  for _, folder := range cfg.OtherFolders() {
    targetFolder := filepath.Join(outputPath, folder)
    spinnerMsg("copying folder %s to %s", folder, targetFolder)
    
    if err := cp.Copy(folder, targetFolder); err != nil {
      log.Fatal("Failed to copy folder %s to %s (%s).", folder, targetFolder, err)
    }
  }

  _ = spinner.Stop()
} 

func exportPhotos(cfg config.Config, outputPath string, progressFunc images.ProgressFunc) []images.Section {
  if err := files.EnsureDirectory(outputPath); err != nil {
    return nil
  }

  return images.ExtractPhotos(cfg, &outputPath, progressFunc)
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
