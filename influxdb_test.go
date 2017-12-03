package influxdb_test

import (
	"testing"
	"time"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/sys/logger"
	"github.com/djthorpe/influxdb"
	"github.com/djthorpe/influxdb/mock"
	"github.com/djthorpe/influxdb/v2"
)

///////////////////////////////////////////////////////////////////////////////

// TODO: Change this when using actual driber to point to your influxdb hostname
func ServerHost() string {
	return "rpi3.lan"
}

func Driver(t *testing.T, db string) influxdb.Driver {
	return ActualDriver(t, db)
}

func MockDriver(t *testing.T, db string) influxdb.Driver {
	configuration := mock.Config{Database: db}
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		t.Fatal("mock client does not implement all the required methods")
	} else {
		return driver
	}
	return nil
}

func ActualDriver(t *testing.T, db string) influxdb.Driver {
	configuration := v2.Config{
		Database: db,
		Host:     ServerHost(),
	}

	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Error(err)
	} else if client, err := gopi.Open(configuration, log.(gopi.Logger)); err != nil {
		t.Error(err)
	} else if driver, ok := client.(influxdb.Driver); ok == false {
		_ = client.(influxdb.Driver)
		t.Fatal("v2 client does not implement all the required methods")
	} else {
		return driver
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func TestOpen_000(t *testing.T) {
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		t.Log(driver)
	}
}
func TestOpen_001(t *testing.T) {
	db := "_internal"
	if driver := Driver(t, db); driver == nil {
		t.Error("nil driver returned")
	} else {
		t.Log(driver)
		if driver.Database() != db {
			t.Errorf("Expected database %v but got %v", db, driver.Database())
		}
	}
}

