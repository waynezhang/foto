package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/utils"
)

type fileConfig struct {
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

	config := fileConfig{v: v}

	v.UnmarshalKey("section", &config.sections)
	v.UnmarshalKey("image", &config.option)
	v.UnmarshalKey("others.folders", &config.otherFolders)

	log.Debug().Msgf("Config parsed: %v", config)

	return config
}

func (cfg fileConfig) GetSectionMetadata() []SectionMetadata {
	return cfg.sections
}

func (cfg fileConfig) GetExtractOption() ExtractOption {
	return cfg.option
}

func (cfg fileConfig) GetOtherFolders() []string {
	return cfg.otherFolders
}

func (cfg fileConfig) AllSettings() map[string]interface{} {
	return cfg.v.AllSettings()
}
