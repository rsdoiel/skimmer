package skimmer

import (
	"bufio"
	"encoding/xml"
	"io"
	"fmt"
	"os"
	"regexp"
	"strings"
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

type UrlsToOpml struct {

}


func NewUrlsToOpml() *UrlsToOpml {
	return new(UrlsToOpml)
}

func (app *UrlsToOpml) Run(out io.Writer, feedTitle string, args []string) error {
	var outlines []Outline
	for _, filePath := range args {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file: %v\n", err)
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
			return fmt.Errorf("Error reading file: %v\n", err)
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
		return fmt.Errorf("Error marshaling XML: %v\n", err)
	}

	// Print XML declaration and OPML content
	fmt.Fprintln(out, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(out, string(output))
	return nil
} 