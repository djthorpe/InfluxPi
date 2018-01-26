// InfluxDB command-line tool
package main

import (
	"errors"
	"os"

	// frameworks
	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"

	// modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/influxdb/v2"
)

const (
	MODULE_NAME = "influx/v2"
)

////////////////////////////////////////////////////////////////////////////////

func MainTask(app *gopi.AppInstance, done chan<- struct{}) error {
	// Create a client
	if client := app.ModuleInstance(MODULE_NAME).(influxdb.Client); client == nil {
		return errors.New("Missing module")
	} else {
		app.Logger.Info("client=%v", client)
	}

	// Signal to other tasks to end
	done <- gopi.DONE

	// Return success
	return nil
}

/*
	if client, err := gopi.Open(GetClientConfig(app), app.Logger); err != nil {
		return err
	} else {
		defer client.Close()

		measurement, _ := app.AppFlags.GetString("measurement")
		if measurement == "" {
			return errors.New("Missing measurement")
		}

		// Construct statement
		statement := client.(*influx.Client).Select(&influx.DataSource{Measurement: measurement})
		if limit, _ := app.AppFlags.GetUint("limit"); limit > 0 {
			statement = statement.Limit(limit)
		}
		if offset, _ := app.AppFlags.GetUint("offset"); offset > 0 {
			statement = statement.Offset(offset)
		}
		if columns, _ := app.AppFlags.GetString("columns"); columns != "" {
			statement = statement.Columns(columns)
		}
		if response, err := client.(*influx.Client).Do(statement); err != nil {
			return err
		} else if err := tablewriter.RenderASCII(response, os.Stdout); err != nil {
			return err
		}
	}

	done <- gopi.DONE
	return nil
}
*/

////////////////////////////////////////////////////////////////////////////////

func main() {
	config := gopi.NewAppConfig(MODULE_NAME)
	os.Exit(gopi.CommandLineTool(config, MainTask))
}
