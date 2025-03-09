package config

import (
	"html/template"
	"sync"
)

type Config interface {
	GetSectionMetadata() []SectionMetadata
	GetExtractOption() ExtractOption
	GetOtherFolders() []string
	AllSettings() map[string]interface{}
}

type ExtractOption struct {
	ThumbnailWidth  int
	ThumbnailHeight int
	OriginalWidth   int
	OriginalHeight  int
	CompressQuality int
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
