// Connect to a remote InfluxDB instance and list the databases
// in the instance
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/influxdb/v2"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *gopi.AppInstance, done chan struct{}) error {
	return nil
}

/*
	host, _ := app.AppFlags.GetString("host")
	port, _ := app.AppFlags.GetUint("port")
	ssl, _ := app.AppFlags.GetBool("ssl")
	username, _ := app.AppFlags.GetString("username")
	password, _ := app.AppFlags.GetString("password")
	timeout, _ := app.AppFlags.GetDuration("timeout")
	db, _ := app.AppFlags.GetString("db")
	if client, err := gopi.Open(influx.Config{
		Host:     host,
		Port:     port,
		SSL:      ssl,
		Username: username,
		Password: password,
		Timeout:  timeout,
		Database: db,
	}, app.Logger); err != nil {
		return err
	} else {
		defer client.Close()
		app.Logger.Debug("influxdb=%v", client)

		// Output version
		fmt.Println("     VERSION:", client.(*influx.Client).GetVersion())

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

	done <- gopi.DONE
	return nil
}
*/

func registerFlags(config gopi.AppConfig) gopi.AppConfig {
	// Register the flags
	// Return config
	return config
}

////////////////////////////////////////////////////////////////////////////////

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
