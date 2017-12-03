/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxpi

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RegExp struct {
	Value string
}

type offsetlimit struct {
	Limit  uint
	Offset uint
}

type columns struct {
	value string
}

type DataSource struct {
	Measurement     string
	Database        string
	RetentionPolicy string
}

type s struct {
	d []*DataSource
	c *columns
	o *offsetlimit
}

type Statement interface {
	// Set columns
	Columns(string) Statement

	// Set offset and limit
	Offset(uint) Statement
	Limit(uint) Statement

	// Write out statement
	Statement() string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTORS

// Select returns a select statement
func (this *Client) Select(from ...*DataSource) Statement {
	if len(from) == 0 {
		this.log.Error("Call to Select requires at least one data source")
		return nil
	}
	return &s{
		d: from,
	}
}

func (this *Client) Do(statement Statement) (*Table, error) {
	if this.client == nil {
		return nil, ErrNotConnected
	}
	// Execute query
	if response, err := this.Query(statement.Statement()); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STATEMENT IMPLEMENTATION FOR SELECT

func (this *s) Limit(limit uint) Statement {
	if this.o == nil {
		this.o = &offsetlimit{Limit: limit}
	} else {
		this.o.Limit = limit
	}
	return this
}

func (this *s) Offset(offset uint) Statement {
	if this.o == nil {
		this.o = &offsetlimit{Offset: offset}
	} else {
		this.o.Offset = offset
	}
	return this
}

func (this *s) Columns(value string) Statement {
	if value == "" {
		this.c = nil
	} else {
		this.c = &columns{value}
	}
	return this
}

func (this *s) Statement() string {
	q := "SELECT "
	// COLUMNS
	if this.c != nil {
		q = q + this.c.String()
	} else {
		q = q + "*"
	}
	// DATA SOURCES
	if len(this.d) > 0 {
		q = q + " FROM "
		for i := range this.d {
			q = q + this.d[i].String() + ","
		}
		q = strings.TrimSuffix(q, ",")
	}
	// LIMIT and OFFSET
	if this.o != nil {
		q = q + " " + this.o.String()
	}

	return q
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (o *offsetlimit) String() string {
	if o.Limit == 0 && o.Offset == 0 {
		return ""
	}
	s := ""
	if o.Limit > 0 {
		s = s + fmt.Sprintf(" LIMIT %v", o.Limit)
	}
	if o.Offset > 0 {
		s = s + fmt.Sprintf(" OFFSET %v", o.Offset)
	}
	return strings.TrimSpace(s)
}

func (r *RegExp) String() string {
	return r.Value
}

func (c *columns) String() string {
	if c.value == "" {
		return "*"
	} else {
		return c.value
	}
}

func (f *DataSource) String() string {
	parts := make([]string, 0, 3)
	if f.Database != "" {
		parts = append(parts, QuoteIdentifier(f.Database))
	}
	if f.RetentionPolicy != "" {
		parts = append(parts, QuoteIdentifier(f.RetentionPolicy))
	} else if f.Database != "" {
		parts = append(parts, "")
	}
	if f.Measurement != "" {
		parts = append(parts, QuoteIdentifier(f.Measurement))
	}
	return strings.Join(parts, ".")
}
