package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)
	configPath := filepath.Join("testdata", "config1.yml")
	config, err := ReadConfig(configPath)
	assert.Equal(err, nil)
	expectedConfig := &Config{
		Name: "Linux",
		WithHistory: []FileEntry{
			{
				Path:     ".config/alacritty/alacritty.yml",
				Mnemonic: "alacritty",
			},
			{
				Path:     ".bashrc",
				Mnemonic: "bashrc",
			},
			{
				Path:     ".config/fish/config.fish",
				Mnemonic: "",
			},
		},
		WithoutHistory: []FileEntry{
			{
				Path:     ".tmux.conf",
				Mnemonic: "tmux",
			},
			{
				Path:     ".config/vscode/settings.json",
				Mnemonic: "",
			},
		},
		StoreLocation: ".config/dotted/store",
	}
	assert.Equal(expectedConfig, config)
}

func TestParseInvalidConfig(t *testing.T) {
	configPath := filepath.Join("testdata", "invalid_config1.yml")
	_, err := ReadConfig(configPath)
	assert.NotNil(t, err)
	configPath = filepath.Join("testdata", "invalid_config2.yml")
	_, err = ReadConfig(configPath)
	assert.NotNil(t, err)
	configPath = filepath.Join("testdata", "invalid_config3.yml")
	_, err = ReadConfig(configPath)
	assert.NotNil(t, err)
	if runtime.GOOS == "windows" {
		configPath = filepath.Join("testdata", "invalid_config4_windows.yml")
	} else {
		configPath = filepath.Join("testdata", "invalid_config4_linux.yml")
	}
	_, err = ReadConfig(configPath)
	assert.NotNil(t, err)
}

func TestParseIncompleteConfig(t *testing.T) {
	assert := assert.New(t)
	configPath := filepath.Join("testdata", "config2.yml")
	config, err := ReadConfig(configPath)
	assert.Equal(err, nil)
	expectedConfig := &Config{
		Name: "Linux",
		WithHistory: []FileEntry{
			{
				Path:     ".config/alacritty/alacritty.yml",
				Mnemonic: "alacritty",
			},
			{
				Path:     ".bashrc",
				Mnemonic: "bashrc",
			},
			{
				Path:     ".config/fish/config.fish",
				Mnemonic: "",
			},
		},
		StoreLocation: ".config/dotted/store",
	}
	assert.Equal(expectedConfig, config)
}
