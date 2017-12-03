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

	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
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
			this.SetPrecision(influxdb.PRECISION_DEFAULT)
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
		this.database = ""
	}
	return nil
}

func (config Config) addr() string {
	method := "http"
	if config.SSL {
		method = "https"
	}
	if config.Port == 0 {
		config.Port = influxdb.DefaultPortHTTP
	}
	return fmt.Sprintf("%v://%v:%v/", method, config.Host, config.Port)
}

////////////////////////////////////////////////////////////////////////////////
// PARAMETERS

// Version returns the version string for the InfluxDB
func (this *Client) Version() string {
	if this.client == nil {
		return ""
	} else {
		return this.version
	}
}

// Precision returns the current precision value
func (this *Client) Precision() string {
	if this.client == nil {
		return ""
	} else {
		return this.precision
	}
}

// SetPrecision sets precision for setting and returning timestamps
func (this *Client) SetPrecision(value string) error {
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

// Database returns the current database string
func (this *Client) Database() string {
	if this.client == nil {
		return ""
	} else {
		return this.database
	}
}

// SetDatabase sets the current database to use, will
// return ErrBadParameter if the database doesn't exist,
// or ErrNotConnected if the server is not connected
func (this *Client) SetDatabase(name string) error {
	if this.client == nil {
		return influxdb.ErrNotConnected
	}
	if results, err := this.Do(influxdb.ShowDatabases()); err != nil {
		return err
	} else if databases, err := results.Column(0, "databases", "name"); err != nil {
		return err
	} else {
		for _, existing_database := range databases {
			if name == existing_database {
				this.database = name
				return nil
			}
		}
	}
	return influxdb.ErrBadParameter
}

////////////////////////////////////////////////////////////////////////////////
// Convenience methods for database and retention policy

func (this *Client) CreateDatabase(name string, policy *influxdb.RetentionPolicy) error {
	if this.client == nil {
		return influxdb.ErrNotConnected
	}
	// Check for existence of database
	if exists, err := this.exists_string(influxdb.ShowDatabases(), "databases", "name", name); err != nil {
		return err
	} else if exists {
		return influxdb.ErrAlreadyExists
	}
	// Perform the creation
	if _, err := this.Do(influxdb.CreateDatabase(name).RetentionPolicy(policy)); err != nil && err != influxdb.ErrEmptyResponse {
		return err
	}
	return nil
}

func (this *Client) DropDatabase(name string) error {
	if this.client == nil {
		return influxdb.ErrNotConnected
	}
	// Perform the drop
	if _, err := this.Do(influxdb.DropDatabase(name)); err != nil && err != influxdb.ErrEmptyResponse {
		return err
	}
	return nil
}

func (this *Client) CreateRetentionPolicy(name string, policy *influxdb.RetentionPolicy) error {
	if this.client == nil {
		return influxdb.ErrNotConnected
	}
	// Check for existence of database
	if exists, err := this.exists_string(influxdb.ShowRetentionPolicies(), "", "name", name); err != nil {
		return err
	} else if exists {
		return influxdb.ErrAlreadyExists
	}
	// Perform the creation
	if _, err := this.Do(influxdb.CreateRetentionPolicy(this.database, name, policy)); err != nil && err != influxdb.ErrEmptyResponse {
		return err
	}
	// Success
	return nil
}

func (this *Client) DropRetentionPolicy(name string) error {
	if this.client == nil {
		return influxdb.ErrNotConnected
	}
	// Perform the drop
	if _, err := this.Do(influxdb.DropRetentionPolicy(this.database, name)); err != nil && err != influxdb.ErrEmptyResponse {
		return err
	}
	return nil
}

func (this *Client) RetentionPolicies() (map[string]*influxdb.RetentionPolicy, error) {
	if this.client == nil {
		return nil, influxdb.ErrNotConnected
	}
	// Perform the query
	if results, err := this.Do(influxdb.ShowRetentionPolicies()); err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, influxdb.ErrUnexpectedResponse
	} else {
		return results[0].ParseRetentionPolicies()
	}
}

////////////////////////////////////////////////////////////////////////////////
// Execute queries and re-format results

func (this *Client) Do(query influxdb.Query) (influxdb.Results, error) {
	if this.client == nil {
		return nil, influxdb.ErrNotConnected
	}
	// Query and sanity check the response
	response, err := this.query(query.String())
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, influxdb.ErrEmptyResponse
	}
	if response.Results[0].Series == nil || len(response.Results[0].Series) == 0 {
		return nil, influxdb.ErrEmptyResponse
	}
	r := make([]*influxdb.Result, 0, len(response.Results))
	for i, result := range response.Results {
		for j, series := range result.Series {
			table := new(influxdb.Result)
			table.Result = i
			table.Series = j
			table.Name = series.Name
			table.Tags = series.Tags
			table.Columns = series.Columns
			table.Values = series.Values
			table.Partial = series.Partial
			r = append(r, table)
		}
	}
	return influxdb.Results(r), nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	if this.client != nil {
		return fmt.Sprintf("influxdb.Client{ connected=true addr=%v%v version=%v precision=%v }", this.addr, this.database, this.Version(), this.precision)
	} else {
		return fmt.Sprintf("influxdb.Client{ connected=false addr=%v%v precision=%v }", this.addr, this.database, this.precision)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Query database and return response or error
func (this *Client) query(query string) (*client.Response, error) {
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

func (this *Client) exists_string(q influxdb.Query, series string, column string, value string) (bool, error) {
	if response, err := this.Do(q); err != nil {
		return false, err
	} else if column, err := response.Column(0, series, column); err != nil {
		return false, err
	} else {
		for _, v := range column {
			if v == value {
				return true, nil
			}
		}
		return false, nil
	}
}
