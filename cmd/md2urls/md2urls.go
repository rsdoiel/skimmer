/*
    md2urls.go is part of Skimmer package. Skimmer is a package for working with feeds and rendering Link Blogs
	Copyright (C) 2025  R. S. Doiel

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

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

{app_name} [OPTIONS] MARKDOWN_FILE [MARKDOWN_FILE ...]

# DESCRIPTION

{app_name} reads in one or more ComomonMark/Markdown files containing a list of links and
returns a new URL file written to instandard output.

# EXAMPLE

Here's an example Markdown file.

~~~Markdown

This is a conceptual Markdown document holding links

- [one](https://one.example.edu)
- [two](https://two.example.edu)
- [three](https://three.example.edu)

~~~

Running the follow will result a URL file formated output for a Markdown
file name "myfile.md".

~~~shell
{app_name} myfile.md >myfile.urls
~~~

This is the resulting "myfile.urls"

~~~text
https://one.example.edu "~one"
https://two.example.edu "~two"
https://three.example.edu "~three"
~~~

`
)

func main() {
	appName := path.Base(os.Args[0])

	showHelp, showLicense, showVersion := false, false, false
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
  
	if len(os.Args) < 2 {
		fmt.Fprintf(eout, "%s\n", skimmer.FmtHelp(helpText, appName, skimmer.Version, skimmer.ReleaseDate, skimmer.ReleaseHash))
		os.Exit(1)
	}

	app := skimmer.NewMdToUrls()
	if err := app.Run(out, eout, args); err != nil {
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	os.Exit(0)
}
