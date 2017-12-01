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
	Name string
	Tags map[string]string
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
		return nil, ErrEmptyResponse
	}
	// iterate through the series
	series := make([]*Series, response.Length())
	for i, v := range response.Values {
		series[i] = &Series{}
		// TODO
		fmt.Println(len(v))
	}
	return series, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s *Series) String() string {
	return fmt.Sprintf("influxdb.Series{ Name=%v %v }", s.Name, s.Tags)
}
