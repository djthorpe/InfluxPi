/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package v2

import (
	"fmt"
	"time"

	"github.com/djthorpe/gopi"

	"github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines the configuration parameters for connecting to Influx Database
type dataset struct {
	name      string
	database  string
	precision string
	fields    []string
	tags      map[string]string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// NewDataset returns an empty dataset object
func (this *Client) NewDataset(name string, fields ...string) influxdb.Dataset {
	d := new(dataset)

	// Set database and measurement name
	d.name = name
	d.database = this.database

	// Set precision
	if this.precision == "" {
		d.precision = influxdb.PRECISION_DEFAULT
	} else {
		d.precision = this.precision
	}

	// tags and fields
	d.tags = make(map[string]string, 0)
	d.fields = make([]string, 0, len(fields))
	d.fields = append(d.fields, fields...)

	// return dataset
	return d
}

func (this *Client) Write(dataset influxdb.Dataset) error {
	// TODO
	return gopi.ErrBadParameter
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// SetTag set a tag value if the value and key are not empty strings
func (this *dataset) SetTag(key, value string) {
	if key != "" && value != "" {
		this.tags[key] = value
	}
}

// Get tag keys
func (this *dataset) Tags() []string {
	tags := make([]string, 0, len(this.tags))
	for k := range this.tags {
		tags = append(tags, k)
	}
	return tags
}

// Tag gets tag value, or empty string if tag key does
// not exist
func (this *dataset) Tag(key string) string {
	if value, ok := this.tags[key]; ok {
		return value
	} else {
		return ""
	}
}

// Fields return the field names
func (this *dataset) Fields() []string {
	fields := make([]string, 0, len(this.fields))
	return append(fields, this.fields...)
}

// Database returns the currently set database
func (this *dataset) Database() string {
	return this.database
}

// SetDatabase sets the database name
func (this *dataset) SetDatabase(name string) {
	this.database = name
}

// SetPrecision sets the precision for the timestamp
func (this *dataset) SetPrecision(value string) {
	if value == "" {
		this.precision = influxdb.PRECISION_DEFAULT
	} else {
		this.precision = value
	}
}

// Name returns the measurement name
func (this *dataset) Name() string {
	return this.name
}

// Len returns the number of rows
func (this *dataset) Len() uint {
	// TODO
	return 0
}

// Partial returns true if either the fetched dataset does
// not contain all rows, or the dataset has not yet been
// written to the client
func (this *dataset) Partial() bool {
	// TODO
	return false
}

func (this *dataset) ValuesAtIndex(uint) (time.Time, []influxdb.Value) {
	// TODO
	return time.Time{}, nil
}

func (this *dataset) AddValues(values ...influxdb.Value) error {
	// TODO
	return gopi.ErrNotImplemented
}

func (this *dataset) AddValuesForTimestamp(ts time.Time, values ...influxdb.Value) error {
	// TODO
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *dataset) String() string {
	// TODO
	return fmt.Sprintf("influxdb.Dataset{  }")
}
