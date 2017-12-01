/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package influxpi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// RetentionPolicy defines a period of time to store measurement data for
type RetentionPolicy struct {
	Name               string
	Duration           time.Duration
	ShardGroupDuration time.Duration
	ReplicationFactor  int
	Default            bool
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetRetentionPolicies enumerates the retention policies for a database
func (this *Client) GetRetentionPolicies() (map[string]*RetentionPolicy, error) {
	response, err := this.Query("SHOW RETENTION POLICIES")
	if err != nil {
		return nil, err
	}
	// check for case where there are no retention policies (there should
	// always be the "autogen" policy?)
	if response.Length() == 0 {
		return nil, ErrEmptyResponse
	}
	// Now iterate through the retention policies
	policies := make(map[string]*RetentionPolicy, response.Length())
	for i := range response.Values {
		var ok bool
		var name string
		hash := response.Row(i)
		policy := new(RetentionPolicy)

		if name, ok = hash["name"].(string); ok == false {
			this.log.Debug("Invalid name parameter")
			return nil, ErrUnexpectedResponse
		} else {
			policy.Name = name
		}
		if policy.Default, ok = hash["default"].(bool); ok == false {
			this.log.Debug("Invalid default parameter")
			return nil, ErrUnexpectedResponse
		}
		if replicationFactor, ok := hash["replicaN"].(json.Number); ok == false {
			this.log.Debug("Invalid replicaN parameter, %T", hash["replicaN"])
			return nil, ErrUnexpectedResponse
		} else if policy.ReplicationFactor, ok = parseInt(string(replicationFactor)); ok == false {
			this.log.Debug("Invalid replicaN parameter, %v", replicationFactor)
			return nil, ErrUnexpectedResponse
		}
		if duration, ok := hash["duration"].(string); ok == false {
			this.log.Debug("Invalid duration parameter, %T", hash["duration"])
			return nil, ErrUnexpectedResponse
		} else if policy.Duration, ok = parseDuration(duration); ok == false {
			this.log.Debug("Invalid duration parameter, %v", duration)
			return nil, ErrUnexpectedResponse
		}
		if shardGroupDuration, ok := hash["shardGroupDuration"].(string); ok == false {
			this.log.Debug("Invalid shardGroupDuration parameter, %T", hash["shardGroupDuration"])
			return nil, ErrUnexpectedResponse
		} else if policy.ShardGroupDuration, ok = parseDuration(shardGroupDuration); ok == false {
			this.log.Debug("Invalid shardGroupDuration parameter, %v", shardGroupDuration)
			return nil, ErrUnexpectedResponse
		}

		policies[name] = policy
	}
	return policies, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *RetentionPolicy) String() string {
	return fmt.Sprintf("influxdb.RetentionPolicy{ Name=%v Duration=%v ShardGroupDuration=%v ReplicationFactor=%v Default=%v }", this.Name, this.Duration, this.ShardGroupDuration, this.ReplicationFactor, this.Default)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func parseInt(value string) (int, bool) {
	if int64Value, err := strconv.ParseInt(value, 10, 64); err != nil {
		return 0, false
	} else {
		return int(int64Value), true
	}
}

func parseDuration(value string) (time.Duration, bool) {
	if duration, err := time.ParseDuration(value); err != nil {
		return 0, false
	} else {
		return duration, true
	}
}
