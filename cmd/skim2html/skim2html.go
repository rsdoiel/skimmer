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

{app_name} [OPTIONS] SKIM_DB_FILENAME [YAML_CONFIG]

# DESCRIPTION

{app_name} reads a skimmer DB and related YAML file (configuration)
writes the saved or tagged items to the display in HTML.  This includes embedded
audio and video elements from Podcasts and Audio casts.

# CONFIGURATION

The configuration file is assumed to have a similarly named YAML file where the
file extension is ".yaml" instead of ".skim". If YAML_CONFIG is provided it'll be
used.

The configuration file supports the following attributes.

title
: (optional) This title of the feed and resulting page

description
: (optional) Any additional description (included in the head, meta element)

CSS
: (optional) Path to CSS file (use `+"`"+`@import`+"`"+` to include other CSS)

Modules
: (optional) a list of paths to ES6 modules to include in the page

Header
: (optional) HTML markup for header (if not included one will be generated from title and timestamp)

Nav
: (optional) HTML markup to be included for navigation

Footer
: (optional) HTML markup to be included in the footer


# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

# EXAMPLE

In the example we fetch URL content, read some feeds, save some interactively
then use {app_name} to generate a webpage of saved or tagged items.

~~~
skimmer myfeeds.urls
skimmer -i myfeeds.skim
{app_name} myfeeds.skim >save_items.html
~~~

You could then process the "saved_items.md" file further with Pandoc.

`

)

func main() {
	appName := path.Base(os.Args[0])

	showHelp, showLicense, showVersion := false, false, false
	title := ""
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
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

	app, err := skimmer.NewSkim2Html(appName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
    }
	app.Title = title

	if err := app.Run(out, eout, args); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
