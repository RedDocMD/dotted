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
			return fmt.Errorf("expected at least one path/mnemonic")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var dotFiles []*file.DotFile
		if allFiles {
			dotFiles = fileStore.Files()
		} else {
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
		}
		withHistory, withoutHistory := splitDotFilesByHistory(dotFiles)
		commitCnt, err := commitDotFiles(withHistory)
		if err != nil {
			return err
		}
		updateCnt, err := updateDotFiles(withoutHistory)
		if err != nil {
			return err
		}
		if commitCnt <= 1 {
			color.Green("Committed %d file", commitCnt)
		} else {
			color.Green("Committed %d files", commitCnt)
		}
		if updateCnt <= 1 {
			color.Green("Updated %d file", updateCnt)
		} else {
			color.Green("Updated %d files", updateCnt)
		}
		return nil
	},
}

var allFiles bool

func initCommitCommand() {
	commitCommand.Flags().BoolVar(&allFiles, "all", false,
		"commit all dot-files that have been changed")
}

func commitDotFiles(dotFiles []*file.DotFile) (int, error) {
	cnt := 0
	for _, dotFile := range dotFiles {
		done, err := dotFile.AddCommit()
		if err != nil {
			return -1, errors.WithMessage(err, "failed to commit")
		}
		if done {
			fmt.Printf("Committed %s\n", dotFile.Path())
			cnt += 1
		}
	}
	return cnt, nil
}

func updateDotFiles(dotFiles []*file.DotFile) (int, error) {
	cnt := 0
	for _, dotFile := range dotFiles {
		done, err := dotFile.UpdateContent()
		if err != nil {
			return -1, errors.WithMessage(err, "failed to commit")
		}
		if done {
			fmt.Printf("Updated %s\n", dotFile.Path())
			cnt += 1
		}
	}
	return cnt, nil
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

func splitDotFilesByHistory(dotFiles []*file.DotFile) ([]*file.DotFile, []*file.DotFile) {
	var withHistory, withoutHistory []*file.DotFile
	for _, dotFile := range dotFiles {
		if dotFile.HasHistory() {
			withHistory = append(withHistory, dotFile)
		} else {
			withoutHistory = append(withoutHistory, dotFile)
		}
	}
	return withHistory, withoutHistory
}
