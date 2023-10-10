package skimmer

import (
	"io"
	"fmt"
	"database/sql"
)

// Skim2Md supports the skim2md cli.
type Skim2Md struct {
	// AppName holds the name of the application
	AppName string `json:"app_name,omitempty"`

	// DbName holds the path to the SQLite3 database
	DBName string `json:"db_name,omitempty"`

	out io.Writer
	eout io.Writer
}

// NewSkim2Md initialized a new Skim2Md struct
func NewSkim2Md(appName string) (*Skim2Md, error) {
	app := new(Skim2Md)
	app.AppName = appName
	return app, nil
}

func (app *Skim2Md) DisplayItem(out io.Writer, link string, title string, description string, updated string, published string, label string, tags string) error {
	// Then see about formatting things.
	pressTime := published
	if updated != "" {
		pressTime = updated
	}
	if len(pressTime) > 10 {
		pressTime = pressTime[0:10]
	}
	if title == "" {
		title = fmt.Sprintf("**@%s** (date: %s)", label, pressTime)
	} else {
		title = fmt.Sprintf("## %s\n\ndate: %s", title, pressTime)
	}
	fmt.Fprintf(out, `--

%s

%s 

<%s>

%s
`, title, description, link, tags)
	return nil
}


// Write, display the contents from database
func (app *Skim2Md) Write(db *sql.DB) error {
	stmt := SQLDisplayItems
	rows, err := db.Query(stmt, "saved")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			link        string
			title       string
			description string
			updated     string
			published   string
			label       string
			tags        string
		)
		if err := rows.Scan(&link, &title, &description, &updated, &published, &label, &tags); err != nil {
			fmt.Fprint(app.eout, "%s\n", err)
			continue
		}
		if err := app.DisplayItem(app.out, link, title, description, updated, published, label, tags); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (app *Skim2Md) Run(out io.Writer, eout io.Writer, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing skimmer database file")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected only one skimmer database file")
	}
	dsn := args[0]
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := app.Write(db); err != nil {
		return err	
	}
	return nil
}
