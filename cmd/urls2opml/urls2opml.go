package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"encoding/xml"

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

{app_name} reads in one or more URL file and writes an OPML file to standard output.

# EXAMPLE

Here's of converting "myfile.urls" to "myfile.opml".


~~~shell
{app_name} myfile.urls >myfile.opml
~~~

`
)

// OPML represents the overall structure of an OPML document
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head contains metadata about the OPML file
type Head struct {
	Title string `xml:"title"`
}

// Outline represents each entry in the OPML document
type Outline struct {
	Text string `xml:"text,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
	HtmlURL  string `xml:"htmlUrl,attr,omitempty"` // This field is not standard in OPML but allows us to include URLs
	Type string `xml:"type,attr,omitempty"`
	Version string `xml:"version,attr,omitempty"`
	XmlURL string `xml:"xmlUrl,attr,omitempty"`

}

// Body contains the outline elements
type Body struct {
	Outlines []Outline `xml:"outline"`
}


func main() {
	appName := path.Base(os.Args[0])

	showHelp, showLicense, showVersion := false, false, false
	feedTitle := ""
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.StringVar(&feedTitle, "title", feedTitle, "set the OPML title for feed")
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
	if feedTitle == "" {
		feedTitle = args[0]
	}

	var outlines []Outline
	for _, filePath := range args {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(`^(\S+)\s+"(~[^"]+)"$`)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Skip empty lines and comments
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				url := matches[1]
				label := matches[2]
				obj := Outline{
					Text: label,
				}
				switch  {
				case strings.HasSuffix(url, ".rss") || strings.HasSuffix(url, ".xml"):
					obj.Type = "RSS"
					obj.Version = "RSS2"
					obj.XmlURL = url
					break;
				case strings.HasSuffix(url, ".atom"):
					obj.Type = "Atom"
					obj.XmlURL = url
				case strings.HasSuffix(url, ".html") || strings.HasSuffix(url, ".htm"):
					obj.Type = "HTML"
					obj.HtmlURL = url
				}
				if (obj.XmlURL != "" || obj.HtmlURL != "" ) && obj.Text != "" {
					outlines = append(outlines, obj)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
	}

	// Create OPML structure
	opml := OPML{
		Version: "2.0",
		Head: Head{
			Title: feedTitle,
		},
		Body: Body{
			Outlines: outlines,
		},
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(opml, "", "  ")
	if err != nil {
		fmt.Fprintf(eout, "Error marshaling XML: %v\n", err)
		os.Exit(1)
	}

	// Print XML declaration and OPML content
	fmt.Fprintln(out, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(out, string(output))
}
