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

{app_name} [OPTIONS] FILENAME [TIME_RANGE]

# DESCRIPTION

{app_name} is a extremely lightweight feed reader. It trys to do two things well.

1. Download feed content
2. Display feed content

To download feed content you need a list of URLs. A list of URLs can be provided
in [newsboat's](https://newsboat.org) urls file format or as an [OPML](http://opml.org/) 
file. 

If either type of these files is provided on the command line then the file will be read
and a similarly named SQLite3 database will be created with a `+"`"+`.skim`+"`"+` extension.
{app_name} will then display the downloaded content.

After populating your skimmer database you can update it using the `+"`"+`-fetach`+"`"+`
option or read it by providing the skimmer file instead of a urls file or OPML file.


{app_name} followed by a skimmer file displays the harvested content in reverse 
chronological order. The content displayed is not "paged" unless you use the 
`+"`"+`-interactive`+"`"+` option. Typically you would use {app_name} in 
conjunction with the POSIX command 'more' or GNU command 'less'. If you provide
a time range then only items published or updated in that time range will be display.
If you only include one timestamp then the items starting with that published or updates
times will be display.

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

-limit N
: display the N most recent items.


-prune 
: The deletes items from the items table for the skimmer file provided. If a time range is provided
then the items in the time range will be deleted. If a single time is provided everything older than
that time is deleted.  A time can be specified in several ways. An alias of "today" would remove all
items older than today (this is the oposite behavior of reading items). If "now" is specified then
all items older then the current time would be removed. Otherwise time can be specified as a date
in YYYY-MM-DD format or timestamp YYYY-MM-DD HH:MM:SS format.

-i, -interactive
: display an item and prompt for next action. e.g. (n)ext, (p)rev, (s)ave, (/)search, (d)elete

-urls
: Output the contents of the SQLite 3 database channels table as a newsboat URLs list

-opml
: Output the contents of the SQLite 3 database channels table as an OPML file.

# EXAMPLE

Create a "my-news.skim" database from "my-news.opml".

~~~
{app_name] my-news.opml
~~~

Now that my-news.skim exists we can read it with

~~~
{app_name} my-news.skim
~~~

Update and read the my-news.skim file.

~~~
skimmer -fetch my-news.skim
skimmer my-news.skim
~~~


{app_name} can prune it's own database and also limit the count of items displayed.
In this example we're pruning all the items older than today and displaying the recent
five items.

~~~
{app_name} -prune today -limit 5 my-news.skim
~~~

If I limit the number of items I am reading to about 100 or so I've found
that this combination works nice.

~~~
{app_name} -limit 100 my-news.skim | pandoc -f markdown -t plain | less -R
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
	fetch, interactive, urls, opml := false, false, false, false
	prune, limit := false, 0
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&fetch, "fetch", fetch, "import feed content into database")
	flag.BoolVar(&interactive, "i", interactive, "interactively display items one at a time in reverse chronologically order")
	flag.BoolVar(&interactive, "interactive", interactive, "interactively display items one at a time in reverse chronologically order")
	flag.BoolVar(&prune, "prune", prune, "remove items in the skimmer file for the time range provided")
	flag.IntVar(&limit, "limit", limit, "limit the number of items output")
	flag.BoolVar(&urls, "urls", urls, "output the substribed feeds in newsboat's urls file format.")
	flag.BoolVar(&opml, "opml", opml, "output the substribed feeds in OPML format.")
	flag.Parse()

	args := flag.Args()

	in := os.Stdin
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
	// Setup our options
	app.Fetch = fetch
	app.Limit = limit
	app.Prune = prune
	app.Interactive = interactive
	app.AsOPML = opml
	app.AsURLs = urls
	if err := app.Run(in, out, eout, args); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
