package config

import (
	"html/template"
	"sync"

	"github.com/spf13/viper"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/log"
	"github.com/waynezhang/foto/internal/utils"
)

type Config interface {
	GetSectionMetadata() []SectionMetadata
	GetExtractOption() ExtractOption
	GetOtherFolders() []string
	AllSettings() map[string]interface{}
}

type ExtractOption struct {
	ThumbnailWidth int
	OriginalWidth  int
}

type SectionMetadata struct {
	Title     string
	Text      template.HTML
	Slug      string
	Folder    string
	Ascending bool
}

var (
	once     sync.Once
	instance Config
)

func Shared() Config {
	once.Do(func() {
		instance = NewFileConfig("./foto.toml")
	})

	return instance
}

// File Config

type FileConfig struct {
	v            *viper.Viper
	option       ExtractOption
	sections     []SectionMetadata
	otherFolders []string
}

func NewFileConfig(file string) Config {
	v := viper.New()
	v.SetConfigFile(file)

	err := v.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse config file foto.toml")

	// Inject PhotoSwipeVersion
	v.Set("PhotoSwipeVersion", constants.PhotoSwipeVersion)

	config := FileConfig{v: v}

	v.UnmarshalKey("section", &config.sections)
	v.UnmarshalKey("image", &config.option)
	v.UnmarshalKey("others.folders", &config.otherFolders)

	log.Debug("Config parsed: %v", config)

	return config
}

func (cfg FileConfig) GetSectionMetadata() []SectionMetadata {
	return cfg.sections
}

func (cfg FileConfig) GetExtractOption() ExtractOption {
	return cfg.option
}

func (cfg FileConfig) GetOtherFolders() []string {
	return cfg.otherFolders
}

func (cfg FileConfig) AllSettings() map[string]interface{} {
	return cfg.v.AllSettings()
}
