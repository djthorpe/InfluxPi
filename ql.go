/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package influxdb

import (
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type q_ShowDatabases struct{}

type q_ShowRetentionPolicies struct {
	database string
}

type q_CreateDatabase struct {
	database   string
	policyName string
	policy     *RetentionPolicy
}

///////////////////////////////////////////////////////////////////////////////
// CONSTRUCTORS

func ShowDatabases() Query {
	return &q_ShowDatabases{}
}

func ShowRetentionPolicies() Query {
	return &q_ShowRetentionPolicies{}
}

func CreateDatabase(name string) Query {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	} else {
		return &q_CreateDatabase{database: name, policyName: "autogen"}
	}
}

///////////////////////////////////////////////////////////////////////////////
// SET DATABASE

func (q *q_CreateDatabase) Database(value string) Query {
	q.database = value
	return q
}

func (q *q_ShowDatabases) Database(value string) Query {
	// Database cannot be set for SHOW DATABASES, ignore
	return q
}

func (q *q_ShowRetentionPolicies) Database(value string) Query {
	q.database = value
	return q
}

///////////////////////////////////////////////////////////////////////////////
// SET RETENTION POLICY

func (q *q_CreateDatabase) RetentionPolicy(value *RetentionPolicy) Query {
	q.policy = value
	return q
}

func (q *q_ShowDatabases) RetentionPolicy(value *RetentionPolicy) Query {
	// Ignore
	return q
}

func (q *q_ShowRetentionPolicies) RetentionPolicy(value *RetentionPolicy) Query {
	// Ignore
	return q
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q *q_ShowDatabases) String() string {
	return "SHOW DATABASES"
}

func (q *q_ShowRetentionPolicies) String() string {
	s := "SHOW RETENTION POLICIES"
	if len(q.database) > 0 {
		s = s + " ON " + Quote(q.database)
	}
	return s
}

func (q *q_CreateDatabase) String() string {
	s := "CREATE DATABASE " + Quote(q.database)
	return s
}
