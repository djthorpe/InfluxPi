/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxpi

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Value is a container for a cell
type Value interface{}

// Table is a query result, which reflects the influxdb model.Row structure
// but which defines a number of additional methods
type Table struct {
	Name    string
	Tags    map[string]string
	Columns []string
	Values  [][]interface{}
	Partial bool
}

////////////////////////////////////////////////////////////////////////////////
// Table methods

func (this *Table) Length() int {
	return len(this.Values)
}

// Row returns a map of row column name to row value
func (this *Table) Row(i int) map[string]interface{} {
	// Sanity check on index parameter
	if i < 0 || i >= len(this.Values) {
		return nil
	}
	if len(this.Values[i]) != len(this.Columns) {
		return nil
	}
	row := make(map[string]interface{}, len(this.Values[i]))
	for j, name := range this.Columns {
		row[name] = Value(this.Values[i][j])
	}
	return row
}

func (this *Table) String() string {
	return fmt.Sprintf("influxdb.Table{ name=%v tags=%v columns=%v values=%v }", this.Name, this.Tags, this.Columns, this.Values)
}

////////////////////////////////////////////////////////////////////////////////
// Value methods
/*
func (value Value) String() string {
	switch value.(type) {
	case string:
		return fmt.Sprintf("<string>%v", value.(string))
	case bool:
		return fmt.Sprintf("<bool>%v", value.(bool))
	case int:
		return fmt.Sprintf("<int>%v", value.(int))
	case uint:
		return fmt.Sprintf("<uint>%v", value.(uint))
	default:
		return fmt.Sprintf("<other>%v", value)
	}
}
*/
