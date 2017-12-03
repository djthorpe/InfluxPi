/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxpi

import (
	"errors"
	"fmt"
	"time"

	gopi "github.com/djthorpe/gopi"
	client "github.com/influxdata/influxdb/client/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config defines the configuration parameters for connecting to Influx Database
type Config struct {
	Host      string
	Port      uint
	SSL       bool
	Database  string
	Username  string
	Password  string
	Precision string
	Timeout   time.Duration
}

// Client defines a connection to an Influx Database
type Client struct {
	log       gopi.Logger
	database  string
	addr      string
	config    client.HTTPConfig
	precision string
	client    client.Client
	version   string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// DefaultPortHTTP defines the default InfluxDB port used for HTTP
	DefaultPortHTTP uint = 8086
)

const (
	PRECISION_NANO    string = "ns"
	PRECISION_MICRO   string = "µ"
	PRECISION_MICRO2  string = "u"
	PRECISION_MILLI   string = "ms"
	PRECISION_SECOND  string = "s"
	PRECISION_MINUTE  string = "m"
	PRECISION_HOUR    string = "h"
	PRECISION_DAY     string = "d"
	PRECISION_WEEK    string = "w"
	PRECISION_DEFAULT        = PRECISION_MILLI
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// ErrNotConnected is returned when database is not connected (Close has been called)
	ErrNotConnected = errors.New("No connection")

	// ErrNotFound is returned when something expected is not found
	ErrNotFound = errors.New("Not Found")

	// ErrAlreadyExists is returned when something already exists
	ErrAlreadyExists = errors.New("Already Exists")

	// ErrUnexpectedResponse is returned when server does not return with expected data
	ErrUnexpectedResponse = errors.New("Unexpected response from server")

	// ErrEmptyResponse is returned when response does not contain any data
	ErrEmptyResponse = errors.New("Empty response from server")

	// ErrBadParameter is returned when some calling parameter is invalid
	ErrBadParameter = errors.New("Bad Parameter")

	// ErrNotSupported is returned if a feature is not yet supported
	ErrNotSupported = errors.New("Not supported")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open returns an InfluxDB client object
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug2("<influxdb.Client>Open{ addr=%v database=%v }", config.addr(), config.Database)

	this := new(Client)
	this.log = log
	this.addr = config.addr()
	this.config = client.HTTPConfig{
		Addr:     this.addr,
		Username: config.Username,
		Password: config.Password,
		Timeout:  config.Timeout,
	}

	var err error
	if this.client, err = client.NewHTTPClient(this.config); err != nil {
		return nil, this.log.Error("%v", err)
	}

	// Ping client to make sure it exists, get InfluxDB version
	var t time.Duration
	if t, this.version, err = this.client.Ping(this.config.Timeout); err != nil {
		this.client.Close()
		this.client = nil
		return nil, this.log.Error("%v", err)
	}
	this.log.Debug("InfluxDB Version=%v Ping=%v", this.version, t)

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
			this.SetPrecision(PRECISION_DEFAULT)
		}
	}

	// Return success
	return this, nil
}

// Close releases any resources associated with the client connection
func (this *Client) Close() error {
	this.log.Debug2("<influxdb.Client>Close")
	if this.client != nil {
		if err := this.client.Close(); err != nil {
			this.client = nil
			return err
		}
		this.client = nil
	}
	return nil
}

