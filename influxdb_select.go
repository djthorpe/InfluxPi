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

type Offset struct {
	Limit  uint
	Offset uint
}

type DataSource struct {
	Measurement     string
	Database        string
	RetentionPolicy string
}

type s struct {
	offset      *Offset
	datasources []*DataSource
	columns     []*Column
}

type Statement interface {
	// Set offset and limit
	Offset(uint) Statement
	Limit(uint) Statement

	// Write out statement
	Statement() string
}

type Column struct{}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTORS

// Select returns a select statement
func (this *Client) Select(from ...*DataSource) Statement {
	if len(from) == 0 {
		this.log.Error("Call to Select requires at least one data source")
		return nil
	}
	return &s{
		datasources: from,
	}
}

func (this *Client) Do(statement Statement) error {
	if this.client == nil {
		return ErrNotConnected
	}
	// Execute query
	if response, err := this.Query(statement.Statement()); err != nil {
		return err
	} else {
		fmt.Println(response)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STATEMENT IMPLEMENTATION FOR SELECT

func (this *s) Limit(limit uint) Statement {
	if this.offset == nil {
		this.offset = &Offset{Limit: limit}
	} else {
		this.offset.Limit = limit
	}
	return this
}

func (this *s) Offset(offset uint) Statement {
	if this.offset == nil {
		this.offset = &Offset{Offset: offset}
	} else {
		this.offset.Offset = offset
	}
	return this
}

func (this *s) Statement() string {
	q := "SELECT "
	// COLUMNS
	if len(this.columns) > 0 {
		for i := range this.columns {
			q = q + this.columns[i].String() + ","
		}
		q = strings.TrimSuffix(q, ",")
	} else {
		q = q + "*"
	}
	// DATA SOURCES
	if len(this.datasources) > 0 {
		q = q + " FROM "
		for i := range this.datasources {
			q = q + this.datasources[i].String() + ","
		}
		q = strings.TrimSuffix(q, ",")
	}
	// LIMIT and OFFSET
	if this.offset != nil {
		q = q + " " + this.offset.String()
	}

	return q
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (o *Offset) String() string {
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

func (c *Column) String() string {
	return "*"
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
