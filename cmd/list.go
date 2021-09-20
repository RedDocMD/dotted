package cmd

import (
	"fmt"
	"os"

	"github.com/RedDocMD/dotted/file"
	"github.com/RedDocMD/dotted/printer"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all the dot-files in the store",
	Run: func(cmd *cobra.Command, args []string) {
		var table printer.TablePrinter = FileTable(fileStore.Files())
		printer.TablePrint(table)
	},
}

type FileTable []*file.DotFile

func (table FileTable) RowCount() int {
	return len(table) + 1
}

func (table FileTable) ColumnCount() int {
	return 3
}

func (table FileTable) Value(row, column int) string {
	var columnHeaders = [3]string{"Path", "Mnemonic", "Has History"}
	if row == 0 {
		return columnHeaders[column]
	} else {
		file := table[row-1]
		if column == 0 {
			return file.Path()
		} else if column == 1 {
			return file.Mnemonic()
		} else if column == 2 {
			hasHistory := file.HasHistory()
			if hasHistory {
				return "🗸"
			} else {
				return "✗"
			}
		}
	}
	// Should not reach the following
	fmt.Fprintln(os.Stderr, "invalid column while printing files")
	os.Exit(1)
	return ""
}

func (table FileTable) Ipad() int {
	return 1
}

func (table FileTable) ColumnAlignment(column int) printer.ColumnAlignment {
	if column == 0 || column == 1 {
		return printer.LeftAlign
	} else if column == 2 {
		return printer.CenterAlign
	}
	// Should not reach the following
	fmt.Fprintln(os.Stderr, "invalid column while printing files")
	os.Exit(1)
	return -1
}
