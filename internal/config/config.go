package config

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/log"
	"github.com/waynezhang/foto/internal/utils"
)

type Config map[string]interface{}

type ExtractOption struct {
	ThumbnailWidth int
	OriginalWidth  int
}

type SectionMetadata map[string]interface{}

var (
	once     sync.Once
	instance Config
)

func Shared() Config {
	once.Do(func() {
		instance = New("./foto.toml")
	})

	return instance
}

func New(file string) Config {
	v := viper.New()
	v.SetConfigFile(file)

	err := v.ReadInConfig()
	utils.CheckFatalError(err, "Failed to parse config file foto.toml")

	instance = loadConfig(v)

	instance["PhotoSwipeVersion"] = constants.PhotoSwipeVersion
	v.WatchConfig()

	return instance
}

func (cfg Config) GetSectionMetadata() []SectionMetadata {
	metadata := []SectionMetadata{}

	sections := cfg["section"].([]any)
	if sections == nil {
		return metadata
	}

	for _, v := range sections {
		metadata = append(metadata, (v.(map[string]interface{})))
	}
	return metadata
}

func (cfg Config) GetExtractOption() ExtractOption {
	imageOptions := cfg["image"].(map[string]interface{})
	return ExtractOption{
		ThumbnailWidth: int(imageOptions["thumbnailwidth"].(int64)),
		OriginalWidth:  int(imageOptions["originalwidth"].(int64)),
	}
}

func (cfg Config) OtherFolders() []string {
	others := cfg["others"]
	if others == nil {
		return nil
	}

	folders := others.(map[string]interface{})["folders"].([]interface{})
	ret := make([]string, len(folders))
	for i, v := range folders {
		ret[i] = v.(string)
	}
	return ret
}

func loadConfig(v *viper.Viper) Config {
	cfg := v.AllSettings()
	log.Debug("Config parsed: %v", cfg)
	return cfg
}
