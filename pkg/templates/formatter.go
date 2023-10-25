package templates

import (
	"os"
	"text/tabwriter"
)

// TODO: add configurable tab indentation, separators, output, other options...

type TableFormatter struct {
}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (f *TableFormatter) Format(columns []string, rows [][]string) error {
	w := tabwriter.NewWriter(os.Stdout, 8, 4, 8, ' ', tabwriter.TabIndent)
	for _, v := range columns {
		w.Write([]byte(v))
		w.Write([]byte("\t"))
	}

	for _, row := range rows {
		w.Write([]byte("\n"))
		for _, cell := range row {
			w.Write([]byte(cell))
			w.Write([]byte("\t"))
		}
	}
	w.Write([]byte("\n"))

	return w.Flush()
}
