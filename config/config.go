package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Name           string      `yaml:"name"`
	WithHistory    []FileEntry `yaml:"withHistory"`
	WithoutHistory []FileEntry `yaml:"withoutHistory"`
	StoreLocation  string      `yaml:"storeLocation"`
}

type FileEntry struct {
	Path     string
	Mnemonic string
}

func ReadConfig(path string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read config")
	}
	var config Config
	err = yaml.UnmarshalStrict(configBytes, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse config")
	}
	if config.validateConfig() {
		return &config, nil
	} else {
		return nil, errors.New("Name or StoreLocation is empty")
	}
}

func (config *Config) validateConfig() bool {
	if config == nil {
		return false
	}
	if len(config.Name) == 0 || len(config.StoreLocation) == 0 {
		return false
	}
	return true
}
