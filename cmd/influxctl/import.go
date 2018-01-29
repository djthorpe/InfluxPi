package influxctl

import (

	// frameworks
	"errors"

	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
)

////////////////////////////////////////////////////////////////////////////////

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
