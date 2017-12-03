package influxdb_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/sys/logger"
	"github.com/djthorpe/influxdb"
	"github.com/djthorpe/influxdb/mock"
)

func TestOpen_000(t *testing.T) {
	configuration := mock.Config{}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else {
		t.Log(client)
	}
}
func TestOpen_001(t *testing.T) {
	db := "test"
	configuration := mock.Config{Database: db}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		t.Error("mock client does not implement all the required methods")
		_ = client.(influxdb.Driver)
	} else {
		t.Log(driver)
		/*if driver.GetDatabase() != db {
			t.Errorf("Expected database %v but got %v", db, driver.GetDatabase())
		}*/
	}
}

func TestOpen_002(t *testing.T) {
	db := "test"
	configuration := mock.Config{Database: db}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		t.Error("mock client does not implement all the required methods")
		_ = client.(influxdb.Driver)
	} else if err := driver.Close(); err != nil {
		t.Error(err)
	} else {
		return_val := driver.SetDatabase(db)
		if return_val != influxdb.ErrNotConnected {
			t.Error("expected ErrNotConnected response")
		}
	}
}

func TestQuote_000(t *testing.T) {
	tests := map[string]string{
		"":                "",
		"a":               "a",
		"t1":              "t1",
		"_":               "_",
		"0":               "\"0\"",
		"a0":              "a0",
		"test string":     "\"test string\"",
		"test \"string\"": "\"test \\\"string\\\"\"",
		"test \\string\"": "\"test \\\\string\\\"\"",
	}
	for k, expected := range tests {
		actual := influxdb.Quote(k)
		if actual != expected {
			t.Errorf("For [%v], expected [%v], got [%v]", k, expected, actual)
		}
	}
}

func TestQueries_003(t *testing.T) {
	query := influxdb.ShowDatabases()
	if query.String() != "SHOW DATABASES" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_004(t *testing.T) {
	query := influxdb.ShowRetentionPolicies()
	if query.String() != "SHOW RETENTION POLICIES" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_005(t *testing.T) {
	query := influxdb.ShowRetentionPolicies().Database("db")
	if query.String() != "SHOW RETENTION POLICIES ON db" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_007(t *testing.T) {
	query := influxdb.CreateDatabase("db")
	if query.String() != "CREATE DATABASE db" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_008(t *testing.T) {
	query := influxdb.CreateDatabase("db 2")
	if query.String() != "CREATE DATABASE \"db 2\"" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestCreateDatabase_001(t *testing.T) {
	db := "test"
	configuration := mock.Config{Database: db}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		t.Error("mock client does not implement all the required methods")
		_ = client.(influxdb.Driver)
	} else {
		defer driver.Close()
		if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		}
		if databases, err := driver.Do(influxdb.ShowDatabases()); err != nil {
			t.Error(err)
		} else {
			// TODO: Check databases
			t.Log(databases)
		}
	}
}

func TestShowDatabases_007(t *testing.T) {
	db := "test"
	configuration := mock.Config{Database: db}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		t.Error("mock client does not implement all the required methods")
		_ = client.(influxdb.Driver)
	} else {
		defer driver.Close()

		if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		}
	}
}
