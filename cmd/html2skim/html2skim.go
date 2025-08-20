package main

import (
	"flag"
	"os"
	"path"
	"fmt"

	// Application package
	"github.com/rsdoiel/skimmer"
)

var (
	helpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME 

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SKIM_DB_FILE URL [CSS_SELECTOR]

# DESCRIPTION

{app_name} provides a way to treat a webpage containing links as a feed
storing the results in a skimmer database. A SKIM_DB_FILE will be
created if it does not already exist. Feed channel information can
be specified via {app_name} options. The URL should be for an 
HTML page that {app_name} will scrape for links. An optional
CSS_SELECTOR can be included to filter a specific section of the
HTML document for links. If none is provided then the selector
`+"`"+`a[href]`+"`"+` will be used.

# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

-title
: this is the channel title to use for the links scraped in the page

-description
: this is the channel description to use for links scraped in the page

-link
: this is the link to associated with the channel. You'd set this if
you were going to take the scaped links and turn them into RSS 2.0
documents.

# EXAMPLE

This is an example of scaping a web page identified
as "https://example.edu/" filter the links found in the 
".college-news" element.

~~~
{app_name}  \
     myfeeds.skim \
	 https://example.edu/ \
     '.college-news > a[href]'
~~~

`

)

func main() {
	appName := path.Base(os.Args[0])

	showHelp, showLicense, showVersion := false, false, false
	title, description, link := "", "", ""
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.StringVar(&title, "title", title, "set the channel title")
	flag.StringVar(&description, "description", description, "set the channel description")
	flag.StringVar(&link, "link", link, "the channel link")
	flag.Parse()

	args := flag.Args()
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", skimmer.FmtHelp(helpText, appName, skimmer.Version, skimmer.ReleaseDate, skimmer.ReleaseHash))
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

	app, err := skimmer.NewHtml2Skim(appName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
    }
	if err := app.Run(out, eout, args, title, description, link); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
