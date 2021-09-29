package cmd

import (
	"fmt"
	"strings"

	"github.com/RedDocMD/dotted/file"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var commitCommand = &cobra.Command{
	Use:   "commit [<path>|<mnemonic>]+",
	Short: "create a commit for the current version of one or more files",
	Args: func(cmd *cobra.Command, args []string) error {
		if allFiles && len(args) != 0 {
			color.Yellow("Warning: Args are ignored with --all flag")
			return nil
		}
		if !allFiles && len(args) == 0 {
			return errors.Wrap(validationError, color.New(color.FgRed).
				SprintFunc()("expected at least one path/mnemonic"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if allFiles {
			return commitDotFiles(fileStore.Files())
		} else {
			var dotFiles []*file.DotFile
			for _, arg := range args {
				if strings.HasPrefix(arg, "/") {
					file, err := dotFileByPath(fileStore.Files(), arg)
					if err != nil {
						return errors.WithMessage(err, "failed to commit")
					}
					dotFiles = append(dotFiles, file)
				} else {
					file, err := dotFileByMnemonic(fileStore.Files(), arg)
					if err != nil {
						return errors.WithMessage(err, "failed to commit")
					}
					dotFiles = append(dotFiles, file)
				}
			}
			return commitDotFiles(dotFiles)
		}
	},
}

var allFiles bool

func initCommitCommand() {
	commitCommand.Flags().BoolVar(&allFiles, "all", false,
		"commit all dot-files that have been changed")
}

func commitDotFiles(dotFiles []*file.DotFile) error {
	// TODO: Fill the actual commit logic
	return nil
}

func dotFileByPath(dotFiles []*file.DotFile, path string) (*file.DotFile, error) {
	for _, dotFile := range dotFiles {
		if dotFile.Path() == path {
			return dotFile, nil
		}
	}
	return nil, fmt.Errorf("failed to find dotfile at path %s", path)
}

func dotFileByMnemonic(dotFiles []*file.DotFile, mnemonic string) (*file.DotFile, error) {
	for _, dotFile := range dotFiles {
		if dotFile.Mnemonic() == mnemonic {
			return dotFile, nil
		}
	}
	return nil, fmt.Errorf("failed to find dotfile with mnemonic %s", mnemonic)
}