func TestOpen_002(t *testing.T) {
	db := "_internal"
	if driver := Driver(t, db); driver == nil {
		t.Error("nil driver returned")
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

func TestQueries_009(t *testing.T) {
	policy := &influxdb.RetentionPolicy{}
	query := influxdb.CreateRetentionPolicy("db", "policy", policy)
	if query.String() != "CREATE RETENTION POLICY \"policy\" ON db REPLICATION 1" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_010(t *testing.T) {
	policy := &influxdb.RetentionPolicy{}
	query := influxdb.CreateRetentionPolicy("db", "policy", policy).Default(true)
	if query.String() != "CREATE RETENTION POLICY \"policy\" ON db REPLICATION 1 DEFAULT" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_011(t *testing.T) {
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 24,
	}
	query := influxdb.CreateRetentionPolicy("db", "policy", policy)
	if query.String() != "CREATE RETENTION POLICY \"policy\" ON db DURATION 24h0m0s REPLICATION 1" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_012(t *testing.T) {
	policy := &influxdb.RetentionPolicy{
		Duration:           time.Hour * 24,
		ShardGroupDuration: time.Hour * 48,
	}
	query := influxdb.CreateRetentionPolicy("db", "policy", policy)
	if query.String() != "CREATE RETENTION POLICY \"policy\" ON db DURATION 24h0m0s REPLICATION 1 SHARD DURATION 48h0m0s" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_013(t *testing.T) {
	query := influxdb.DropRetentionPolicy("db", "policy")
	if query.String() != "DROP RETENTION POLICY \"policy\" ON db" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_014(t *testing.T) {
	query := influxdb.ShowSeries()
	if query.String() != "SHOW SERIES" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_015(t *testing.T) {
	query := influxdb.ShowSeries().Database("db")
	if query.String() != "SHOW SERIES ON db" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_016(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").OffsetLimit(0, 10)
	if query.String() != "SHOW SERIES ON db LIMIT 10" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_017(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").OffsetLimit(10, 0)
	if query.String() != "SHOW SERIES ON db OFFSET 10" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_018(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").OffsetLimit(10, 10)
	if query.String() != "SHOW SERIES ON db LIMIT 10 OFFSET 10" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_019(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").Measurement(&influxdb.Measurement{Name: "test"})
	if query.String() != "SHOW SERIES ON db FROM test" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_020(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").Measurement(&influxdb.Measurement{Name: "test", Database: "db"})
	if query.String() != "SHOW SERIES ON db FROM db..test" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_021(t *testing.T) {
	query := influxdb.ShowSeries().Database("db").Measurement(&influxdb.Measurement{Name: "test", Database: "db", Policy: "policy"})
	if query.String() != "SHOW SERIES ON db FROM db.\"policy\".test" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_022(t *testing.T) {
	query := influxdb.AlterRetentionPolicy("db", "policy", nil)
	if query.String() != "ALTER RETENTION POLICY \"policy\" ON db" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_023(t *testing.T) {
	query := influxdb.AlterRetentionPolicy("db", "policy", nil).Default(true)
	if query.String() != "ALTER RETENTION POLICY \"policy\" ON db DEFAULT" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_024(t *testing.T) {
	query := influxdb.AlterRetentionPolicy("db", "policy", &influxdb.RetentionPolicy{Duration: time.Hour})
	if query.String() != "ALTER RETENTION POLICY \"policy\" ON db DURATION 1h0m0s" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_025(t *testing.T) {
	query := influxdb.AlterRetentionPolicy("db", "policy", &influxdb.RetentionPolicy{ReplicationFactor: 2})
	if query.String() != "ALTER RETENTION POLICY \"policy\" ON db REPLICATION 2" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_026(t *testing.T) {
	query := influxdb.AlterRetentionPolicy("db", "policy", &influxdb.RetentionPolicy{ShardGroupDuration: time.Minute})
	if query.String() != "ALTER RETENTION POLICY \"policy\" ON db SHARD DURATION 1m0s" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_027(t *testing.T) {
	query := influxdb.Select(&influxdb.Measurement{Name: "test"})
	if query.String() != "SELECT * FROM test" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestQueries_028(t *testing.T) {
	query := influxdb.Select(&influxdb.Measurement{Name: "test1"}, &influxdb.Measurement{Name: "test2"})
	if query.String() != "SELECT * FROM test1,test2" {
		t.Errorf("Unexpected query: %v", query.String())
	}
}

func TestCreateDatabase_001(t *testing.T) {
	db := "TestCreateDatabase_001"
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if databases, err := driver.Do(influxdb.ShowDatabases()); err != nil {
			t.Error(err)
		} else if values, err := databases.Column(0, "databases", "name"); err != nil {
			t.Error(err)
		} else {
			// Check database was created
			found := false
			for _, v := range values {
				if v == db {
					found = true
				}
			}
			if found == false {
				t.Error("database", db, "wasn't created")
			}
		}
	}
}

func TestCreateDatabase_002(t *testing.T) {
	db := "TestCreateDatabase_002"
	policy := &influxdb.RetentionPolicy{}
	// Empty policy
	if q := influxdb.CreateDatabase(db).RetentionPolicy(policy); q.String() != "CREATE DATABASE "+db {
		t.Error("Unexpected query:", q.String())
	}
	// Policy with replication factor
	policy.ReplicationFactor = 1
	if q := influxdb.CreateDatabase(db).RetentionPolicy(policy); q.String() != "CREATE DATABASE "+db+" WITH REPLICATION 1 NAME autogen" {
		t.Error("Unexpected query:", q.String())
	}
	// Policy with duration
	policy.ReplicationFactor = 0
	policy.Duration = 1 * time.Second
	if q := influxdb.CreateDatabase(db).RetentionPolicy(policy); q.String() != "CREATE DATABASE "+db+" WITH DURATION 1s NAME autogen" {
		t.Error("Unexpected query:", q.String())
	}
	// Policy with group shard duration
	policy.ReplicationFactor = 0
	policy.Duration = 1 * time.Second
	policy.ShardGroupDuration = 1 * time.Minute
	if q := influxdb.CreateDatabase(db).RetentionPolicy(policy); q.String() != "CREATE DATABASE "+db+" WITH DURATION 1s SHARD DURATION 1m0s NAME autogen" {
		t.Error("Unexpected query:", q.String())
	}
}

func TestCreateDatabase_003(t *testing.T) {
	db := "TestCreateDatabase_003"
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 1,
	}
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, policy); err != nil {
			t.Error(err)
		} else if databases, err := driver.Do(influxdb.ShowDatabases()); err != nil {
			t.Error(err)
		} else if values, err := databases.Column(0, "databases", "name"); err != nil {
			t.Error(err)
		} else {
			// Check database was created
			found := false
			for _, v := range values {
				if v == db {
					found = true
				}
			}
			if found == false {
				t.Error("database", db, "wasn't created")
			}
		}
	}
}

func TestCreateDatabase_004(t *testing.T) {
	db := "TestCreateDatabase_004"
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if err := driver.SetDatabase(db); err != nil {
			t.Error(err)
		} else {
			// Create retention policy
			if err := driver.CreateRetentionPolicy("policy", &influxdb.RetentionPolicy{
				Duration: time.Hour * 5,
			}); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestCreateRetentionPolicy_001(t *testing.T) {
	db := "TestCreateRetentionPolicy_001"
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 1,
	}
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if err := driver.SetDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateRetentionPolicy("policy", policy); err != nil {
			t.Error(err)
		} else if err := driver.CreateRetentionPolicy("policy", policy); err != nil {
			if err != influxdb.ErrAlreadyExists {
				t.Error(err)
			}
		}
	}
}

func TestCreateRetentionPolicy_002(t *testing.T) {
	db := "TestCreateRetentionPolicy_002"
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 1,
	}
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if err := driver.SetDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateRetentionPolicy("policy", policy); err != nil {
			t.Error(err)
		} else if err := driver.DropRetentionPolicy("policy"); err != nil {
			t.Error(err)
		}
	}
}

func TestCreateRetentionPolicy_003(t *testing.T) {
	db := "TestCreateRetentionPolicy_003"
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 1,
	}
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if err := driver.SetDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateRetentionPolicy("policy", policy); err != nil {
			t.Error(err)
		} else if policies, err := driver.RetentionPolicies(); err != nil {
			t.Error(err)
		} else if policy, exists := policies["policy"]; exists == false {
			t.Error("Missing policy after being created")
		} else if policy.Duration != time.Hour*1 {
			t.Error("Invalid policy time, unexpected value %v", policy.Duration)
		}
	}
}

func TestDropRetentionPolicy_001(t *testing.T) {
	db := "TestDropRetentionPolicy_001"
	policy := &influxdb.RetentionPolicy{
		Duration: time.Hour * 1,
	}
	if driver := Driver(t, ""); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if err := driver.DropDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateDatabase(db, nil); err != nil {
			t.Error(err)
		} else if err := driver.SetDatabase(db); err != nil {
			t.Error(err)
		} else if err := driver.CreateRetentionPolicy("policy", policy); err != nil {
			t.Error(err)
		} else if policies, err := driver.RetentionPolicies(); err != nil {
			t.Error(err)
		} else if _, exists := policies["policy"]; exists == false {
			t.Error("Missing policy after being created")
		} else if err := driver.DropRetentionPolicy("policy"); err != nil {
			t.Error(err)
		} else if policies, err := driver.RetentionPolicies(); err != nil {
			t.Error(err)
		} else if _, exists := policies["policy"]; exists == true {
			t.Error("Drop retention policy failed")
		}
	}
}

func TestSelect_001(t *testing.T) {
	db := "_internal"
	if driver := Driver(t, db); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if response, err := driver.Do(influxdb.Select(&influxdb.Measurement{Name: "httpd"}).OffsetLimit(0, 10)); err != nil {
			t.Error(err)
		} else {
			if len(response) != 1 {
				t.Error("Expected one result from query, got", len(response))
			}
		}
	}
}

func TestSelect_002(t *testing.T) {
	db := "_internal"
	if driver := Driver(t, db); driver == nil {
		t.Error("nil driver returned")
	} else {
		defer driver.Close()

		if response, err := driver.Do(influxdb.Select(&influxdb.Measurement{Name: "httpd"}, &influxdb.Measurement{Name: "shard"}).OffsetLimit(0, 10)); err != nil {
			t.Error(err)
		} else {
			if len(response) != 2 {
				t.Error("Expected one result from query, got", len(response))
			}
			if response[0].Name != "httpd" {
				t.Error("Expected result 0 to be httpd, got", response[0].Name)
			}
			if response[1].Name != "shard" {
				t.Error("Expected result 1 to be shard, got", response[0].Name)
			}
		}
	}
}

func TestWhere_001(t *testing.T) {
	if where := influxdb.TagEquals("name", "value"); where.String() != "name = \"value\"" {
		t.Error("Expected string, got", where.String())
	}
	if where := influxdb.TagNotEquals("name", "value"); where.String() != "name != \"value\"" {
		t.Error("Expected string, got", where.String())
	}
	if where := influxdb.TagNotEquals("name with space", "value"); where.String() != "\"name with space\" != \"value\"" {
		t.Error("Expected string, got", where.String())
	}
	if where := influxdb.TagNotEquals("name", "\"value\""); where.String() != "name != \"\\\"value\\\"\"" {
		t.Error("Expected string, got", where.String())
	}
	if where := influxdb.TagEquals("name", "a", "b"); where.String() != "name IN (\"a\",\"b\")" {
		t.Error("Expected string, got", where.String())
	}
}
