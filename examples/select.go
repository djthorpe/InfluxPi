// Connect to a remote InfluxDB instance and query
// values
package main

import (
	"errors"
	"fmt"
	"os"

	tablewriter "github.com/djthorpe/InfluxPi/tablewriter"
	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
	influx "github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////

func GetClientConfig(app *gopi.AppInstance) influx.Config {
	host, _ := app.AppFlags.GetString("host")
	port, _ := app.AppFlags.GetUint("port")
	ssl, _ := app.AppFlags.GetBool("ssl")
	username, _ := app.AppFlags.GetString("username")
	password, _ := app.AppFlags.GetString("password")
	timeout, _ := app.AppFlags.GetDuration("timeout")
	precision, _ := app.AppFlags.GetString("precision")
	db, _ := app.AppFlags.GetString("db")
	return influx.Config{
		Host:      host,
		Port:      port,
		SSL:       ssl,
		Username:  username,
		Password:  password,
		Timeout:   timeout,
		Database:  db,
		Precision: precision,
	}
}

func RunLoop(app *gopi.AppInstance, done chan struct{}) error {
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

func registerFlags(config gopi.AppConfig) gopi.AppConfig {
	// Register the flags
	config.AppFlags.FlagString("host", "localhost", "InfluxDB hostname")
	config.AppFlags.FlagUint("port", 0, "InfluxDB port, or 0 to use the default")
	config.AppFlags.FlagBool("ssl", false, "Use SSL")
	config.AppFlags.FlagString("username", "", "InfluxDB username")
	config.AppFlags.FlagString("password", "", "InfluxDB password")
	config.AppFlags.FlagDuration("timeout", 0, "Connection timeout")
	config.AppFlags.FlagString("db", "", "Name of database to create")
	config.AppFlags.FlagString("measurement", "", "Measurement")
	config.AppFlags.FlagUint("limit", 0, "Limit number of rows returned")
	config.AppFlags.FlagUint("offset", 0, "Offset")
	config.AppFlags.FlagString("columns", "", "Columns to return")
	config.AppFlags.FlagString("precision", "", "Time precision")

	// Return config
	return config
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the application with an empty configuration
	app, err := gopi.NewAppInstance(registerFlags(gopi.NewAppConfig()))
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
