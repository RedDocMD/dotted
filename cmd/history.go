package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history <path|mnemonic>",
	Short: "view the history of a file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected exactly one path/mnemonic as arg")
		}
		if list == view {
			return fmt.Errorf("expected exactly one of list or view")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var list, view bool

func initHistoryCommand() {
	historyCmd.Flags().BoolVar(&list, "list", false, "list all commits of the file")
	historyCmd.Flags().BoolVar(&view, "view", false, "view the file at a specified commit")
}
