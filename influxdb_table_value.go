/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxpi

import (
	"encoding/json"
	"fmt"
	"time"
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

// Length returns the number of rows in the table
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
	row := make(map[string]interface{}, len(this.Columns))
	for j, name := range this.Columns {
		row[name] = toValue(name, this.Values[i][j])
	}
	return row
}

// RowArray returns an array of values for a row
func (this *Table) RowArray(i int) []interface{} {
	// Sanity check on index parameter
	if i < 0 || i >= len(this.Values) {
		return nil
	}
	if len(this.Values[i]) != len(this.Columns) {
		return nil
	}
	row := make([]interface{}, len(this.Columns))
	for j, name := range this.Columns {
		row[j] = toValue(name, this.Values[i][j])
	}
	return row
}

func (this *Table) String() string {
	return fmt.Sprintf("influxdb.Table{ name=%v tags=%v columns=%v values=%v }", this.Name, this.Tags, this.Columns, this.Values)
}

////////////////////////////////////////////////////////////////////////////////
// Value methods

func toValue(col string, value interface{}) Value {
	switch value.(type) {
	case json.Number:
		if col == "time" {
			if n, err := value.(json.Number).Float64(); err == nil {
				secs := n / 1E3
				return Value(time.Unix(int64(secs), 0))
			}
		} else if n, err := value.(json.Number).Float64(); err == nil {
			return Value(n)
		}
	case string:
		if col == "time" {
			if t, err := time.Parse(time.RFC3339Nano, value.(string)); err == nil {
				fmt.Println("time=", value, t)
				return Value(t)
			}
		}
	}
	return Value(value)
}

/*
	fmt.Println(col, value)
	if col == "time" {
		if str, ok := value.(string); ok == false {
			return Value(value)
		} else if t, err := time.Parse(time.RFC1123Z, str); err != nil {
			return Value(value)
		} else {
			return Value(t)
		}
	} else {
		return Value(value)
	}
}
*/
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
