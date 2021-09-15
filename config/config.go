package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
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
	return &config, nil
}
