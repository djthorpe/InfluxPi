/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/
package influxdb

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r Results) String() string {
	if len(r) == 0 {
		return "<influxdb.Result>{}"
	} else if len(r) == 1 {
		return r[0].String()
	} else {
		results := ""
		for _, result := range r {
			results = results + "," + result.String()
		}
		return fmt.Sprintf("<influxdb.Results>[ %v ]", strings.TrimSuffix(results, ","))
	}
}

func (r *Result) String() string {
	return fmt.Sprintf("<influxdb.Result>{ result=%v series=%v name=%v columns=%v number_of_rows=%v partial=%v }", r.Result, r.Series, r.Name, r.Columns, len(r.Values), r.Partial)
}

func (this *RetentionPolicy) String() string {
	return fmt.Sprintf("<influxdb.RetentionPolicy>{ Duration=%v ShardGroupDuration=%v ReplicationFactor=%v Default=%v }", this.Duration, this.ShardGroupDuration, this.ReplicationFactor, this.Default)
}
