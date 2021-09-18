package config

import (
	"fmt"
	"os"

	"github.com/RedDocMD/dotted/fs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Fs = fs.OsFs
var Afs = fs.OsAfs

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
	configBytes, err := Afs.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}
	var config Config
	err = yaml.UnmarshalStrict(configBytes, &config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}
	if err := config.validateConfig(); err != nil {
		return nil, err
	} else {
		return &config, nil
	}
}

func (config Config) validateConfig() error {
	if len(config.Name) == 0 {
		return errors.New("invalid config: empty name")
	}
	if len(config.StoreLocation) == 0 {
		return errors.New("invalid config: empty store location")
	}
	for _, entry := range config.WithHistory {
		if Fs.IsAbs(entry.Path) {
			return errors.New(fmt.Sprintf("invalid config: %s is an absolute path, all paths must be relative to $HOME", entry.Path))
		}
	}
	return nil
}

func (config *Config) IsStoreAvailable() bool {
	stat, err := Fs.Stat(config.StoreLocation)
	if os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}
