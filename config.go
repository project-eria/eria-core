package eria

import (
	configmanager "github.com/project-eria/eria-core/config-manager"
	zlog "github.com/rs/zerolog/log"
)

var eriaConfig = struct {
	Host        string `yaml:"host,omitempty"`
	Port        uint   `yaml:"port" default:"80"`
	Location    string `yaml:"location" required:"true"`
	ExposedAddr string `yaml:"exposedAddr,omitempty"`
}{}

// LoadConfig Loads the config file into a struct
func LoadConfig(config interface{}) *configmanager.ConfigManager {
	cm, err := configmanager.Init(*_configPath, config, &eriaConfig)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			zlog.Fatal().Msg("[core:LoadConfig] Config file do not exists...")
		} else {
			zlog.Fatal().Str("filePath", *_configPath).Err(err).Msg("[core:LoadConfig]")
		}
	}

	if err := cm.Load(); err != nil {
		zlog.Fatal().Err(err).Msg("[core:LoadConfig]")
	}
	return cm
}
