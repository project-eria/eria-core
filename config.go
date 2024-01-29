package eria

import (
	"github.com/project-eria/eria-core/automations"
	configmanager "github.com/project-eria/eria-core/config-manager"
	zlog "github.com/rs/zerolog/log"
)

var eriaConfig = struct {
	Host         string                         `yaml:"host,omitempty"`
	Port         uint                           `yaml:"port" default:"80"`
	Location     string                         `yaml:"location" required:"true"`
	ExposedAddr  string                         `yaml:"exposedAddr,omitempty"`
	Automations  []automations.AutomationConfig `yaml:"automations,omitempty"`
	ContextsRef  string                         `yaml:"contextsRef,omitempty"`
	RemoteThings []struct {
		Url  string   `yaml:"url" required:"true"`
		Tags []string `yaml:"tags,omitempty"`
	} `yaml:"remoteThings,omitempty"`
}{}

// loadConfig Loads the config file into a struct
func loadConfig(config interface{}) *configmanager.ConfigManager {
	cm, err := configmanager.Init(*_configPath, config, &eriaConfig)
	if err != nil {
		if configmanager.IsFileMissing(err) {
			zlog.Fatal().Str("filePath", *_configPath).Msg("[core:LoadConfig] Config file do not exists...")
		} else {
			zlog.Fatal().Str("filePath", *_configPath).Err(err).Msg("[core:LoadConfig]")
		}
	}

	if err := cm.Load(); err != nil {
		zlog.Fatal().Err(err).Msg("[core:LoadConfig]")
	}
	return cm
}
