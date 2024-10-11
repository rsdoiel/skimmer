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
	helpText = `%{app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME 

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SKIM_DB_FILENAME

# DESCRIPTION

{app_name} reads a skimmer DB and writes the saved or tagged items 
to the display in a Markdown friendly way.  This includes embedded
audio and video elements from Podcasts and Audio casts.

# OPTIONS

-help
: display help

-version
: display version info

-license
: display license

-title
: Set a page title to be included in the output of saved items

-frontmatter
: add frontmatter to Markdown output

-pocket
: add "save to pocket" link for each RSS Item displayed

# EXAMPLE

In the example we fetch URL content, read some feeds, save some interactively
then use {app_name} to generate a webpage of saved or tagged items.

~~~
skimmer myfeeds.urls
skimmer -i myfeeds.skim
{app_name} myfeeds.skim >save_items.md
~~~

You could then process the "saved_items.md" file further with Pandoc.

`

)

func main() {
	appName := path.Base(os.Args[0])

	showHelp, showLicense, showVersion := false, false, false
	frontmatter, pocket := false, false
	title := ""
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.StringVar(&title, "title", title, "set the page title for output")
	flag.BoolVar(&frontmatter, "frontmatter", frontmatter, "add frontmatter to output")
	flag.BoolVar(&pocket, "pocket", frontmatter, "add 'save to pocket' link for RSS items")
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

	app, err := skimmer.NewSkim2Md(appName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
    }
	app.Title = title

	if err := app.Run(out, eout, args, frontmatter, pocket); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
