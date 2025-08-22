/*
    skimdups.go is part of Skimmer package. Skimmer is a package for working with feeds and rendering Link Blogs
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
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"

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

{app_name} [OPTIONS] URL_FILE [URL_FILE ...]

# DESCRIPTION

{app_name} reads in one or more URL files and scans them for duplicate URLs. It writes a report
to standard output indication the URL duplicated along with the files and line number where the
duplication is found.

# EXAMPLE

Look for duplicate URLs in page1.urls and page2.urls.

~~~shell
{app_name} page1.urls page2.urls
~~~

`
)

// URLLocation holds the file path and line number of a URL occurrence
type URLLocation struct {
	filePath   string
	lineNumber int
}

func parseFile(filePath string) (map[string][]URLLocation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urlMap := make(map[string][]URLLocation)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := line

		// Skip empty lines and comments
		if len(trimmedLine) == 0 || trimmedLine[0] == '#' {
			continue
		}

		// Regex to match URL and optional label
		re := regexp.MustCompile(`^(\S+)\s*(?:"([^"]*)")?$`)
		matches := re.FindStringSubmatch(trimmedLine)

		if len(matches) >= 2 {
			url := matches[1]
			location := URLLocation{filePath: filePath, lineNumber: lineNumber}
			urlMap[url] = append(urlMap[url], location)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urlMap, nil
}

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

	app := skimmer.NewSkimDups()
	if err := app.Run(out, eout, args); err != nil {
		fmt.Fprintln(eout, err)
		os.Exit(1)
	}
	os.Exit(0)
}
