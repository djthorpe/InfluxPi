package influxctl

import (
	"errors"
	"os"

	// frameworks
	gopi "github.com/djthorpe/gopi"
	"github.com/djthorpe/influxdb"
	"github.com/djthorpe/influxdb/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func ListDatabases(client influxdb.Client, app *gopi.AppInstance) error {
	// Return a table of databases
	q := influxdb.ShowDatabases()
	if r, err := client.Do(q); err != nil {
		return err
	} else {
		for _, dataset := range r {
			tablewriter.RenderASCII(dataset, os.Stdout)
		}
		return nil
	}
}

func CreateDatabase(client influxdb.Client, app *gopi.AppInstance) error {
	// Set database
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.CreateDatabase(db, nil); err != nil {
		return err
	}

	return ListDatabases(client, app)
}

func DropDatabase(client influxdb.Client, app *gopi.AppInstance) error {
	// Get database flag
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.DropDatabase(db); err != nil {
		return err
	}

	return ListDatabases(client, app)
}

func CreateRetentionPolicy(client influxdb.Client, app *gopi.AppInstance) error {
	// Get database flag and policy name
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if policy_name, err := GetOneArg(app, "Policy Name"); err != nil {
		return err
	} else if policy, err := GetPolicyValue(app); err != nil {
		return err
	} else if err := client.SetDatabase(db); err != nil {
		return err
	} else if err := client.CreateRetentionPolicy(policy_name, policy); err != nil {
		return err
	} else {
		return ListRetentionPolicies(client, app)
	}
}

func DropRetentionPolicy(client influxdb.Client, app *gopi.AppInstance) error {
	// Get database flag and policy name
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if policy_name, err := GetOneArg(app, "Policy Name"); err != nil {
		return err
	} else if err := client.SetDatabase(db); err != nil {
		return err
	} else if err := client.DropRetentionPolicy(policy_name); err != nil {
		return err
	} else {
		return ListRetentionPolicies(client, app)
	}
	return nil
}

func ListRetentionPolicies(client influxdb.Client, app *gopi.AppInstance) error {
	// Set database
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.SetDatabase(db); err != nil {
		return err
	}

	// Return a table of retention policies
	q := influxdb.ShowRetentionPolicies()
	if r, err := client.Do(q); err != nil {
		return err
	} else {
		for _, dataset := range r {
			tablewriter.RenderASCII(dataset, os.Stdout)
		}
		return nil
	}
}

func ListSeries(client influxdb.Client, app *gopi.AppInstance) error {
	// Set database
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.SetDatabase(db); err != nil {
		return err
	}

	// Return a table of series
	q := influxdb.ShowSeries()
	if r, err := client.Do(q); err != nil {
		return err
	} else {
		for _, dataset := range r {
			tablewriter.RenderASCII(dataset, os.Stdout)
		}
		return nil
	}
}

func ListMeasurements(client influxdb.Client, app *gopi.AppInstance) error {
	// Set database
	db, _ := app.AppFlags.GetString("db")
	if db == "" {
		return errors.New("-db flag required")
	} else if err := client.SetDatabase(db); err != nil {
		return err
	}

	// Return a table of measurements
	q := influxdb.ShowMeasurements()
	if r, err := client.Do(q); err != nil {
		return err
	} else {
		for _, dataset := range r {
			tablewriter.RenderASCII(dataset, os.Stdout)
		}
		return nil
	}
}
