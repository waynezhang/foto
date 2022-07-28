package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/images"
	"github.com/waynezhang/foto/internal/log"
	"github.com/waynezhang/foto/internal/utils"
	"github.com/waynezhang/foto/internal/url"
)

var port = 5000

var PreviewCmd = func() *cobra.Command {
  cmd := &cobra.Command {
    Use: "preview",
    Short: "Preview in local environment",
    Run: preview,
  }
  cmd.Flags().IntVarP(&port, "port", "p", 5000, "Port")

  return cmd
}()

func preview(cmd *cobra.Command, args []string) {
  log.Debug("Creating Preview...")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    handleRoot(
      config.Shared(),
      images.ExtractPhotos(config.Shared(), nil, nil),
      w,
      r,
    )
  })

  http.HandleFunc(url.PhotosPath, func(w http.ResponseWriter, r *http.Request) {
    handleImage(
      strings.TrimPrefix(r.URL.Path, url.PhotosPath),
      config.Shared(),
      images.ExtractPhotos(config.Shared(), nil, nil),
      w,
      r,
    )
  })

  otherFolders := config.Shared().OtherFolders()
  for _, folder := range otherFolders {
    dir := http.FileServer(http.Dir(folder))
    path := "/" + folder + "/"
    http.Handle(path, http.StripPrefix(path, dir))
  }

  log.Info("Server started -> http://localhost:%d", port)
  err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
  utils.CheckFatalError(err, "Failed to listen the port")
}

func handleRoot(cfg map[string]interface{}, sections []images.Section, w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles(files.TemplateFilePath))
  _ = tmpl.Execute(w, struct {
    Config map[string]interface{}
    Sections []images.Section
  } {
    cfg,
    sections,
  })
} 

func handleImage(path string, cfg map[string]interface{}, sections []images.Section, w http.ResponseWriter, r *http.Request) {
  comps := strings.Split(path, "/")
  if len(comps) != 3 {
    http.NotFound(w, r)
    return
  }

  var dir *string
  for _, s := range sections {
    if s.Slug == comps[0] {
      dir = &s.Folder
      break
    }
  }
  
  if dir == nil {
    http.NotFound(w, r)
    return
  }

  key := comps[1] + "width"
  width := cfg["image"].(map[string]interface{})[key].(int64)
  data, err := images.ResizeData(*dir + "/" + comps[2], int(width))
  if err != nil {
    http.NotFound(w, r)
    return
  }

  w.Header().Set("Content-Type", "image/jpeg")
  w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
  _, _ = w.Write(data.Bytes())
} 
