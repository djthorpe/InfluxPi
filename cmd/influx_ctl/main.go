// InfluxDB command-line tool
package main

import (
	"errors"
	"fmt"
	"os"

	// frameworks
	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
	"github.com/djthorpe/influxdb/tablewriter"

	// modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/influxdb/v2"
)

const (
	MODULE_NAME = "influx/v2"
)

////////////////////////////////////////////////////////////////////////////////

type CommandFunc func(client influxdb.Client, app *gopi.AppInstance) error

var (
	Commands = map[string]CommandFunc{
		"Databases":      ListDatabases,
		"CreateDatabase": CreateDatabase,
		"DropDatabase":   DropDatabase,
		"Policies":       ListRetentionPolicies,
		"CreatePolicy":   CreateRetentionPolicy,
		"DropPolicy":     DropRetentionPolicy,
		"Series":         ListSeries,
		"Measurements":   ListMeasurements,
		"Query":          Query,
		"Import":         Import,
	}
)

////////////////////////////////////////////////////////////////////////////////

func Query(client influxdb.Client, app *gopi.AppInstance) error {
	// Get flags
	db, _ := app.AppFlags.GetString("db")
	offset, _ := app.AppFlags.GetUint("offset")
	limit, _ := app.AppFlags.GetUint("limit")

	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.SetDatabase(db); err != nil {
		return err
	} else if measurement, err := GetOneArg(app, "Measurement"); err != nil {
		return err
	} else {
		q := influxdb.Select(GetMeasurement(measurement)).OffsetLimit(offset, limit)
		if r, err := client.Do(q); err != nil {
			return err
		} else {
			for _, dataset := range r {
				tablewriter.RenderASCII(dataset, os.Stdout)
			}
			return nil
		}
	}
}

func Import(client influxdb.Client, app *gopi.AppInstance) error {
	// Get flags
	db, _ := app.AppFlags.GetString("db")

	// Select database, retrieve measurement name
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.SetDatabase(db); err != nil {
		return err
	} else if measurement, err := GetOneArg(app, "Measurement"); err != nil {
		return err
	} else if dataset, err := client.NewDataset(measurement, []string{"tag1", "tag2"}, []string{"field1", "field2"}); err != nil {
		return err
	} else {
		if err := client.Write(dataset); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetOneArg(app *gopi.AppInstance, param1 string) (string, error) {
	if args := app.AppFlags.Args(); len(args) < 2 {
		return "", fmt.Errorf("Missing \"%v\" command-line argument", param1)
	} else if len(args) > 2 {
		return "", fmt.Errorf("Too many command-line arguments")
	} else {
		return args[1], nil
	}
}

func GetPolicyValue(app *gopi.AppInstance) (*influxdb.RetentionPolicy, error) {
	return &influxdb.RetentionPolicy{}, nil
}

func GetMeasurement(arg string) *influxdb.Measurement {
	return &influxdb.Measurement{
		Name: arg,
	}
}

////////////////////////////////////////////////////////////////////////////////

func MainTask(app *gopi.AppInstance, done chan<- struct{}) error {
	// Call command
	if args := app.AppFlags.Args(); len(args) < 1 {
		return gopi.ErrHelp
	} else if c, ok := Commands[args[0]]; ok == false {
		return errors.New("Invalid command")
	} else if client := app.ModuleInstance(MODULE_NAME).(influxdb.Client); client == nil {
		return errors.New("Missing module")
	} else if err := c(client, app); err != nil {
		return err
	}

	// Signal to other tasks to end
	done <- gopi.DONE

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Configuration
	config := gopi.NewAppConfig(MODULE_NAME)
	config.AppFlags.FlagString("db", "", "Database name")
	config.AppFlags.FlagUint("limit", 1000, "Row limit")
	config.AppFlags.FlagUint("offset", 0, "Row offset")

	// Run Command-Line Tool
	os.Exit(gopi.CommandLineTool(config, MainTask))
}
