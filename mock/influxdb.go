/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/
package mock

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines the configuration parameters for the mock influx database
type Config struct {
	Database  string
	Precision string
}

// Driver defines a connection to an Influx Database
type Driver struct {
	log       gopi.Logger
	connected bool
	database  string
	precision string
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open returns an InfluxDB client object
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<influxdb.mock.Client>Open{ database=%v }", config.Database)

	this := new(Driver)
	this.log = log

	// Connected
	this.connected = true

	// Set database
	if config.Database != "" {
		if err := this.SetDatabase(config.Database); err != nil {
			return nil, this.log.Error("Unknown database: %v", config.Database)
		}
	}

	// Set precision
	if config.Precision != "" {
		if err := this.SetPrecision(config.Precision); err != nil {
			return nil, err
		} else {
			this.SetPrecision(influxdb.PRECISION_DEFAULT)
		}
	}

	// Return success
	return this, nil
}

// Close releases any resources associated with the client connection
func (this *Driver) Close() error {
	this.log.Debug2("<influxdb.mock.Client>Close")
	this.connected = false
	this.database = ""
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GET AND SET PARAMETERS

func (this *Driver) Version() string {
	if this.connected == false {
		return ""
	} else {
		return "<influxdb.Driver.mock>"
	}
}

func (this *Driver) Database() string {
	return this.database
}

func (this *Driver) SetDatabase(value string) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	if value == "" {
		return influxdb.ErrBadParameter
	} else {
		this.database = value
	}
	return nil
}

func (this *Driver) Precision() string {
	return this.precision
}

func (this *Driver) SetPrecision(value string) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	switch value {
	case
		influxdb.PRECISION_NANO, influxdb.PRECISION_MICRO, influxdb.PRECISION_MILLI,
		influxdb.PRECISION_SECOND, influxdb.PRECISION_MINUTE, influxdb.PRECISION_HOUR,
		influxdb.PRECISION_DAY, influxdb.PRECISION_WEEK:
		this.precision = value
	case influxdb.PRECISION_MICRO2:
		this.precision = influxdb.PRECISION_MICRO
	default:
		return influxdb.ErrBadParameter
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DATABASE AND RETENTION POLICIES

func (this *Driver) CreateDatabase(name string, policy *influxdb.RetentionPolicy) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	if _, err := this.Do(influxdb.CreateDatabase(name).RetentionPolicy(policy)); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Driver) CreateRetentionPolicy(name string, policy *influxdb.RetentionPolicy) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	return influxdb.ErrNotSupported
}

func (this *Driver) DropDatabase(name string) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	return influxdb.ErrNotSupported
}

func (this *Driver) DropRetentionPolicy(name string) error {
	if this.connected == false {
		return influxdb.ErrNotConnected
	}
	return influxdb.ErrNotSupported
}

////////////////////////////////////////////////////////////////////////////////
// PERFORM QUERY

func (this *Driver) Do(query influxdb.Query) (influxdb.Results, error) {
	if this.connected == false {
		return nil, influxdb.ErrNotConnected
	}
	this.log.Debug2("Do(%v)", query.String())
	return nil, influxdb.ErrNotSupported
}
