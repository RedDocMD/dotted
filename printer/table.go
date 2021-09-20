package printer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ColumnAlignment = int

const (
	CenterAlign ColumnAlignment = iota
	RightAlign
	LeftAlign
)

type TablePrinter interface {
	RowCount() int
	ColumnCount() int
	Value(row, column int) string
	Ipad() int
	ColumnAlignment(column int) ColumnAlignment
}

func TablePrint(table TablePrinter) {
	columnWidths := make([]int, table.ColumnCount())
	rows := table.RowCount()
	cols := table.ColumnCount()
	ipad := table.Ipad()
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			width := utf8.RuneCountInString(table.Value(j, i))
			if columnWidths[i] < width {
				columnWidths[i] = width
			}
		}
	}
	for row := 0; row < rows; row++ {
		if row == 0 {
			fmt.Print("\u250C")
			for col := 0; col < cols; col++ {
				fmt.Print(rule(columnWidths[col] + 2*ipad))
				if col != cols-1 {
					fmt.Print("\u252C")
				}
			}
			fmt.Println("\u2510")
		} else {
			fmt.Print("\u251C")
			for col := 0; col < cols; col++ {
				fmt.Print(rule(columnWidths[col] + 2*ipad))
				if col != cols-1 {
					fmt.Print("\u253C")
				}
			}
			fmt.Println("\u2524")
		}
		for col := 0; col < cols; col++ {
			fmt.Printf("\u2502%s",
				pad(table.Value(row, col),
					columnWidths[col], ipad,
					table.ColumnAlignment(col)))
		}
		fmt.Println("\u2502")
	}
	fmt.Print("\u2514")
	for col := 0; col < cols; col++ {
		fmt.Print(rule(columnWidths[col] + 2*ipad))
		if col != cols-1 {
			fmt.Print("\u2534")
		}
	}
	fmt.Println("\u2518")
}

func spacer(width int) string {
	return strings.Repeat(" ", width)
}

func rule(width int) string {
	return strings.Repeat("\u2500", width)
}

func pad(s string, width, ipad int, alignment ColumnAlignment) string {
	var res string

	rest := width - utf8.RuneCountInString(s)
	switch alignment {
	case RightAlign:
		res = spacer(rest) + s
	case LeftAlign:
		res = s + spacer(rest)
	case CenterAlign:
		left := rest / 2
		right := rest - left
		res = spacer(left) + s + spacer(right)
	}
	pad := spacer(ipad)
	return pad + res + pad
}
