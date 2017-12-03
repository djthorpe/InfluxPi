/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxdb

import (
	"encoding/json"
	"fmt"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return column values for a particular series of values for a particular
// query result. When using single statements, the result parameter should
// be zero
func (r Results) Column(result int, series string, column string) ([]Value, error) {
	if result >= len(r) {
		return nil, ErrBadParameter
	}
	for _, resultset := range r {
		if resultset.Result == result && resultset.Name == series {
			return resultset.column(column)
		}
	}
	return nil, ErrBadParameter
}

// ParseRetentionPolicies returns retention policies from a server
// response
func (r *Result) ParseRetentionPolicies() (map[string]*RetentionPolicy, error) {
	policies := make(map[string]*RetentionPolicy, len(r.Values))
	for _, row := range r.Values {
		if len(row) != 5 {
			return nil, ErrUnexpectedResponse
		}
		if name, ok := row[0].(string); ok == false {
			return nil, ErrUnexpectedResponse
		} else if duration, ok := row[1].(string); ok == false {
			return nil, ErrUnexpectedResponse
		} else if shard_duration, ok := row[2].(string); ok == false {
			return nil, ErrUnexpectedResponse
		} else if replication_factor, ok := row[3].(json.Number); ok == false {
			return nil, ErrUnexpectedResponse
		} else if default_policy, ok := row[4].(bool); ok == false {
			return nil, ErrUnexpectedResponse
		} else if duration2, err := time.ParseDuration(duration); err != nil {
			return nil, ErrUnexpectedResponse
		} else if shard_duration2, err := time.ParseDuration(shard_duration); err != nil {
			return nil, ErrUnexpectedResponse
		} else if replication_factor2, err := replication_factor.Int64(); err != nil {
			return nil, ErrUnexpectedResponse
		} else {
			policies[name] = &RetentionPolicy{
				Duration:           duration2,
				ShardGroupDuration: shard_duration2,
				ReplicationFactor:  int(replication_factor2),
				Default:            default_policy,
			}
		}
	}
	return policies, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (r *Result) columnindex(column string) int {
	for i := range r.Columns {
		if r.Columns[i] == column {
			return i
		}
	}
	return -1
}

func (r *Result) column(column string) ([]Value, error) {
	if i := r.columnindex(column); i >= 0 && i < len(r.Columns) {
		c := make([]Value, len(r.Values))
		for j := range r.Values {
			c[j] = toValue(column, r.Values[j][i])
		}
		return c, nil
	} else {
		return nil, ErrBadParameter
	}
}

func toValue(col string, value interface{}) Value {
	switch value.(type) {
	case json.Number:
		if col == "time" {
			// TODO
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
