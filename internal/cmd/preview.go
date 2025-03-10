package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/indexer"
	"github.com/waynezhang/foto/internal/utils"
)

var port = 5000

var PreviewCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview in local environment",
		Run:   preview,
	}
	cmd.Flags().IntVarP(&port, "port", "p", 5000, "Port")

	return cmd
}()

func preview(cmd *cobra.Command, args []string) {
	log.Debug().Msg("Creating Preview...")

	config := config.Shared()
	index, err := indexer.Build(config.GetSectionMetadata(), config.GetExtractOption())
	utils.CheckFatalError(err, "Failed to build index")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRoot(
			config,
			index,
			w,
			r,
		)
	})

	http.HandleFunc(constants.PhotosURLPath, func(w http.ResponseWriter, r *http.Request) {
		handleImage(
			strings.TrimPrefix(r.URL.Path, constants.PhotosURLPath),
			config,
			index,
			w,
			r,
		)
	})

	otherFolders := config.GetOtherFolders()
	for _, folder := range otherFolders {
		dir := http.FileServer(http.Dir(folder))
		path := "/" + folder + "/"
		http.Handle(path, http.StripPrefix(path, dir))
	}

	log.Info().Msgf("Server started -> http://localhost:%d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	utils.CheckFatalError(err, "Failed to listen the port")
}

func handleRoot(cfg config.Config, sections []indexer.Section, w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(constants.TemplateFilePath))
	_ = tmpl.Execute(w, struct {
		Config   map[string]interface{}
		Sections []indexer.Section
	}{
		cfg.AllSettings(),
		sections,
	})
}

func handleImage(path string, cfg config.Config, sections []indexer.Section, w http.ResponseWriter, r *http.Request) {
	comps := strings.Split(path, "/")
	if len(comps) != 3 {
		http.NotFound(w, r)
		return
	}

	slug := comps[0]
	key := comps[1]
	file := comps[2]

	var file_path string
	var size images.ImageSize
	for _, s := range sections {
		if s.Slug == slug {
			for _, is := range s.ImageSets {
				if is.FileName == file {
					file_path = filepath.Join(s.Folder, file)
					if key == "thumbnail" {
						size = is.ThumbnailSize
					} else if key == "original" {
						size = is.OriginalSize
					}
					break
				}
			}
		}
	}

	if file_path == "" || size.Width == 0 || size.Height == 0 {
		http.NotFound(w, r)
		return
	}

	data, err := images.ResizeData(file_path, size.Width, size.Height, cfg.GetExtractOption().CompressQuality)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	_, _ = w.Write(data.Bytes())
}
