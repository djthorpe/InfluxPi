/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

// Connect to a remote InfluxDB instance and list the databases in the instance
package main

import (
	"errors"
	"fmt"
	"os"

	// Interfaces
	gopi "github.com/djthorpe/gopi"
	influxdb "github.com/djthorpe/influxdb"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/influxdb/v2"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *gopi.AppInstance, done chan struct{}) error {
	if client, ok := app.Module("influxdb/v2").(influxdb.Driver); ok == false {
		return errors.New("No influxdb driver")
	}

	// Output version
	fmt.Println("     VERSION:", client.Version())

	// Successful completion
	done <- gopi.DONE
	return nil
}

/*
		// Retrieve databases
		if databases, err := client.(*influx.Client).ShowDatabases(); err != nil {
			return err
		} else {
			fmt.Println("   DATABASES:", databases)
		}

		// Retrieve measurements if database is set
		if measurements, err := client.(*influx.Client).GetMeasurements(); err != nil && err != influx.ErrEmptyResponse {
			return err
		} else {
			fmt.Println("MEASUREMENTS:", measurements)
		}

		// Retrieve retention policies
		if policies, err := client.(*influx.Client).GetRetentionPolicies(); err != nil && err != influx.ErrEmptyResponse {
			return err
		} else {
			fmt.Println("RETENTION POLICIES:", policies)
		}

		// Show series
		if series, err := client.(*influx.Client).ShowSeries(); err != nil {
			return err
		} else {
			fmt.Println("SERIES:", series)
		}

		// Show measurements
		if measurements, err := client.(*influx.Client).ShowMeasurements(nil, nil); err != nil {
			return err
		} else {
			fmt.Println("MEASUREMENTS:", measurements)
		}

	}
*/

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP THE APPLICATION

func registerFlags(config gopi.AppConfig) gopi.AppConfig {
	// Register the flags & return the configuration
	return config
}

func main_inner() int {
	// Set application configuration
	config := gopi.NewAppConfig("influxdb/v2")
	// Create the application with an empty configuration
	app, err := gopi.NewAppInstance(registerFlags(config))
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(RunLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
