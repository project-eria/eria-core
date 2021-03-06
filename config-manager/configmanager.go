package configmanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

// ConfigManager struct
type ConfigManager struct {
	filepath string
	yaml     []byte
	s        interface{}
}

const (
	eNotFound = "Config file missing"
)

// Init config manager with filename, and a struct
func Init(filePath string, s interface{}) (*ConfigManager, error) {
	configManager := &ConfigManager{
		filepath: filePath,
		s:        s,
	}
	log.Trace().Str("filePath", filePath).Msg("[configmanager:Init] Looking for file")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return configManager, errors.New(eNotFound)
	}

	return configManager, nil
}

// IsFileMissing return true if this is a file missing error
func IsFileMissing(e error) bool {
	return e.Error() == eNotFound
}

// Load config from file, based on the configmanger parameters
func (c *ConfigManager) Load() error {
	log.Trace().Str("filePath", c.filepath).Msg("[configmanager:Load] Loading config")
	bytes, err := ioutil.ReadFile(c.filepath)
	if err != nil {
		// TODO What to do if file doesn't exists
		return err
	}

	// Save as yaml string
	c.yaml = bytes

	if err := yaml.Unmarshal(bytes, c.s); err != nil {
		// TODO What to do if not valid file
		return err
	}

	if err := processTags(c.s); err != nil {
		return err
	}
	if e := log.Trace(); e.Enabled() {
		content, _ := json.MarshalIndent(c.s, "", "  ")
		e.Msgf("[configmanager:Load] %s", string(content))
	}
	return nil
}

// Save config to file, based on the configmanger parameters
/*
func (c *ConfigManager) Save() error {
	c.Lock()

	log.Trace().Str("module", "configmanager").Str("filename", c.filepath).Msg("Saving config")

	bytes, err := yaml.Marshal(c.s)
	if err != nil {
		return err
	}

	c.Unlock()
	return ioutil.WriteFile(c.filepath, bytes, 0644)
}
*/
func processTags(config interface{}) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			// Set default configuration if blank
			if value := fieldStruct.Tag.Get("default"); value != "" {
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			} else if fieldStruct.Tag.Get("required") == "true" {
				// return error if it is required but blank
				return errors.New(fieldStruct.Name + " is required, but blank")
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := processTags(field.Addr().Interface()); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := processTags(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
