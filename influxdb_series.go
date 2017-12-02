/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxpi

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Tag defines a single key/value pair
type Series struct {
	Measurement string
	Tags        map[string]string
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Client) ShowSeries() ([]*Series, error) {
	if this.client == nil {
		return nil, ErrNotConnected
	}
	response, err := this.Query("SHOW SERIES")
	if err != nil {
		return nil, err
	}
	// check for case where there are no series
	if response.Length() == 0 {
		return []*Series{}, ErrEmptyResponse
	}
	// iterate through the series
	series := make([]*Series, response.Length())
	for i := range response.Values {
		hash := response.Row(i)
		series[i] = &Series{}
		// TODO
		fmt.Printf("%v\n", hash)
	}
	return series, nil
}

func (this *Client) ShowMeasurements(regexp *RegExp, offset *Offset) ([]string, error) {
	//	SHOW MEASUREMENTS [ON <database_name>] [WITH MEASUREMENT <regular_expression>] [WHERE <tag_key> <operator> ['<tag_value>' | <regular_expression>]] [LIMIT_clause] [OFFSET_clause]
	if this.client == nil {
		return nil, ErrNotConnected
	}
	q := "SHOW MEASUREMENTS"
	if this.database != "" {
		q = q + " ON " + QuoteIdentifier(this.database)
	}
	if regexp != nil {
		q = q + " WITH MEASUREMENT " + regexp.String()
	}
	if offset != nil {
		q = q + " " + offset.String()
	}
	if response, err := this.Query(q); err != nil {
		return nil, err
	} else {
		// check for case where there are no measurements
		if response.Length() == 0 {
			return []string{}, ErrEmptyResponse
		}
		// iterate through
		measurements := make([]string, response.Length())
		for i := range response.Values {
			hash := response.Row(i)
			if str, ok := hash["name"].(string); ok == false {
				return nil, ErrUnexpectedResponse
			} else {
				measurements[i] = str
			}
		}
		return measurements, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s *Series) String() string {
	return fmt.Sprintf("influxdb.Series{ Measurement=%v Tags=%v }", s.Measurement, s.Tags)
}
