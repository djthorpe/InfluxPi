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

///////////////////////////////////////////////////////////////////////////////
// TYPES

type q_ShowDatabases struct{}

type q_CreateDatabase struct {
	database   string
	policyName string
	policy     *RetentionPolicy
}

type q_DropDatabase struct {
	database string
}

type q_DropRetentionPolicy struct {
	database string
	name     string
}

type q_AlterRetentionPolicy struct {
	database string
	name     string
	policy   *RetentionPolicy
	defalt   bool
}

type q_ShowRetentionPolicies struct {
	database string
}

type q_CreateRetentionPolicy struct {
	database string
	name     string
	policy   *RetentionPolicy
	defalt   bool
}

type q_ShowSeries struct {
	database    string
	measurement *Measurement
	limit       uint
	offset      uint
}

type q_Select struct {
	measurement []*Measurement
	where       []Predicate
	limit       uint
	offset      uint
}

type p_TagClause struct {
	name  string
	value []string
	op    string
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCT QUERIES

func ShowDatabases() Query {
	return &q_ShowDatabases{}
}

func ShowRetentionPolicies() Query {
	return &q_ShowRetentionPolicies{}
}

func ShowSeries() Query {
	return &q_ShowSeries{}
}

func CreateDatabase(name string) Query {
	return &q_CreateDatabase{database: name, policyName: "autogen"}
}

func CreateRetentionPolicy(database string, name string, policy *RetentionPolicy) Query {
	return &q_CreateRetentionPolicy{database: database, name: name, policy: policy}
}

func DropDatabase(name string) Query {
	return &q_DropDatabase{database: name}
}

func DropRetentionPolicy(database string, name string) Query {
	return &q_DropRetentionPolicy{database: database, name: name}
}

func AlterRetentionPolicy(database string, name string, policy *RetentionPolicy) Query {
	return &q_AlterRetentionPolicy{database: database, name: name, policy: policy}
}

func Select(measurements ...*Measurement) Query {
	return &q_Select{measurement: measurements}
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCT PREDICATES

func TagEquals(name string, value ...string) Predicate {
	return &p_TagClause{name: name, value: value, op: "="}
}

func TagNotEquals(name, value string) Predicate {
	return &p_TagClause{name: name, value: []string{value}, op: "!="}
}

func TagMatches(name, regexp string) Predicate {
	return &p_TagClause{name: name, value: []string{regexp}, op: "=~"}
}

///////////////////////////////////////////////////////////////////////////////
// SET DATABASE

func (q *q_CreateDatabase) Database(value string) Query        { q.database = value; return q }
func (q *q_DropDatabase) Database(value string) Query          { q.database = value; return q }
func (q *q_ShowSeries) Database(value string) Query            { q.database = value; return q }
func (q *q_ShowDatabases) Database(value string) Query         { return q }
func (q *q_ShowRetentionPolicies) Database(value string) Query { q.database = value; return q }
func (q *q_CreateRetentionPolicy) Database(value string) Query { q.database = value; return q }
func (q *q_DropRetentionPolicy) Database(value string) Query   { q.database = value; return q }
func (q *q_AlterRetentionPolicy) Database(value string) Query  { q.database = value; return q }
func (q *q_Select) Database(value string) Query                { return q }

///////////////////////////////////////////////////////////////////////////////
// SET RETENTION POLICY

func (q *q_CreateDatabase) RetentionPolicy(value *RetentionPolicy) Query        { q.policy = value; return q }
func (q *q_DropDatabase) RetentionPolicy(value *RetentionPolicy) Query          { return q }
func (q *q_ShowDatabases) RetentionPolicy(value *RetentionPolicy) Query         { return q }
func (q *q_ShowSeries) RetentionPolicy(value *RetentionPolicy) Query            { return q }
func (q *q_ShowRetentionPolicies) RetentionPolicy(value *RetentionPolicy) Query { return q }
func (q *q_CreateRetentionPolicy) RetentionPolicy(value *RetentionPolicy) Query {
	q.policy = value
	return q
}
func (q *q_DropRetentionPolicy) RetentionPolicy(value *RetentionPolicy) Query { return q }
func (q *q_AlterRetentionPolicy) RetentionPolicy(value *RetentionPolicy) Query {
	q.policy = value
	return q
}
func (q *q_Select) RetentionPolicy(value *RetentionPolicy) Query { return q }

///////////////////////////////////////////////////////////////////////////////
// SET DEFAULT

func (q *q_CreateDatabase) Default(value bool) Query        { return q }
func (q *q_DropDatabase) Default(value bool) Query          { return q }
func (q *q_ShowDatabases) Default(value bool) Query         { return q }
func (q *q_ShowSeries) Default(value bool) Query            { return q }
func (q *q_ShowRetentionPolicies) Default(value bool) Query { return q }
func (q *q_CreateRetentionPolicy) Default(value bool) Query { q.defalt = true; return q }
func (q *q_DropRetentionPolicy) Default(value bool) Query   { return q }
func (q *q_AlterRetentionPolicy) Default(value bool) Query  { q.defalt = true; return q }
func (q *q_Select) Default(value bool) Query                { return q }

///////////////////////////////////////////////////////////////////////////////
// SET OFFSET AND LIMIT

func (q *q_CreateDatabase) OffsetLimit(offset uint, limit uint) Query        { return q }
func (q *q_DropDatabase) OffsetLimit(offset uint, limit uint) Query          { return q }
func (q *q_ShowDatabases) OffsetLimit(offset uint, limit uint) Query         { return q }
func (q *q_ShowRetentionPolicies) OffsetLimit(offset uint, limit uint) Query { return q }
func (q *q_ShowSeries) OffsetLimit(offset uint, limit uint) Query {
	q.offset = offset
	q.limit = limit
	return q
}
func (q *q_CreateRetentionPolicy) OffsetLimit(offset uint, limit uint) Query { return q }
func (q *q_DropRetentionPolicy) OffsetLimit(offset uint, limit uint) Query   { return q }
func (q *q_AlterRetentionPolicy) OffsetLimit(offset uint, limit uint) Query  { return q }
func (q *q_Select) OffsetLimit(offset uint, limit uint) Query {
	q.offset = offset
	q.limit = limit
	return q
}

///////////////////////////////////////////////////////////////////////////////
// MEASUREMENT

func (q *q_CreateDatabase) Measurement(value ...*Measurement) Query        { return q }
func (q *q_DropDatabase) Measurement(value ...*Measurement) Query          { return q }
func (q *q_ShowDatabases) Measurement(value ...*Measurement) Query         { return q }
func (q *q_ShowRetentionPolicies) Measurement(value ...*Measurement) Query { return q }
func (q *q_CreateRetentionPolicy) Measurement(value ...*Measurement) Query { return q }
func (q *q_AlterRetentionPolicy) Measurement(value ...*Measurement) Query  { return q }
func (q *q_DropRetentionPolicy) Measurement(value ...*Measurement) Query   { return q }
func (q *q_ShowSeries) Measurement(value ...*Measurement) Query {
	if len(value) > 0 {
		q.measurement = value[0]
	} else {
		q.measurement = nil
	}
	return q
}
func (q *q_Select) Measurement(value ...*Measurement) Query {
	q.measurement = value
	return q
}

///////////////////////////////////////////////////////////////////////////////
// FILTER

func (q *q_CreateDatabase) Filter(value ...Predicate) Query        { return q }
func (q *q_DropDatabase) Filter(value ...Predicate) Query          { return q }
func (q *q_ShowDatabases) Filter(value ...Predicate) Query         { return q }
func (q *q_ShowRetentionPolicies) Filter(value ...Predicate) Query { return q }
func (q *q_CreateRetentionPolicy) Filter(value ...Predicate) Query { return q }
func (q *q_AlterRetentionPolicy) Filter(value ...Predicate) Query  { return q }
func (q *q_DropRetentionPolicy) Filter(value ...Predicate) Query   { return q }
func (q *q_ShowSeries) Filter(value ...Predicate) Query            { return q }
func (q *q_Select) Filter(value ...Predicate) Query {
	q.where = value
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *RetentionPolicy) query(name string) string {
	if p.Duration == 0 && p.ReplicationFactor == 0 && p.ShardGroupDuration == 0 {
		return ""
	}
	s := make([]string, 0, 4)
	if p.Duration != 0 {
		s = append(s, fmt.Sprintf("DURATION %v", p.Duration))
	}
	if p.ReplicationFactor != 0 {
		s = append(s, fmt.Sprintf("REPLICATION %v", p.ReplicationFactor))
	}
	if p.ShardGroupDuration != 0 {
		s = append(s, fmt.Sprintf("SHARD DURATION %v", p.ShardGroupDuration))
	}
	if name != "" {
		s = append(s, fmt.Sprintf("NAME %v", name))
	}
	return strings.Join(s, " ")
}

func (p *p_TagClause) String() string {
	if len(p.value) == 0 {
		panic("Expected at least one value")
	}
	switch p.op {
	case "=":
		if len(p.value) > 1 {
			return Quote(p.name) + " IN (TODO)"
		}
	case "=~":
		v := strings.Trim(p.value[0], "/")
		return Quote(p.name) + " " + p.op + " /" + v + "/"
	}
	return Quote(p.name) + " " + p.op + " " + QuoteString(p.value[0])
}

func (m Measurement) String() string {
	if m.Database == "" && m.Policy == "" {
		return Quote(m.Name)
	} else {
		return Quote(m.Database) + "." + Quote(m.Policy) + "." + Quote(m.Name)
	}
}

func (q *q_ShowDatabases) String() string {
	return "SHOW DATABASES"
}

func (q *q_ShowSeries) String() string {
	s := "SHOW SERIES"
	if len(q.database) > 0 {
		s = s + " ON " + Quote(q.database)
	}
	if q.measurement != nil {
		s = s + " FROM " + q.measurement.String()
	}
	if q.limit > 0 {
		s = s + " LIMIT " + fmt.Sprint(q.limit)
	}
	if q.offset > 0 {
		s = s + " OFFSET " + fmt.Sprint(q.offset)
	}
	return s
}

func (q *q_ShowRetentionPolicies) String() string {
	s := "SHOW RETENTION POLICIES"
	if len(q.database) > 0 {
		s = s + " ON " + Quote(q.database)
	}
	return s
}

func (q *q_DropRetentionPolicy) String() string {
	s := "DROP RETENTION POLICY " + Quote(q.name)
	if len(q.database) > 0 {
		s = s + " ON " + Quote(q.database)
	}
	return s
}

func (q *q_CreateRetentionPolicy) String() string {
	s := "CREATE RETENTION POLICY " + Quote(q.name)
	if q.database != "" {
		s = s + " ON " + Quote(q.database)
	}
	if q.policy != nil {
		if q.policy.ReplicationFactor == 0 {
			q.policy.ReplicationFactor = 1
		}
		if p := q.policy.query(""); p != "" {
			s = s + " " + p
		}
	}
	if q.defalt {
		s = s + " DEFAULT"
	}
	return s
}

func (q *q_AlterRetentionPolicy) String() string {
	s := "ALTER RETENTION POLICY " + Quote(q.name)
	if q.database != "" {
		s = s + " ON " + Quote(q.database)
	}
	if q.policy != nil {
		if p := q.policy.query(""); p != "" {
			s = s + " " + p
		}
	}
	if q.defalt {
		s = s + " DEFAULT"
	}
	return s
}

func (q *q_CreateDatabase) String() string {
	s := "CREATE DATABASE " + Quote(q.database)
	if q.policy != nil {
		if policy := q.policy.query(q.policyName); policy != "" {
			s = s + " WITH " + policy
		}
	}
	return s
}

func (q *q_DropDatabase) String() string {
	s := "DROP DATABASE " + Quote(q.database)
	return s
}

func (q *q_Select) String() string {
	s := "SELECT * FROM "
	for i, m := range q.measurement {
		s = s + m.String()
		if (i + 1) < len(q.measurement) {
			s = s + ","
		}
	}
	if len(q.where) > 0 {
		s = s + " WHERE "
	}
	if q.limit > 0 {
		s = s + " LIMIT " + fmt.Sprint(q.limit)
	}
	if q.offset > 0 {
		s = s + " OFFSET " + fmt.Sprint(q.offset)
	}
	return s
}
