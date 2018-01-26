/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package v2

import (
	"github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines the configuration parameters for connecting to Influx Database
type dataset struct {
	name      string
	database  string
	precision string
	tags      map[string]string
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewDataset returns an empty dataset object
func (this *Client) NewDataset(name string) influxdb.Dataset {
	d := new(dataset)
	d.name = name
	d.tags = make(map[string]string, 0)
	d.database = this.database
	d.precision = this.precision
	return d
}

// Set a tag value
func (this *dataset) SetTag(key, value string) {
	this.Tags[key] = value
}
