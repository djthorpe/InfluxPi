/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxdb

import (
	"regexp"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS & CONSTS

const (
	reserved_words = `
		ALL          ALTER        AS           ASC          BEGIN        BY
		CREATE       CONTINUOUS   DATABASE     DATABASES    DEFAULT      DELETE
		DESC         DROP         DURATION     END          EXISTS       EXPLAIN
		FIELD        FROM         GRANT        GROUP        IF           IN
		INNER        INSERT       INTO         KEY          KEYS         LIMIT
		SHOW         MEASUREMENT  MEASUREMENTS NOT          OFFSET       ON           
		ORDER        PASSWORD     POLICY       POLICIES     PRIVILEGES   QUERIES      
		QUERY        READ         REPLICATION  RETENTION    REVOKE       SELECT       
		SERIES       SERVER       SHARD        SLIMIT       SOFFSET      TAG          
		TO           USER         USERS        VALUES       WHERE        WITH      
		WRITE
	`
)

var (
	regexpBareIdentifier = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
	reservedWords        = make(map[string]bool, 0)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Quote returns a query-safe version of an identifier (a database, series
// or measurement name)
func Quote(value string) string {
	if value == "" {
		return value
	} else if isReservedWord(value) {
		return "\"" + escapeString(value) + "\""
	} else if isBareIdentifier(value) {
		return value
	} else {
		return "\"" + escapeString(value) + "\""
	}
}

// Returns a version of the string which has double-quotes
// around it, and will replace baskslashes with double backslashes
// and double quotes with backslash + double quote
func QuoteString(value string) string {
	return "\"" + escapeString(value) + "\""
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func isBareIdentifier(value string) bool {
	return regexpBareIdentifier.MatchString(value)
}

func isReservedWord(value string) bool {
	w := strings.TrimSpace(strings.ToUpper(value))
	_, exists := reservedWords[w]
	return exists
}

// Append \ to every newline double quote and backslash
func escapeString(value string) string {
	value = strings.Replace(value, "\\", "\\\\", -1)
	value = strings.Replace(value, "\"", "\\\"", -1)
	return value
}

////////////////////////////////////////////////////////////////////////////////
// RESERVED WORDS HASH

func init() {
	for _, word := range strings.Fields(reserved_words) {
		reservedWords[strings.ToUpper(word)] = true
	}
}
