package shell

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func printTable(header, footer []string, data [][]string, border bool) {
	table := tablewriter.NewWriter(color.Output)

	if header != nil {
		table.SetHeader(header)
	}
	if footer != nil {
		table.SetFooter(footer)
	}

	table.SetBorder(border)

	for _, v := range data {
		table.Append(v)
	}

	shell.Print("\n")
	table.Render()
	shell.Print("\n")
}
