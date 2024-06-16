package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
)

// TODO: I haven't figure out how to abstract this
func PrintTable(headers []string, rows [][]string) error {
	t := table.New().
		Headers(headers...).
		Rows(rows...)

	fmt.Println(t)
	return nil
}
