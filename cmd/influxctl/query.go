package main

import (

	// frameworks
	"errors"
	"os"

	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
	"github.com/djthorpe/influxdb/tablewriter"
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
