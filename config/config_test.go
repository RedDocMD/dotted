package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RedDocMD/dotted/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct{ suite.Suite }

func (suite *Config) SetupSuite() {
	Fs = fs.MockFs
	Afs = fs.MockAfs
}

func (suite *Config) TearDownSuite() {
	Fs = fs.OsFs
	Afs = fs.OsAfs
}

func (suite *ConfigSuite) TestParseConfig() {
	assert := assert.New(suite.T())
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
		StoreLocation: Fs.Abs(".config/dotted/store"),
	}
	assert.Equal(expectedConfig, config)
}

func (suite *ConfigSuite) TestParseInvalidConfig() {
	t := suite.T()
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

func (suite *ConfigSuite) TestParseIncompleteConfig() {
	assert := assert.New(suite.T())
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
		StoreLocation: Fs.Abs(".config/dotted/store"),
	}
	assert.Equal(expectedConfig, config)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &ConfigSuite{})
}
