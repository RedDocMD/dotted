package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RedDocMD/dotted/config"
	"github.com/RedDocMD/dotted/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dtd",
	Short: "Dotted is a dot-file manager with built-in version control and online backup",
	Long: `Dotted is a fast and reliable dot-file manager, which 
gives full version control of individual files 
(along with implicit branching).
Supports multiple backup and restore options.
★ Inspired by Git. Guided by stars. ★`,
}

var configs *config.Config
var fileStore *store.Store

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := fileStore.SaveToDisk(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfigAndStore)
	rootCmd.AddCommand(listCmd)
}

func initConfigAndStore() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	viper.SetConfigName("dotted")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".dotted"))
	viper.AddConfigPath(filepath.Join(home, ".config", "dotted"))

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Fprintln(os.Stderr, "failed to find config file")
		os.Exit(1)
	}
	path := viper.GetViper().ConfigFileUsed()
	configs, err = config.ReadConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fileStore, err = store.LoadFromDisk(configs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
