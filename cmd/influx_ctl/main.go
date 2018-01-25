// InfluxDB command-line tool
package main

import (
	"os"

	// frameworks
	gopi "github.com/djthorpe/gopi"
	influx "github.com/djthorpe/influxdb"

	// modules
	_ "github.com/djthorpe/gopi/sys/logger"
	v2 "github.com/djthorpe/influxdb/v2"
)

////////////////////////////////////////////////////////////////////////////////
/*
func GetClientConfig(app *gopi.AppInstance) influx.Config {
	return
}
*/

func GetClient(app *gopi.AppInstance) (influx.Driver, error) {
	host, _ := app.AppFlags.GetString("influx.host")
	port, _ := app.AppFlags.GetUint("influx.port")
	ssl, _ := app.AppFlags.GetBool("influx.ssl")
	username, _ := app.AppFlags.GetString("influx.user")
	password, _ := app.AppFlags.GetString("influx.pass")
	timeout, _ := app.AppFlags.GetDuration("influx.timeout")
	db, _ := app.AppFlags.GetString("influx.db")

	if client, err := gopi.Open(v2.Config{
		Host:     host,
		Port:     port,
		SSL:      ssl,
		Username: username,
		Password: password,
		Timeout:  timeout,
		Database: db,
	}, app.Logger); err != nil {
		return nil, err
	} else {
		return client.(influx.Driver), nil
	}
}

func MainTask(app *gopi.AppInstance, done chan<- struct{}) error {
	// Create a client
	if client, err := GetClient(app); err != nil {
		return err
	} else {
		defer client.Close()
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

func registerFlags(config gopi.AppConfig) gopi.AppConfig {
	// Register the flags
	config.AppFlags.FlagString("influx.host", "localhost", "hostname")
	config.AppFlags.FlagUint("influx.port", 0, "port")
	config.AppFlags.FlagBool("influx.ssl", false, "Use SSL")
	config.AppFlags.FlagString("influx.user", "", "Username")
	config.AppFlags.FlagString("influx.pass", "", "Password")
	config.AppFlags.FlagDuration("influx.timeout", 0, "Connection timeout")
	config.AppFlags.FlagString("influx.db", "", "Database")
	/*	config.AppFlags.FlagString("measurement", "", "Measurement")
		config.AppFlags.FlagUint("limit", 0, "Limit number of rows returned")
		config.AppFlags.FlagUint("offset", 0, "Offset")
		config.AppFlags.FlagString("columns", "", "Columns to return")
		config.AppFlags.FlagString("precision", "", "Time precision")
	*/

	// Return config
	return config
}

func main() {
	config := registerFlags(gopi.NewAppConfig())
	os.Exit(gopi.CommandLineTool(config, MainTask))
}
