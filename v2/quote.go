/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/
package influxdb

import (
	"fmt"
	"regexp"
	"strings"
	"text/scanner"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS & CONSTS

const (
	StateInit = iota
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

// UnquoteLine returns the measurement name
func UnquoteLine(line string) (string, error) {
	// TODO
	var scan scanner.Scanner
	scan.Init(strings.NewReader(line))
	state := StateInit
	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		fmt.Printf("%s: %s\n", scan.Position, scan.TokenText())
		switch state {
		case StateInit:
			continue
		}
	}
	return "", ErrNotConnected
}
