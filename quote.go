/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/
package influxpi

import (
	"regexp"
	"strings"
)

var (
	regexpBareIdentifier = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func isBareIdentifier(value string) bool {
	return regexpBareIdentifier.MatchString(value)
}

// Append \ to every newline double quote and backslash
func escapeString(value string) string {
	value = strings.Replace(value, "\\", "\\\\", -1)
	value = strings.Replace(value, "\"", "\\\"", -1)
	return value
}

// QuoteIdentifier returns a query-safe version of an identifier (a database, series
// or measurement name)
func QuoteIdentifier(value string) string {
	if value == "" {
		return value
	} else if isBareIdentifier(value) {
		return value
	} else {
		return "\"" + escapeString(value) + "\""
	}
}
