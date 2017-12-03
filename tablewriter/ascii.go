package tablewriter

import (
	"fmt"
	"io"

	"github.com/djthorpe/influxdb"
	"github.com/olekukonko/tablewriter"
)

func RenderASCII(result *influxdb.Result, writer io.Writer) error {
	out := tablewriter.NewWriter(writer)
	out.SetHeader(result.Columns)
	out.SetAutoMergeCells(true)
	out.SetCaption(true, result.Name)
	out.SetAutoFormatHeaders(false)
	row := make([]string, len(result.Columns))
	for i := range result.Values {
		out.Append(asStringArray(result.Values[i], row))
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
