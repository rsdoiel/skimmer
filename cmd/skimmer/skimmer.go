package main

import (
	"os"
	"flag"
	"fmt"
	"path"

	// Application package
	"github.com/rsdoiel/skimmer"
)

var (
	helpText = `%{app_name}(1) skimmer user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS]

# DESCRIPTION

{app_name} is a extremely lightweight feed reader. It can do two things. Retrieve
feed items from a list of urls and store the contents in a SQLite 3 database
(see -fetch option). It can write the contents of the SQLite 3 database to 
standard output (the default behavior if no options provided).

{app_name} displays the harvested content in reverse chronological order. The content
displayed is not "paged" so typically you would use {app_name} in conjunction with
the POSIX command 'more' or GNU command 'less'.

The url list file uses the newsboat url file format and {app_name} looks in the 
`+"`"+`$HOME/.skimmer/skimmer.urls`+"`"+` directory for that file.

By default (i.e. no options specified) {app_name} just outputs the 
contents it find in the database. If you use the "fetch" it will download
the contents identified in the feeds an update the database. Pulling things
from the web can be slow so you need to explicitly invoke the "fetch" option
to do this.

The output format uses Pandoc's style of markdown markup. 
- "--" starts the item record
- This is followed by "##" (aka H2), followed by a date 
(updated or published), followed by title or if none an
"@<LABEL>" where LABEL is the feed's title. 
- Next is the link line (link is in angle brackets)
- Finally this is wrapped up by the description content.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version and release hash

-fetch
: download the latest feed content

-display
: display the downloaded contents

-limit N
: display the N most recent items.

-prune TIMESTAMP
: Remove items in database older then TIMESTAMP. If TIMESTAMP "today" the
it removes items before today's date, if it is "now" prunes everything with
a timestamp before the current time. Timestamps can be day or day hour in
the YYYY-MM-DD or YYYY-MM-DD HH:MM:SS formats.

# EXAMPLE

{app_name} update your feed database and read your feed items
and page items with "less -R".

~~~
{app_name] -fetch
{app_name} | less -R
~~~

The fetch read the url file and download any feed items found.

~~~
{app_name} -fetch
~~~

Display the current contents of feed item database in reverse chronological order.

~~~
{app_name} -display
~~~

{app_name} can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
{app_name} -prune today -limit 5
~~~

If I limit the number of items I am reading to about 100 or so I've found
that this combination works nice.

~~~
{app_name} -limit 100 | pandoc -f markdown -t plain | less -R
~~~


# Acknowledgments

This project is an experiment it would not be possible without the authors of
newsboat, SQLite3, Pandoc and mmcdole's gofeed package have provided excellent
free software to the planet. - RSD, 2023-10-07

`
)

func main() {
	appName := path.Base(os.Args[0])
	showHelp, showVersion, showLicense := false, false, false
	fetch, display := false, false
	prune, limit := "", 0
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&fetch, "fetch", fetch, "import feed content into database")
	flag.BoolVar(&display, "display", display, "show the items in reverse chronologically order")
	flag.StringVar(&prune, "prune", prune, "remove items older than this value from the database")
	flag.IntVar(&limit, "limit", limit, "limit the number of items output")
	flag.Parse()

	args := flag.Args()

	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintln(out, skimmer.FmtHelp(helpText, appName, skimmer.Version, skimmer.ReleaseDate, skimmer.ReleaseHash))
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintf(out, "%s\n", skimmer.LicenseText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, skimmer.Version, skimmer.ReleaseHash)
		os.Exit(0)
	}

	app, err := skimmer.NewSkimmer(out, eout, appName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
	app.Fetch = fetch
	app.Display = display
	app.Limit = limit
	app.Prune = prune
	if err := app.Run(out, eout, args); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
