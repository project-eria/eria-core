package configmanager

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

type testStruct struct {
	A string `default:"A"`
	B uint   `default:"1"`
	C bool   `default:"true"`
	D struct {
		D1 string
	}
	E []struct {
		E1 string
	}
	F string `required:"true"`
}

func generateDefaultConfig() testStruct {
	config := testStruct{
		A: "A",
		B: 1,
		C: true,
		D: struct {
			D1 string
		}{
			D1: "Y",
		},
		E: []struct {
			E1 string
		}{
			{
				E1: "Z",
			},
			{
				E1: "W",
			},
		},
		F: "V",
	}
	return config
}

func TestInit(t *testing.T) {
	currentEnv := os.Getenv("ERIA_CONF_PATH") // Save current env

	// Create dummy file for file exist check
	file, err := ioutil.TempFile("", "test.yaml")
	if err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		file.Write([]byte{0})
	}

	path := os.TempDir()
	fileName := strings.TrimPrefix(file.Name(), path)

	var result testStruct

	type args struct {
		fileName   string
		config     interface{}
		eriaConfig interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *ConfigManager
		wantErr bool
		env     string
	}{
		{
			name:    "Existing file",
			args:    args{fileName: file.Name(), config: &result, eriaConfig: nil},
			want:    &ConfigManager{filepath: file.Name(), config: &result, eriaConfig: nil},
			wantErr: false,
			env:     path,
		},
		{
			name:    "Missing file",
			args:    args{fileName: "notest.yaml", config: &result, eriaConfig: nil},
			wantErr: true,
			env:     path,
		},
		{
			name:    "Missing env",
			args:    args{fileName: fileName, config: &result, eriaConfig: nil},
			wantErr: true,
			env:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ERIA_CONF_PATH", tt.env)
			got, err := Init(tt.args.fileName, tt.args.config, tt.args.eriaConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %+v, want %+v", got, tt.want)
			}
		})
	}
	os.Setenv("ERIA_CONF_PATH", currentEnv) // Restore current env
}

// func TestConfigManager_Load_ValidYaml(t *testing.T) {
// 	config := generateDefaultConfig()
// 	if bytes, err := yaml.Marshal(config); err == nil {
// 		if file, err := ioutil.TempFile("", "test.yaml"); err == nil {
// 			defer file.Close()
// 			defer os.Remove(file.Name())
// 			file.Write(bytes)

// 			var result testStruct

// 			configmanager := &ConfigManager{
// 				filepath: file.Name(),
// 				config:   &result,
// 			}

// 			if err := configmanager.Load(); err != nil {
// 				t.Errorf("Load_ValidYaml() Error: %s", err)
// 			}
// 			if !reflect.DeepEqual(result, config) {
// 				t.Errorf("Load_ValidYaml() %+v, want %+v", result, config)
// 			}
// 		}
// 	} else {
// 		t.Errorf("Load_ValidYaml() failed to marshal config")
// 	}
// }

// func TestConfigManager_Load_InvalidJson(t *testing.T) {
// 	if file, err := ioutil.TempFile("", "test.json"); err == nil {
// 		defer file.Close()
// 		defer os.Remove(file.Name())
// 		file.Write([]byte{0})

// 		var result testStruct
// 		configmanager := &ConfigManager{
// 			filepath: file.Name(),
// 			config:   &result,
// 		}

// 		if err := configmanager.Load(); err == nil {
// 			t.Errorf("Load_InvalidJson() should return an error")
// 		}
// 	}
// }

// func TestConfigManager_Load_Required(t *testing.T) {
// 	config := generateDefaultConfig()
// 	config.F = ""
// 	if bytes, err := yaml.Marshal(config); err == nil {
// 		if file, err := ioutil.TempFile("", "test.json"); err == nil {
// 			defer file.Close()
// 			defer os.Remove(file.Name())
// 			file.Write(bytes)

// 			var result testStruct

// 			configmanager := &ConfigManager{
// 				filepath: file.Name(),
// 				config:   &result,
// 			}
// 			if err := configmanager.Load(); err.Error() != "F is required, but blank" {
// 				t.Errorf("Load_Required() Doesn't returns the correct error: %s", err)
// 			}
// 		}
// 	} else {
// 		t.Errorf("Load_Required() failed to marshal config")
// 	}
// }

// func TestConfigManager_Load_Default(t *testing.T) {
// 	config := generateDefaultConfig()
// 	config.A = ""
// 	config.B = 0
// 	config.C = false
// 	if bytes, err := yaml.Marshal(config); err == nil {
// 		if file, err := ioutil.TempFile("", "test.json"); err == nil {
// 			defer file.Close()
// 			defer os.Remove(file.Name())
// 			file.Write(bytes)

// 			var result testStruct

// 			configmanager := &ConfigManager{
// 				filepath: file.Name(),
// 				config:   &result,
// 			}

// 			if err := configmanager.Load(); err != nil {
// 				t.Errorf("Load_Default() Error: %s", err)
// 			}
// 			if !reflect.DeepEqual(result, generateDefaultConfig()) {
// 				t.Errorf("Load_Default() %+v, want %+v", result, generateDefaultConfig())
// 			}
// 		}
// 	} else {
// 		t.Errorf("Load_Default() failed to marshal config")
// 	}
// }
