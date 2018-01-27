/*
	InfluxDB client
	(c) Copyright David Thorpe 2017
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE file
*/

package v2

import (
	gopi "github.com/djthorpe/gopi"
	influxdb "github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////
// MODULE INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name: "influx/v2",
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("influx.host", "localhost", "Host")
			config.AppFlags.FlagUint("influx.port", influxdb.DefaultPortHTTP, "Port")
			config.AppFlags.FlagBool("influx.ssl", false, "Use SSL")
			config.AppFlags.FlagBool("influx.ssl.verify", true, "Verify SSL Certificate")
			config.AppFlags.FlagString("influx.user", "", "User")
			config.AppFlags.FlagString("influx.password", "", "Password")
			config.AppFlags.FlagDuration("influx.timeout", 0, "Communication timeout")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			host, _ := app.AppFlags.GetString("influx.host")
			port, _ := app.AppFlags.GetUint("influx.port")
			ssl, _ := app.AppFlags.GetBool("influx.ssl")
			sslverify, _ := app.AppFlags.GetBool("influx.ssl.verify")
			user, _ := app.AppFlags.GetString("influx.user")
			password, _ := app.AppFlags.GetString("influx.password")
			timeout, _ := app.AppFlags.GetDuration("influx.timeout")
			return gopi.Open(Config{
				Host:      host,
				Port:      port,
				SSL:       ssl,
				SSLVerify: sslverify,
				Username:  user,
				Password:  password,
				Timeout:   timeout,
			}, app.Logger)
		},
	})
}
