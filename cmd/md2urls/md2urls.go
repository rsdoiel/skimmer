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

	for _, filePath := range args {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)

		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindAllStringSubmatch(line, -1)

			for _, match := range matches {
				linkText := match[1]
				url := match[2]
				fmt.Printf("%s \"~%s\"\n", url, linkText)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		}
	}
}
