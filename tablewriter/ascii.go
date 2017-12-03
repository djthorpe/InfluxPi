package tablewriter

import (
	"fmt"
	"io"

	influx "github.com/djthorpe/InfluxPi"
	"github.com/olekukonko/tablewriter"
)

func RenderASCII(table *influx.Table, writer io.Writer) error {
	out := tablewriter.NewWriter(writer)
	out.SetHeader(table.Columns)
	out.SetAutoMergeCells(true)
	out.SetCaption(true, table.Name)
	out.SetAutoFormatHeaders(false)
	row := make([]string, len(table.Columns))
	for i := range table.Values {
		out.Append(asStringArray(table.RowArray(i), row))
	}
	out.Render()
	return nil
}

func asStringArray(in []interface{}, out []string) []string {
	if len(out) != len(in) {
		panic("out != in")
	}
	for i := range in {
		if in[i] == nil {
			out[i] = "<nil>"
		} else {
			out[i] = fmt.Sprintf("%v", in[i])
		}
	}
	return out
}
