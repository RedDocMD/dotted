package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)
	configPath := filepath.Join("testdata", "config1.yml")
	config, err := ReadConfig(configPath)
	assert.Equal(err, nil)
	expectedConfig := &Config{
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