func (config Config) addr() string {
	method := "http"
	if config.SSL {
		method = "https"
	}
	if config.Port == 0 {
		config.Port = DefaultPortHTTP
	}
	return fmt.Sprintf("%v://%v:%v/", method, config.Host, config.Port)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetVersion returns the version string for the InfluxDB
func (this *Client) GetVersion() string {
	if this.client == nil {
		return ""
	} else {
		return this.version
	}
}

func (this *Client) SetPrecision(value string) error {
	switch value {
	case
		PRECISION_NANO, PRECISION_MICRO, PRECISION_MILLI,
		PRECISION_SECOND, PRECISION_MINUTE, PRECISION_HOUR,
		PRECISION_DAY, PRECISION_WEEK:
		this.precision = value
	case PRECISION_MICRO2:
		this.precision = PRECISION_MICRO
	default:
		return ErrBadParameter
	}
	return nil
}

// GetDatabase returns the current database
func (this *Client) GetDatabase() string {
	if this.client == nil {
		return ""
	} else {
		return this.database
	}
}

// SetDatabase sets the current database to use, will
// return ErrEmptyResponse if the database doesn't exist
func (this *Client) SetDatabase(name string) error {
	if this.client == nil {
		return ErrNotConnected
	}
	if databases, err := this.ShowDatabases(); err != nil {
		return err
	} else {
		for _, value := range databases {
			if value == name {
				this.database = value
				return nil
			}
		}
	}
	return ErrNotFound
}

// ShowDatabases enumerates the databases
func (this *Client) ShowDatabases() ([]string, error) {
	if this.client == nil {
		return nil, ErrNotConnected
	}
	if values, err := this.queryScalar("SHOW DATABASES", "databases", "name"); err != nil {
		return nil, err
	} else {
		return values, nil
	}
}

// DatabaseExists returns a boolean value. It will return false
// if an error occurred
func (this *Client) DatabaseExists(name string) bool {
	if this.client == nil {
		return false
	}
	if databases, err := this.ShowDatabases(); err != nil {
		return false
	} else {
		for _, database := range databases {
			if database == name {
				return true
			}
		}
	}
	return false
}

// GetMeasurements enumerates the measurements for a database
func (this *Client) GetMeasurements() ([]string, error) {
	if this.client == nil {
		return nil, ErrNotConnected
	}
	if values, err := this.queryScalar("SHOW MEASUREMENTS", "measurements", "name"); err != nil {
		return nil, err
	} else {
		return values, nil
	}
}

// CreateDatabase with an optional retention policy. The retention policy will
// always have the name 'default'
func (this *Client) CreateDatabase(name string, policy *RetentionPolicy) error {
	if this.client == nil {
		return ErrNotConnected
	}
	if this.DatabaseExists(name) {
		return ErrAlreadyExists
	}
	q := "CREATE DATABASE " + QuoteIdentifier(name)
	if policy != nil {
		q = q + " WITH"
		if policy.Duration != 0 {
			q = q + " DURATION " + fmt.Sprintf("%v", policy.Duration)
		}
		if policy.ReplicationFactor != 0 {
			q = q + " REPLICATION " + fmt.Sprintf("%v", policy.ReplicationFactor)
		}
		if policy.ShardGroupDuration != 0 {
			q = q + " SHARD DURATION " + fmt.Sprintf("%v", policy.ShardGroupDuration)
		}
		q = q + " NAME " + QuoteIdentifier("default")
	}
	if _, err := this.query(q); err != nil {
		return err
	}
	return nil
}

// DropDatabase will delete a database. It will return
// ErrNotFound if the database does not exist
func (this *Client) DropDatabase(name string) error {
	if this.client == nil {
		return ErrNotConnected
	}
	if this.DatabaseExists(name) == false {
		return ErrNotFound
	}
	q := "DROP DATABASE " + QuoteIdentifier(name)
	if _, err := this.query(q); err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	if this.client != nil {
		return fmt.Sprintf("influxdb.Client{ connected=true addr=%v%v version=%v precision=%v }", this.addr, this.database, this.GetVersion(), this.precision)
	} else {
		return fmt.Sprintf("influxdb.Client{ connected=false addr=%v%v precision=%v }", this.addr, this.database, this.precision)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Query database and return response or error
func (this *Client) query(query string) (*client.Response, error) {
	if this.client == nil {
		return nil, ErrNotConnected
	}
	if this.database != "" {
		this.log.Debug("<influxdb.Query>{ database=%v, q=%v }", this.database, query)
	} else {
		this.log.Debug("<influxdb.Query>{ database=<nil>, q=%v }", query)
	}
	response, err := this.client.Query(client.Query{
		Command:   query,
		Database:  this.database,
		Precision: this.precision,
	})
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return response, nil
}

// queryTable returns a table structure
func (this *Client) Query(query string) (*Table, error) {
	// Query and sanity check the response
	response, err := this.query(query)
	if err != nil {
		return nil, err
	}
	if len(response.Results) != 1 {
		return nil, ErrEmptyResponse
	}
	if response.Results[0].Series == nil || len(response.Results[0].Series) == 0 {
		return nil, ErrEmptyResponse
	}

	// Don't support multiple resultsets
	if len(response.Results[0].Series) > 1 {
		this.log.Error("Multiple Results is not supported yet")
		return nil, ErrNotSupported
	}

	// Copy model.Row over to Table structure
	series := response.Results[0].Series[0]
	table := new(Table)
	table.Name = series.Name
	table.Tags = series.Tags
	table.Columns = series.Columns
	table.Values = series.Values
	table.Partial = series.Partial
	return table, nil
}

// queryScalar returns a single column of string values
func (this *Client) queryScalar(query, dataset, column string) ([]string, error) {
	// Query and sanity check the response
	table, err := this.Query(query)
	if err != nil {
		return nil, err
	}
	// Sanity check the data returned
	if table.Name != dataset {
		return nil, ErrUnexpectedResponse
	}
	if len(table.Columns) != 1 && table.Columns[0] != column {
		return nil, ErrUnexpectedResponse
	}

	values := make([]string, 0, len(table.Values))
	for _, column := range table.Values {
		if len(column) != 1 {
			return nil, ErrUnexpectedResponse
		}
		if value, ok := column[0].(string); ok == false {
			return nil, ErrUnexpectedResponse
		} else {
			values = append(values, value)
		}

	}

	return values, nil
}
