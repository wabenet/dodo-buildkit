package config

import (
	"fmt"

	"github.com/dodo/dodo-build/pkg/config/decoder"
	"github.com/dodo/dodo-build/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/sahilm/fuzzy"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoadBackdrop(backdrop string) (*types.Backdrop, error) {
	config := loadConfig()
	if result, ok := config.Backdrops[backdrop]; ok {
		return &result, nil
	}

	matches := fuzzy.Find(backdrop, config.Names())
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", backdrop)
	}
	return nil, fmt.Errorf("backdrop '%s' not found, did you mean '%s'?", backdrop, matches[0].Str)
}

func LoadImage(image string) (*types.BuildInfo, error) {
	config := loadConfig()
	for _, backdrop := range config.Backdrops {
		if backdrop.Build != nil && backdrop.Build.ImageName == image {
			return backdrop.Build, nil
		}
	}
	return nil, fmt.Errorf("could not find any backdrop configuration that would produce image '%s'", image)
}

func loadConfig() *decoder.Group {
	var result decoder.Group
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			var mapType map[interface{}]interface{}
			if err := yaml.Unmarshal(configFile.Content, &mapType); err != nil {
				log.WithFields(log.Fields{"file": configFile.Path}).Warn("invalid YAML syntax in file")
				return false
			}

			decoder := decoder.NewDecoder(configFile.Path)
			config, err := decoder.DecodeGroup(configFile.Path, mapType)
			if err != nil {
				log.WithFields(log.Fields{"file": configFile.Path, "reason": err}).Warn("invalid config file")
				return false
			}

			result.Merge(&config)
			return false
		},
	})
	return &result
}