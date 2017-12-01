// Connect to a remote InfluxDB instance and list the databases
// in the instance
package main

import (
	"fmt"
	"os"

	influx "github.com/djthorpe/InfluxPi"
	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func GetClientConfig(app *gopi.AppInstance) influx.Config {
	host, _ := app.AppFlags.GetString("host")
	port, _ := app.AppFlags.GetUint("port")
	ssl, _ := app.AppFlags.GetBool("ssl")
	username, _ := app.AppFlags.GetString("username")
	password, _ := app.AppFlags.GetString("password")
	timeout, _ := app.AppFlags.GetDuration("timeout")
	return influx.Config{
		Host:     host,
		Port:     port,
		SSL:      ssl,
		Username: username,
		Password: password,
		Timeout:  timeout,
	}
}

func RunLoop(app *gopi.AppInstance, done chan struct{}) error {
	db, exists := app.AppFlags.GetString("db")
	if client, err := gopi.Open(GetClientConfig(app), app.Logger); err != nil {
		return err
	} else {
		defer client.Close()
		if exists == false {
			return app.Logger.Error("Requires -db flag")
		}
		if err := client.(*influx.Client).DropDatabase(db); err != nil {
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
