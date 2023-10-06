package main

import (
	"io"
	"os"
	"flag"
	"fmt"
	"path"

	// Application package
	"github.com/rsdoiel/skimmer"
)

var (
	helpText := `

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS]

# DESCRIPTION

{app_name} is a extremely lightweight feed reader. It does two things. Retrieves
feeds from a list of urls and stores the contents in a SQLite 3 database. 
It displays the harvested content in reverse chronological order. The content
display is not "paged" so typically you would use {app_name} in conjunction with
the POSIX command 'more' or GNU command 'less'.

The url list file uses the newsboat format and {app_name} looks in the 
`+"`"+`$HOME/.newsboat/urls`+"`"+` directory for that file.

A note about the display. RSS 2 items do not require titles. For those items
an "@<LABEL>" is printed before displaying the description element. 

# OPTIONS

-help
: display help

-license
: display license

-version
: display version and release hash

-fetch
: download the latest feed content

-read
: display the downloaded contents

# EXAMPLE

{app_name} perform a "fetch" and "read" action based on the contents
of your newsboat url list.

~~~
{app_name}
~~~

Using the fetch option the url will be read and items will be downloaded
into the database. No read option will be performed. This is so you can
run this on a cron so the database is updated as frequently as you like.

~~~
{app_name} -fetch
~~~

Using the read option the contents of the database are displayed 
in reverse chronological order. No fetch is performed.

~~~
{app_name} -read
~~~

# Acknowledgements

This experiment would not be possible without the authors of
newsboat, sqlite3, Pandoc and mmcdole's gofeed pacakge.

`
)

func main() {
	appName := path.Base(os.Argv[0])
	showHelp, showVersion, showLicense := false, false, false
	fetch, read := false, false
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showVersion, "version", showVeresion, "display version")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&fetch, "fetch", fetch, "import feed content into database")
	flag.BoolVar(&read, "read", read, "display database content in reverse chronologically order")
	flag.Parse()

	args := flag.Args()

	var err error

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprint(out, FmtHelp(helpText, appName, skimmer.Version, skimmer.ReleaseDate, skimmer.ReleaseHash)
		os.Exit(0)
	}

	if showLicense {
		flag.Fprintf(out, "%s\n", skimmer.LicenseText)
		os.Exit(0)
	}

	if showVersion {
		flag.Fprint(out, "%s %s %s\n", appName, skimmer.Version, skimmer.ReleaseHash)
	}

	app := new(skimmer.Skimmer)
	app.Fetch := fetch
	app.Read := read
	if err := app.Run(in, out, eout, args); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
