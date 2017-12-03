/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/
package influxdb

import (
	"errors"
	"time"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// DefaultPortHTTP defines the default InfluxDB port used for HTTP
	DefaultPortHTTP uint = 8086
)

const (
	// Precision defines the precision used for time values on inserting
	// data and returning it
	PRECISION_NANO    string = "ns"
	PRECISION_MICRO   string = "µ"
	PRECISION_MICRO2  string = "u"
	PRECISION_MILLI   string = "ms"
	PRECISION_SECOND  string = "s"
	PRECISION_MINUTE  string = "m"
	PRECISION_HOUR    string = "h"
	PRECISION_DAY     string = "d"
	PRECISION_WEEK    string = "w"
	PRECISION_DEFAULT string = PRECISION_MILLI
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// ErrNotConnected is returned when database is not connected (Close has been called)
	ErrNotConnected = errors.New("Not connected")

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
// TYPES

// RetentionPolicy defines a period of time to store measurement data for
type RetentionPolicy struct {
	Duration           time.Duration
	ShardGroupDuration time.Duration
	ReplicationFactor  int
	Default            bool
}

// Result reflects the influxdb model.Row structure but which defines a number
// of additional methods
type Result struct {
	Result  int
	Series  int
	Name    string
	Tags    map[string]string
	Columns []string
	Values  [][]interface{}
	Partial bool
}

// Results is a set of results (usually one, but may be more if more than one measure
// is selected)
type Results []*Result

// Value is a value returned by influxdb
type Value interface{}

// Measurement defines a measurement
type Measurement struct {
	Name     string
	Database string
	Policy   string
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Driver is the abstract driver interface
type Driver interface {
	gopi.Driver

	// Get and set parameters
	Version() string
	Database() string
	SetDatabase(value string) error
	Precision() string
	SetPrecision(value string) error

	// Convenience methods for database and retention policy
	CreateDatabase(name string, policy *RetentionPolicy) error
	CreateRetentionPolicy(name string, policy *RetentionPolicy) error
	DropDatabase(name string) error
	DropRetentionPolicy(name string) error
	RetentionPolicies() (map[string]*RetentionPolicy, error)

	// Excute a query
	Do(query Query) (Results, error)
}

// Query is the abstract InfluxQL statement interface
type Query interface {
	// Set parameters
	Database(value string) Query
	RetentionPolicy(value *RetentionPolicy) Query
	Default(value bool) Query
	Measurement(values ...*Measurement) Query
	OffsetLimit(offset uint, limit uint) Query
	Filter(values ...Predicate) Query

	// Return the query as a string
	String() string
}

// Predicate is an abstract predicate (a tag, a field or a function)
type Predicate interface {
	// Return the predicate as a string
	String() string
}
